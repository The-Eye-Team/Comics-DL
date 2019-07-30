package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	hosts["e-hentai.org"] = HostVal{2, s03GetComic}
}

func s03GetComic(wg *sync.WaitGroup, b *BarProxy, host string, id string, path string, outputDir string) {
	defer wg.Done()

	d := getDoc("https://" + host + path + "?p=0")
	g := d.Find(".ptds").Eq(0).Parent().Children().Length() - 2
	n := fixTitleForFilename(id + " -- " + trim(d.Find("#gn").Text()))
	f := 0

	b.AddToTotal(1)
	lp := s03GetListPage(id, d, f, b, outputDir)
	if lp == -1 {
		b.FinishNow()
		return
	}
	f += lp

	for i := 1; i < g; i++ {
		is := strconv.FormatInt(int64(i), 10)
		gd := getDoc("https://" + host + path + "?p=" + is)
		f += s03GetListPage(id, gd, f, b, outputDir)
	}

	dir1 := fmt.Sprintf("%s/jpg/%s/", outputDir, n)
	dir2 := fmt.Sprintf("%s/cbz/", outputDir)
	finp := dir2 + n + ".cbz"
	packCbzArchive(dir1, finp, b)
}

func s03GetListPage(id string, d *goquery.Document, from int, b *BarProxy, outputDir string) int {
	s := d.Find(".gdtm a")
	l := s.Length()
	n := trim(d.Find("#gn").Text())
	n = id + " -- " + n
	n = strings.Replace(n, "|", "-", -1)

	dir2 := fmt.Sprintf("%s/cbz/", outputDir)
	os.MkdirAll(dir2, os.ModePerm)
	finp := dir2 + n + ".cbz"

	if doesFileExist(finp) {
		return -1
	}

	dir1 := fmt.Sprintf("%s/jpg/%s/", outputDir, n)
	os.MkdirAll(dir1, os.ModePerm)

	b.AddToTotal(l)
	s.Each(func(i int, el *goquery.Selection) {
		defer b.Increment(1)
		v, _ := el.Attr("href")
		fp := fmt.Sprintf("%s%03d.jpg", dir1, from+i)
		if doesFileExist(fp) {
			return
		}
		s03GetPage(v, fp)
	})

	return l
}

func s03GetPage(urlS string, fpath string) {
	d := getDoc(urlS)
	s, _ := d.Find("#img").Attr("src")
	r := doRequest(s)
	b, _ := ioutil.ReadAll(r.Body)
	ioutil.WriteFile(fpath, b, os.ModePerm)
}
