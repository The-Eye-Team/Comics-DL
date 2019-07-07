package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func s01GetComic(id string) {
	log("Saving comic: readcomicsonline.ru /", id)

	d := getDoc("https://" + s01Host + "/comic/" + id)
	s := d.Find("ul.chapters li")
	n := trim(d.Find("h2.listmanga-header").Eq(0).Text())
	n = strings.Replace(n, ":", "", -1)
	n = strings.Replace(n, "/", "-", -1)
	log("Found", s.Length(), "issues of", n)

	setupUIList(n, id)

	s.Each(func(i int, el *goquery.Selection) {
		is0, _ := el.Children().First().Children().First().Attr("href")
		is1 := strings.Split(is0, "/")
		is2 := is1[len(is1)-1]
		is3, _ := url.ParseQuery("x=" + is2)
		waitgroup.Add(1)
		count++
		go s01GetIssue(id, n, is3["x"][0], findNextOpenRow(is2))
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

func s01GetIssue(id string, name string, issue string, row int) {
	setRowText(row, F("[%s] Preparing...", issue))
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
			setRowText(row, F("[%s] Downloading Issue %s, Page %d", issue, issue, j))
			bys, _ := ioutil.ReadAll(res.Body)
			ioutil.WriteFile(pth, bys, os.ModePerm)
		}
		//
		setRowText(row, F("[%s] Packing archive..", issue))
		packCbzArchive(dir, finp)
	}
	setRowText(row, F("[x] Completed Issue %s.", issue))
	count--
	waitgroup.Done()
}
