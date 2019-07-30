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
	flag "github.com/spf13/pflag"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type HostVal struct {
	idPathIndex  int
	downloadFunc func(*sync.WaitGroup, *BarProxy, string, string, string, string)
}

var (
	hosts     = map[string]HostVal{}
	rootDir   string
	waitgroup = new(sync.WaitGroup)
	concurr   int
	count     int
	keepJpg   bool
)

var (
	doneWg    = new(sync.WaitGroup)
	progress  = mpb.New(mpb.WithWidth(64), mpb.WithWaitGroup(doneWg))
	bars      []*BarProxy
	taskIndex = 1
)

func main() {
	flagConcur := flag.IntP("concurrency", "c", 4, "The number of files to download simultaneously.")
	flagOutDir := flag.StringP("output-dir", "o", "./results", "Output directory")
	flagKeepJpg := flag.BoolP("keep-jpg", "k", false, "Flag to keep/delete .jpg files of individual pages.")
	flagURL := flag.StringP("url", "u", "", "URL of comic to download.")
	flag.Parse()

	outDir, _ := filepath.Abs(*flagOutDir)
	outDir = strings.Replace(outDir, string(filepath.Separator), "/", -1)
	outDir += "/"
	rootDir = outDir

	concurr = *flagConcur
	keepJpg = *flagKeepJpg

	if len(*flagURL) > 0 {
		urlO, err := url.Parse(*flagURL)
		if err != nil {
			log("URL parse error. Aborting!")
			return
		}
		doSite(urlO)
	}

	progress.Wait()
}

func doSite(place *url.URL) {
	h, ok := hosts[place.Host]
	if !ok {
		log("Site not supported!")
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	job := strings.Split(place.Path, "/")[h.idPathIndex]
	bar := createBar(job)
	go h.downloadFunc(wg, &bar, place.Host, job, place.Path, rootDir+place.Host)
	bars = append(bars, &bar)
}

func createBar(name string) BarProxy {
	task := fmt.Sprintf("Task #%d:", taskIndex)
	taskIndex++
	return BarProxy{
		0,
		progress.AddBar(0,
			mpb.PrependDecorators(
				decor.Name(task, decor.WC{W: len(task) + 1, C: decor.DidentRight}),
				decor.Name(name, decor.WCSyncSpaceR),
				decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
			),
			mpb.AppendDecorators(
				decor.OnComplete(decor.Percentage(decor.WC{W: 5}), ""),
			),
		),
	}
}

func getDoc(urlS string) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(doRequest(urlS).Body)
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
	fmt.Print("[" + time.Now().UTC().String()[5:19] + "] ")
	fmt.Println(message...)
}

func doRequest(urlS string) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, urlS, strings.NewReader(""))
	req.Header.Add("User-Agent", "The-Eye-Team/Comics-DL/1.0")
	res, _ := http.DefaultClient.Do(req)
	return res
}

func doesDirectoryExist(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !s.IsDir() {
		return false
	}
	return true
}

// F is an shorthand alias to fmt.Sprintf
func F(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func packCbzArchive(dirIn string, fileOut string, bar *BarProxy) {
	outf, _ := os.Create(fileOut)
	outz := zip.NewWriter(outf)
	files, _ := ioutil.ReadDir(dirIn)
	for _, item := range files {
		zw, _ := outz.Create(item.Name())
		bs, _ := ioutil.ReadFile(dirIn + item.Name())
		zw.Write(bs)
	}
	outz.Close()
	if !keepJpg {
		os.RemoveAll(dirIn)
	}
	bar.Increment(1)
}

func fixTitleForFilename(t string) string {
	n := trim(t)
	n = strings.Replace(n, ":", "", -1)
	n = strings.Replace(n, "\\", "-", -1)
	n = strings.Replace(n, "/", "-", -1)
	n = strings.Replace(n, "*", "-", -1)
	n = strings.Replace(n, "?", "-", -1)
	n = strings.Replace(n, "\"", "-", -1)
	n = strings.Replace(n, "<", "-", -1)
	n = strings.Replace(n, ">", "-", -1)
	n = strings.Replace(n, "|", "-", -1)
	return n
}
