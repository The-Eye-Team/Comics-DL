package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	hosts["myreadingmanga.info"] = HostVal{1, s04GetComic}
}

func s04GetComic(b *BarProxy, host string, id string, path string, outputDir string) {

	from := 0
	n := ""
	for i := 1; true; i++ {
		t, o := s04GetComicList(host, id, i, from, outputDir)
		if n == "" {
			n = o
		}
		if t == -2 {
			return
		}
		if t == -1 {
			break
		}
		from += t
	}

	dir1 := fmt.Sprintf("%s/jpg/%s/", outputDir, id)
	dir2 := fmt.Sprintf("%s/cbz/", outputDir)
	finp := dir2 + n + ".cbz"
	packCbzArchive(dir1, finp, b)
}

func s04GetComicList(host string, id string, page int, from int, outputDir string) (int, string) {
	d := getDoc(F("https://%s/%s/%d/", host, id, page))
	n := fixTitleForFilename(d.Find("h1.entry-title").Text())

	dir2 := outputDir + "/cbz/"
	os.Mkdir(dir2, os.ModePerm)

	finp := dir2 + n + ".cbz"
	if doesFileExist(finp) {
		return -2, n
	}

	dir1 := outputDir + "/jpg/" + id + "/"
	os.MkdirAll(dir1, os.ModePerm)

	g := d.Find("div.entry-content div")
	e := false
	p := g.Length()
	f := from
	g.Each(func(i int, el *goquery.Selection) {
		cl, o := el.Attr("class")
		if o && cl == "entry-pagination" {
			p--
			e = true
			if el.Children().Last().Is("span") {
				e = false
			}
			return
		}

		pfn := F(dir1+"%04d.jpg", f)
		if doesFileExist(pfn) {
			return
		}
		f++
		u, ex := el.Children().Eq(0).Attr("data-lazy-src")
		if !ex {
			return
		}
		res := doRequest(u)
		bys, _ := ioutil.ReadAll(res.Body)
		ioutil.WriteFile(pfn, bys, os.ModePerm)
	})

	if e {
		return p, n
	}
	return -1, n
}
