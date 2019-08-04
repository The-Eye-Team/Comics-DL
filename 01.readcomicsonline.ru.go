package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	hosts["readcomicsonline.ru"] = HostVal{2, s01GetComic}
}

func s01GetComic(b *BarProxy, host string, id string, path string, outputDir string) {
	defer guard.Release(1)

	d := getDoc("https://" + host + "/comic/" + id)
	s := d.Find("ul.chapters li")
	n := fixTitleForFilename(trim(d.Find("h2.listmanga-header").Eq(0).Text()))

	s.Each(func(i int, el *goquery.Selection) {
		is0, _ := el.Children().First().Children().First().Attr("href")
		is1 := strings.Split(is0, "/")
		is2 := is1[len(is1)-1]
		is3, _ := url.ParseQuery("x=" + is2)
		b.AddToTotal(1)
		go s01GetIssue(id, n, is3["x"][0], b, outputDir)
	})
	if s.Length() == 0 {
		b.FinishNow()
	}
}

func s01GetIssue(id string, name string, issue string, b *BarProxy, outputDir string) {
	defer guard.Release(1)
	bs := createBar("readcomicsonline.ru", F("%s #%s", id, issue))
	dir2 := F(outputDir+"/cbz/%s/", name)
	os.MkdirAll(dir2, os.ModePerm)
	finp := F("%sIssue %s.cbz", dir2, issue)

	dir := F(outputDir+"/jpg/%s/Issue %s/", name, issue)
	if !doesFileExist(finp) {
		os.MkdirAll(dir, os.ModePerm)
		bs.AddToTotal(1)
		for j := 1; true; j++ {
			pth := F("%s%03d.jpg", dir, j)
			if doesFileExist(pth) {
				continue
			}
			bs.AddToTotal(1)
			u := F("https://readcomicsonline.ru/uploads/manga/%s/chapters/%s/%02d.jpg", id, issue, j)
			res := doRequest(u)
			if res.StatusCode >= 400 {
				break
			}
			bys, _ := ioutil.ReadAll(res.Body)
			bytesDLd += int64(len(bys))
			ioutil.WriteFile(pth, bys, os.ModePerm)
			bs.Increment(1)
		}
		packCbzArchive(dir, finp, &bs)
	}
	bs.FinishNow()
	b.Increment(1)
}
