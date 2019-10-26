package main

import (
	"archive/zip"
	"bufio"
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
	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
	flag "github.com/spf13/pflag"
)

type HostVal struct {
	idPathIndex  int
	downloadFunc func(string, string, string, string) func(*mbpp.BarProxy, *sync.WaitGroup)
}

var (
	hosts   = map[string]HostVal{}
	keepJpg bool
)

func main() {
	flagConcur := flag.IntP("concurrency", "c", 10, "The number of files to download simultaneously.")
	flagOutDir := flag.StringP("output-dir", "o", "./results/", "Output directory")
	flagKeepJpg := flag.BoolP("keep-jpg", "k", false, "Flag to keep/delete .jpg files of individual pages.")
	flagURL := flag.StringP("url", "u", "", "URL of comic to download.")
	flagFile := flag.StringP("file", "f", "", "Path to txt file with list of links to download.")
	flag.Parse()

	outDir, _ := filepath.Abs(*flagOutDir)
	outDir = strings.Replace(outDir, string(filepath.Separator), "/", -1)

	mbpp.Init(*flagConcur)
	keepJpg = *flagKeepJpg

	if len(*flagURL) > 0 {
		urlO, err := url.Parse(*flagURL)
		if err != nil {
			util.Log("URL parse error. Aborting!")
			return
		}
		doSite(urlO, outDir)
	}

	if len(*flagFile) > 0 {
		pth, _ := filepath.Abs(*flagFile)
		if !util.DoesFileExist(pth) {
			util.Log("Unable to reach file!")
			return
		}
		file, _ := os.Open(pth)
		scan := bufio.NewScanner(file)

		for scan.Scan() {
			line := scan.Text()
			urlO, err := url.Parse(line)
			if err != nil {
				return
			}
			doSite(urlO, outDir)
		}
	}

	time.Sleep(time.Second / 2)
	mbpp.Wait()

	fmt.Println("Completed download with", mbpp.GetTaskCount(), "tasks and", util.ByteCountIEC(mbpp.GetTaskDownloadSize()), "bytes.")
}

func doSite(place *url.URL, rootDir string) {
	h, ok := hosts[place.Host]
	if !ok {
		return
	}
	id := strings.Split(place.Path, "/")[h.idPathIndex]
	job := place.Host + " / " + id
	mbpp.CreateJob(job, h.downloadFunc(place.Host, id, place.Path, rootDir+"/"+place.Host))
}

func getDoc(urlS string) *goquery.Document {
	res, _ := http.Get(urlS)
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	return doc
}

func trim(x string) string {
	return strings.Trim(x, " \n\r\t")
}

func packCbzArchive(dirIn string, fileOut string, bar *mbpp.BarProxy) {
	mbpp.CreateJob("Packing "+fileOut, func(b *mbpp.BarProxy, _ *sync.WaitGroup) {
		outf, _ := os.Create(fileOut)
		outz := zip.NewWriter(outf)
		files, _ := ioutil.ReadDir(dirIn)
		b.AddToTotal(int64(len(files) + 2))
		for _, item := range files {
			zw, _ := outz.Create(item.Name())
			bs, _ := ioutil.ReadFile(dirIn + "/" + item.Name())
			zw.Write(bs)
			b.Increment(1)
		}
		outz.Close()
		b.Increment(1)
		if !keepJpg {
			os.RemoveAll(dirIn)
		}
		b.Increment(1)
		bar.Increment(1)
	})
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

func padPgNum(n int) string {
	return fmt.Sprintf("%04d", n)
}
