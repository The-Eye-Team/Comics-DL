package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/The-Eye-Team/Comics-DL/pkg/idata"
	"github.com/The-Eye-Team/Comics-DL/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
	flag "github.com/spf13/pflag"

	_ "github.com/The-Eye-Team/Comics-DL/pkg/sites"
)

func main() {
	flagConcur := flag.IntP("concurrency", "c", 5, "The number of files to download simultaneously.")
	flagOutDir := flag.StringP("output-dir", "o", "./results/", "Output directory")
	flagKeepJpg := flag.BoolP("keep-jpg", "k", false, "Flag to keep/delete .jpg files of individual pages.")
	flagURL := flag.StringP("url", "u", "", "URL of comic to download.")
	flagFile := flag.StringP("file", "f", "", "Path to txt file with list of links to download.")
	flag.Parse()

	outDir, _ := filepath.Abs(*flagOutDir)
	outDir = strings.Replace(outDir, string(filepath.Separator), "/", -1)

	mbpp.Init(*flagConcur)

	idata.KeepJpg = *flagKeepJpg

	idata.Wg = new(sync.WaitGroup)
	idata.C = 0

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
			for idata.C == *flagConcur-10 {
				idata.Wg.Wait()
				idata.C = 0
			}
			idata.C++
			idata.Wg.Add(1)
			go iutil.DoSite(urlO, outDir)
		}
	}

	time.Sleep(time.Second / 2)
	mbpp.Wait()
	onClose()
}

func onClose() {
	fmt.Println("Completed download with", mbpp.GetTaskCount(), "tasks and", util.ByteCountIEC(mbpp.GetTaskDownloadSize()), "bytes.")
}
