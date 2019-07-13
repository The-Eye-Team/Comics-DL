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
)

const (
	s01Host = "readcomicsonline.ru"
	s02Host = "www.tsumino.com"
	s03Host = "e-hentai.org"
	s04Host = "myreadingmanga.info"
)

var (
	outputDir string
	waitgroup *sync.WaitGroup
	concurr   int
	count     = 0
	keepJpg   bool
)

func main() {
	flagComicID := flag.String("comic-id", "", "readcomicsonline.ru comic ID")
	flagConcur := flag.IntP("concurrency", "c", 4, "The number of files to download simultaneously.")
	flagOutDir := flag.StringP("output-dir", "o", "./results", "Output directory")
	flagKeepJpg := flag.BoolP("keep-jpg", "k", false, "Flag to keep/delete .jpg files of individual pages.")
	flagURL := flag.StringP("url", "u", "", "URL of comic to download.")
	flag.Parse()

	//

	outputDir, _ = filepath.Abs(*flagOutDir)
	outputDir = strings.Replace(outputDir, string(filepath.Separator), "/", -1)
	outputDir += "/"
	log("Saving all files to", outputDir)

	concurr = *flagConcur

	wg := sync.WaitGroup{}
	waitgroup = &wg

	keepJpg = *flagKeepJpg

	//

	if len(*flagComicID) > 0 {
		*flagURL = "https://readcomicsonline.ru/comic/" + *flagComicID
	}

	//

	urlO, err := url.Parse(*flagURL)
	if err != nil {
		log("URL parse error. Aborting!")
		return
	}

	switch urlO.Host {
	case s01Host:
		outputDir += s01Host
		s01GetComic(strings.Split(urlO.Path, "/")[2])
	case s02Host:
		outputDir += s02Host
		s02GetComic(strings.Split(urlO.Path, "/")[3])
	case s03Host:
		outputDir += s03Host
		s03GetComic(strings.Split(urlO.Path, "/")[2], urlO.Path)
	case s04Host:
		outputDir += s04Host
		s04GetComic(strings.Split(urlO.Path, "/")[1])
	default:
		log("Site not supported!")
		return
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
	fmt.Print("[" + time.Now().UTC().String()[0:19] + "] ")
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

func packCbzArchive(dirIn string, fileOut string) {
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
