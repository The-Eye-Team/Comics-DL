package main

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/The-Eye-Team/Comics-DL/pkg/idata"
	"github.com/The-Eye-Team/Comics-DL/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
	flag "github.com/spf13/pflag"

	_ "github.com/The-Eye-Team/Comics-DL/pkg/sites"
)

var (
	Version = "vMASTER"
)

func main() {
	util.Log("Starting up Comics-DL " + Version + "...")
	util.Log("Brought to you by The-Eye.eu")

	flagConcur := flag.IntP("concurrency", "c", 10, "The number of files to download simultaneously.")
	flagOutDir := flag.StringP("output-dir", "o", "./results/", "Output directory")
	flagKeepJpg := flag.BoolP("keep-jpg", "k", false, "Flag to keep/delete .jpg files of individual pages.")
	flagURL := flag.StringP("url", "u", "", "URL of comic to download.")
	flagFile := flag.StringP("file", "f", "", "Path to txt file with list of links to download.")
	flag.Parse()

	outDir, _ := filepath.Abs(*flagOutDir)
	outDir = strings.Replace(outDir, string(filepath.Separator), "/", -1)

	mbpp.Init(*flagConcur)

	idata.KeepJpg = *flagKeepJpg

	util.RunOnClose(onClose)

	if len(*flagURL) > 0 {
		urlO, err := url.Parse(*flagURL)
		if err != nil {
			util.Log("URL parse error. Aborting!")
			return
		}
		iutil.DoSite(urlO, outDir)
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
			iutil.DoSite(urlO, outDir)
		}
	}

	mbpp.Wait()
	time.Sleep(time.Second)
	onClose()
}

func onClose() {
	util.Log(mbpp.GetCompletionMessage())
}
