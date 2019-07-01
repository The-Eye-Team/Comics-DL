package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	flag "github.com/spf13/pflag"
)

const (
	domain = "https://readcomicsonline.ru"
)

var (
	outputDir string
	waitgroup *sync.WaitGroup
	uilist    *widgets.List
	count     = 0
)

func main() {
	flagComic := flag.String("comic-id", "", "")
	flagConcur := flag.Int("concurrency", 4, "The number of files to download simultaneously.")
	flagOutDir := flag.String("output-dir", "./results", "Output directory")
	flag.Parse()

	id := *flagComic
	if len(id) == 0 {
		log("Must send a valid comic ID")
		log(">If you'd like to download https://readcomicsonline.ru/comic/justice-league-2016")
		log(">then pass --comic-id justice-league-2016")
		return
	}
	log("Saving comic:", id)

	outputDir, _ = filepath.Abs(*flagOutDir)
	log("Saving all files to", outputDir)

	//

	if err := termui.Init(); err != nil {
		log("failed to initialize termui:", err)
	}
	defer termui.Close()

	//

	d := getDoc(domain + "/comic/" + id)
	s := d.Find("ul.chapters li")
	n := trim(d.Find("h2.listmanga-header").Eq(0).Text())
	log("Found", s.Length(), "issues of", n)

	//

	uilist = widgets.NewList()
	uilist.Title = "Comics-DL Progress of " + n + " [" + id + "] ---- " + outputDir + " "
	uilist.Rows = strings.Split(strings.Repeat("[x] ,", *flagConcur), ",")
	uilist.WrapText = false
	uilist.SetRect(0, 0, 100, *flagConcur*3)
	termui.Render(uilist)

	//

	wg := sync.WaitGroup{}
	waitgroup = &wg
	s.Each(func(i int, el *goquery.Selection) {
		is0, _ := el.Children().First().Children().First().Attr("href")
		is1 := strings.Split(is0, "/")
		is2 := is1[len(is1)-1]
		is3, _ := url.ParseQuery("x=" + is2)
		waitgroup.Add(1)
		count++
		go getIssue(id, n, is3["x"][0], findNextOpenRow(is2))
		if count == *flagConcur {
			waitgroup.Wait()
		}
	})
	waitgroup.Wait()
	log("Done!")
}

func getIssue(id string, name string, issue string, row int) {
	setRowText(row, fmt.Sprintf("[%s] Preparing...", issue))
	dir2 := fmt.Sprintf(outputDir+"/cbz/%s/", name)
	os.MkdirAll(dir2, os.ModePerm)
	finp := fmt.Sprintf("%sIssue %s.cbz", dir2, issue)

	if !doesFileExist(finp) {
		dir := fmt.Sprintf(outputDir+"/jpg/%s/Issue %s/", name, issue)
		os.MkdirAll(dir, os.ModePerm)
		for j := 1; true; j++ {
			pth := fmt.Sprintf("%s%02d.jpg", dir, j)
			if doesFileExist(pth) {
				continue
			}
			u := fmt.Sprintf("https://readcomicsonline.ru/uploads/manga/%s/chapters/%s/%02d.jpg", id, issue, j)
			res := doRequest(u)
			if res.StatusCode >= 400 {
				break
			}
			setRowText(row, fmt.Sprintf("[%s] Downloading Issue %s, Page %02d", issue, issue, j))
			bys, _ := ioutil.ReadAll(res.Body)
			ioutil.WriteFile(pth, bys, os.ModePerm)
		}
		//
		outf, _ := os.Create(finp)
		outz := zip.NewWriter(outf)
		files, _ := ioutil.ReadDir(dir)
		for _, item := range files {
			zw, _ := outz.Create(item.Name())
			bs, _ := ioutil.ReadFile(dir + item.Name())
			zw.Write(bs)
		}
		outz.Close()
	}
	setRowText(row, fmt.Sprintf("[x] Completed Issue %s.", issue))
	count--
	waitgroup.Done()
}

func getDoc(lru string) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(doRequest(lru).Body)
	return doc
}

func trim(x string) string {
	return strings.Trim(x, " \n\r\t")
}

func doesFileExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func log(message ...interface{}) {
	fmt.Print("[" + time.Now().UTC().String()[0:19] + "] ")
	fmt.Println(message...)
}

func doRequest(urlAsText string) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, urlAsText, strings.NewReader(""))
	res, _ := http.DefaultClient.Do(req)
	return res
}

func setRowText(row int, text string) {
	uilist.Rows[row] = text
	termui.Render(uilist)
}

func findNextOpenRow(iss string) int {
	for i, v := range uilist.Rows {
		if strings.HasPrefix(v, "[x]") {
			uilist.Rows[i] = "[r] Reserved for " + iss
			return i
		}
	}
	return -1
}
