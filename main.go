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

type HostVal struct {
	idPathIndex  int
	downloadFunc func(string, string, string)
}

var (
	hosts     = map[string]HostVal{}
	outputDir string
	waitgroup = new(sync.WaitGroup)
	concurr   int
	count     int
	keepJpg   bool
)

func main() {
	flagConcur := flag.IntP("concurrency", "c", 4, "The number of files to download simultaneously.")
	flagOutDir := flag.StringP("output-dir", "o", "./results", "Output directory")
	flagKeepJpg := flag.BoolP("keep-jpg", "k", false, "Flag to keep/delete .jpg files of individual pages.")
	flagURL := flag.StringP("url", "u", "", "URL of comic to download.")
	flag.Parse()

	outputDir, _ = filepath.Abs(*flagOutDir)
	outputDir = strings.Replace(outputDir, string(filepath.Separator), "/", -1)
	outputDir += "/"
	log("Saving all files to", outputDir)

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
}

func doSite(place *url.URL) {
	h, ok := hosts[place.Host]
	if !ok {
		log("Site not supported!")
		return
	}
	go h.downloadFunc(wg, bar, place.Host, strings.Split(place.Path, "/")[h.idPathIndex], place.Path)
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
