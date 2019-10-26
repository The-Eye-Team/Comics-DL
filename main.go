package main

import (
	"archive/zip"
	"bufio"
	"context"
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

	"golang.org/x/sync/semaphore"
)

type HostVal struct {
	idPathIndex  int
	downloadFunc func(*BarProxy, string, string, string, string)
}

var (
	hosts     = map[string]HostVal{}
	keepJpg   bool
	doneWg    = new(sync.WaitGroup)
	progress  = mpb.New(mpb.WithWidth(64), mpb.WithWaitGroup(doneWg))
	taskIndex = 0
	guard     *semaphore.Weighted
	ctx       = context.TODO()
	bytesDLd  int64
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
	outDir += "/"

	guard = semaphore.NewWeighted(int64(*flagConcur))
	keepJpg = *flagKeepJpg

	if len(*flagURL) > 0 {
		urlO, err := url.Parse(*flagURL)
		if err != nil {
			log("URL parse error. Aborting!")
			return
		}
		doSite(urlO, outDir)
	}

	if len(*flagFile) > 0 {
		if !doesFileExist(*flagFile) {
			log("Unable to reach file!")
			return
		}
		pth, _ := filepath.Abs(*flagFile)
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

	progress.Wait()

	fmt.Println("Completed download after:")
	fmt.Println(F("\t%d tasks", taskIndex))
	fmt.Println(F("\t%s saved", byteCountIEC(bytesDLd)))
}

func doSite(place *url.URL, rootDir string) {
	h, ok := hosts[place.Host]
	if !ok {
		log("Site not supported!")
		return
	}

	job := strings.Split(place.Path, "/")[h.idPathIndex]
	bar := createBar(place.Host, job)
	go h.downloadFunc(&bar, place.Host, job, place.Path, rootDir+place.Host)
}

func createBar(host string, name string) BarProxy {
	guard.Acquire(ctx, 1)
	taskIndex++
	task := F("Task #%d:", taskIndex)
	return BarProxy{
		0,
		progress.AddBar(0,
			mpb.BarRemoveOnComplete(),
			mpb.PrependDecorators(
				decor.Name(task, decor.WC{W: len(task) + 1, C: decor.DidentRight}),
				decor.Name(host, decor.WCSyncSpaceR),
				decor.Name(name, decor.WCSyncSpaceR),
				decor.Name(": ", decor.WC{W: 2}),
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

func reduceNumber(input int64, unit int64, base string, prefixes string) string {
	if input < unit {
		return F("%d "+base, input)
	}
	div, exp := int64(unit), 0
	for n := input / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return F("%.1f %ci", float64(input)/float64(div), prefixes[exp]) + base
}

func byteCountIEC(b int64) string {
	return reduceNumber(b, 1024, "B", "KMGTPEZY")
}
