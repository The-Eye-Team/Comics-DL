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

	d := getDoc("https://" + host + "/comic/" + id)
	s := d.Find("ul.chapters li")
	n := fixTitleForFilename(trim(d.Find("h2.listmanga-header").Eq(0).Text()))

	s.Each(func(i int, el *goquery.Selection) {
		is0, _ := el.Children().First().Children().First().Attr("href")
		is1 := strings.Split(is0, "/")
		is2 := is1[len(is1)-1]
		is3, _ := url.ParseQuery("x=" + is2)
		waitgroup.Add(1)
		count++
		go s01GetIssue(id, n, is3["x"][0], b, outputDir)
		if count == concurr {
			waitgroup.Wait()
		}
	})
	waitgroup.Wait()

	if !keepJpg {
		di := F(outputDir+"/jpg/%s/", n)
		if doesDirectoryExist(di) {
			os.RemoveAll(di)
		}
	}
}

func s01GetIssue(id string, name string, issue string, b *BarProxy, outputDir string) {
	dir2 := F(outputDir+"/cbz/%s/", name)
	os.MkdirAll(dir2, os.ModePerm)
	finp := F("%sIssue %s.cbz", dir2, issue)

	dir := F(outputDir+"/jpg/%s/Issue %s/", name, issue)
	if !doesFileExist(finp) {
		os.MkdirAll(dir, os.ModePerm)
		for j := 1; true; j++ {
			pth := F("%s%03d.jpg", dir, j)
			if doesFileExist(pth) {
				continue
			}
			u := F("https://readcomicsonline.ru/uploads/manga/%s/chapters/%s/%02d.jpg", id, issue, j)
			res := doRequest(u)
			if res.StatusCode >= 400 {
				break
			}
			bys, _ := ioutil.ReadAll(res.Body)
			ioutil.WriteFile(pth, bys, os.ModePerm)
		}
		//
		packCbzArchive(dir, finp, b)
	}
	count--
	waitgroup.Done()
}
