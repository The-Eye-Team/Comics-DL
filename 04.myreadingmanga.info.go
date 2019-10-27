package main

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

func init() {
	hosts["myreadingmanga.info"] = itypes.HostVal{1, func(host string, id string, path string, outputDir string) func(*mbpp.BarProxy, *sync.WaitGroup) {

		savePage := func(p int, d *goquery.Document, b *mbpp.BarProxy, dir string) {
			s := d.Find("div.entry-content img.img-myreadingmanga")
			b.AddToTotal(int64(s.Length()))

			s.Each(func(i int, el *goquery.Selection) {
				pn := padPgNum(p) + "_" + padPgNum(i)
				//
				urlS, _ := el.Attr("data-lazy-src")
				mbpp.CreateDownloadJob(urlS, dir+"/"+pn+".jpg", mbpp.BlankWaitGroup(), b)
			})
		}

		return func(bar *mbpp.BarProxy, _ *sync.WaitGroup) {

			d := getDoc("https://" + host + "/" + id + "/")
			t := strings.TrimSuffix(d.Find("title").Text(), " - MyReadingManga")

			dir := outputDir + "/" + fixTitleForFilename(t)
			out := dir + ".cbz"
			if util.DoesFileExist(out) {
				return
			}
			os.MkdirAll(dir, os.ModePerm)

			savePage(1, d, bar, dir)
			c := d.Find("a.post-page-numbers")
			if c.Length() > 0 {
				end := parseInt(c.Eq(c.Length() - 2).Text())
				for i := 2; i < end; i++ {
					is := strconv.Itoa(i)
					savePage(i, getDoc("https://"+host+"/"+id+"/"+is+"/"), bar, dir)
				}
			}

			bar.Wait()
			bar.AddToTotal(1)
			packCbzArchive(dir, out, bar)
		}
	}}
}
