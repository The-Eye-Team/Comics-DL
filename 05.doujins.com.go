package main

import (
	"os"
	"strings"
	"sync"

	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

func init() {
	hosts["doujins.com"] = itypes.HostVal{2, func(host string, id string, path string, outputDir string) func(*mbpp.BarProxy, *sync.WaitGroup) {
		return func(bar *mbpp.BarProxy, _ *sync.WaitGroup) {

			d := getDoc("https://" + host + path)
			t := d.Find("title").Text()
			s := d.Find("img.doujin")

			dir := outputDir + "/" + fixTitleForFilename(t)
			out := dir + ".cbz"
			if util.DoesFileExist(out) {
				return
			}
			os.MkdirAll(dir, os.ModePerm)

			bar.AddToTotal(int64(s.Length()))
			s.Each(func(i int, el *goquery.Selection) {
				urlS, _ := el.Attr("data-file")
				urlS = strings.ReplaceAll(urlS, "&amp;", "&")
				pth := dir + "/" + padPgNum(i+1) + ".jpg"
				go mbpp.CreateDownloadJob(urlS, pth, mbpp.BlankWaitGroup(), bar)
			})

			bar.Wait()
			bar.AddToTotal(1)
			packCbzArchive(dir, out, bar)
		}
	}}
}
