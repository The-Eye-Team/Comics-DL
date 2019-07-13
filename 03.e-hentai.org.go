package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func s03GetComic(id string, path string) {
	log("Saving comic: e-hentai.org /", path)

	d := getDoc("https://" + s03Host + path + "?p=0")
	g := d.Find(".ptds").Eq(0).Parent().Children().Length() - 2
	n := fixTitleForFilename(id + " -- " + trim(d.Find("#gn").Text()))
	f := 0

	log("Preparing...")

	lp := s03GetListPage(id, d, f)
	if lp == -1 {
		return
	}
	f += lp

	for i := 1; i < g; i++ {
		is := strconv.FormatInt(int64(i), 10)
		gd := getDoc("https://" + s03Host + path + "?p=" + is)
		f += s03GetListPage(id, gd, f)
	}

	log(F("Packing archive.."))
	dir1 := fmt.Sprintf("%s/jpg/%s/", outputDir, n)
	dir2 := fmt.Sprintf("%s/cbz/", outputDir)
	finp := dir2 + n + ".cbz"
	packCbzArchive(dir1, finp)

	log("Completed.")
}

func s03GetListPage(id string, d *goquery.Document, from int) int {
	s := d.Find(".gdtm a")
	l := s.Length()
	n := trim(d.Find("#gn").Text())
	n = id + " -- " + n
	n = strings.Replace(n, "|", "-", -1)

	log(F("Found %d pages in %s", l, n))

	dir2 := fmt.Sprintf("%s/cbz/", outputDir)
	os.MkdirAll(dir2, os.ModePerm)
	finp := dir2 + n + ".cbz"

	if doesFileExist(finp) {
		log("Comic already saved, skipping!")
		return -1
	}

	dir1 := fmt.Sprintf("%s/jpg/%s/", outputDir, n)
	os.MkdirAll(dir1, os.ModePerm)
	s.Each(func(i int, el *goquery.Selection) {
		v, _ := el.Attr("href")
		fp := fmt.Sprintf("%s%03d.jpg", dir1, from+i)
		log(F("Downloading Page %d/%d", from+i+1, from+l))
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
