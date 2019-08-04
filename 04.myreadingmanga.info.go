package main

import (
	"io/ioutil"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	hosts["myreadingmanga.info"] = HostVal{1, s04GetComic}
}

func s04GetComic(b *BarProxy, host string, id string, path string, outputDir string) {
	defer guard.Release(1)

	b.AddToTotal(1)
	from := 0
	n := ""
	for i := 1; true; i++ {
		t, o := s04GetComicList(host, id, i, from, outputDir, b)
		if n == "" {
			n = o
		}
		if t == -2 {
			b.FinishNow()
			return
		}
		if t == -1 {
			break
		}
		from += t
	}

	dir1 := F("%s/jpg/%s/", outputDir, id)
	dir2 := F("%s/cbz/", outputDir)
	finp := dir2 + n + ".cbz"
	packCbzArchive(dir1, finp, b)
	b.FinishNow()
}

func s04GetComicList(host string, id string, page int, from int, outputDir string, b *BarProxy) (int, string) {
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

	g := d.Find("div.entry-content img")
	e := false
	p := g.Length()
	f := from
	b.AddToTotal(p)
	g.Each(func(i int, el *goquery.Selection) {
		defer b.Increment(1)

		pfn := F(dir1+"%04d.jpg", f)
		if doesFileExist(pfn) {
			return
		}
		f++
		u, ex := el.Attr("data-lazy-src")
		if !ex {
			return
		}
		res := doRequest(u)
		if res == nil {
			return
		}
		bys, _ := ioutil.ReadAll(res.Body)
		bytesDLd += int64(len(bys))
		ioutil.WriteFile(pfn, bys, os.ModePerm)
	})
	nc1 := d.Find("div.entry-content .entry-pagination")
	if nc1.Length() > 0 {
		if nc1.Last().Is("span") {
			e = true
		}
	}

	if e {
		return p, n
	}
	return -1, n
}
