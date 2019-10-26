package main

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
)

func init() {
	hosts["e-hentai.org"] = HostVal{2, s03GetComic}
}

func s03GetComic(host string, id string, path string, outputDir string) func(*mbpp.BarProxy, *sync.WaitGroup) {

	savePage := func(p int, d *goquery.Document, b *mbpp.BarProxy, dir string) {
		s := d.Find("div.gdtm a img")
		b.AddToTotal(int64(s.Length()))

		s.Each(func(i int, el *goquery.Selection) {
			pn := (40 * p) + i + 1
			//
			url1, _ := el.Parent().Attr("href")
			url2, _ := getDoc(url1).Find("#img").Attr("src")
			mbpp.CreateDownloadJob(url2, dir+"/"+padPgNum(pn)+".jpg", mbpp.BlankWaitGroup(), b)

		})
	}

	return func(bar *mbpp.BarProxy, _ *sync.WaitGroup) {

		d := getDoc("https://" + host + path + "?p=0")
		t := strings.TrimSuffix(d.Find("title").Text(), " - E-Hentai Galleries")
		bar.AddToTotal(1)

		dir := outputDir + "/" + id + "." + fixTitleForFilename(t)
		os.MkdirAll(dir, os.ModePerm)

		savePage(0, d, bar, dir)
		c := d.Find("table.ptt td")
		if c.Length() > 3 {
			pu, _ := c.Eq(c.Length() - 2).Children().Eq(0).Attr("href")
			pl, _ := strconv.ParseInt(strings.Split(pu, "=")[1], 10, 32)

			pli := int(pl)
			for i := 1; i <= pli; i++ {
				savePage(i, getDoc("https://"+host+path+"?p="+strconv.Itoa(pli)), bar, dir)
			}
		}

		packCbzArchive(dir, dir+".cbz", bar)
	}
}
