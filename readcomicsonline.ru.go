package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func s01GetComic(id string) {
	log("Saving comic: readcomicsonline.ru /", id)

	d := getDoc("https://" + s01Host + "/comic/" + id)
	s := d.Find("ul.chapters li")
	n := trim(d.Find("h2.listmanga-header").Eq(0).Text())
	log("Found", s.Length(), "issues of", n)

	//

	uilist = widgets.NewList()
	uilist.Title = "Comics-DL ---- " + n + " [" + id + "] ---- " + outputDir + " "
	uilist.Rows = strings.Split(strings.Repeat("[x] ,", concurr), ",")
	uilist.WrapText = false
	uilist.SetRect(0, 0, 100, concurr*2)
	termui.Render(uilist)

	//

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
		di := fmt.Sprintf(outputDir+"/"+s01Host+"/jpg/%s/", n)
		if doesDirectoryExist(di) {
			os.RemoveAll(di)
		}
	}
}

func s01GetIssue(id string, name string, issue string, row int) {
	setRowText(row, fmt.Sprintf("[%s] Preparing...", issue))
	dir2 := fmt.Sprintf(outputDir+"/"+s01Host+"/cbz/%s/", name)
	os.MkdirAll(dir2, os.ModePerm)
	finp := fmt.Sprintf("%sIssue %s.cbz", dir2, issue)

	dir := fmt.Sprintf(outputDir+"/"+s01Host+"/jpg/%s/Issue %s/", name, issue)
	if !doesFileExist(finp) {
		os.MkdirAll(dir, os.ModePerm)
		for j := 1; true; j++ {
			pth := fmt.Sprintf("%s%03d.jpg", dir, j)
			if doesFileExist(pth) {
				continue
			}
			u := fmt.Sprintf("https://readcomicsonline.ru/uploads/manga/%s/chapters/%s/%02d.jpg", id, issue, j)
			res := doRequest(u)
			if res.StatusCode >= 400 {
				break
			}
			setRowText(row, fmt.Sprintf("[%s] Downloading Issue %s, Page %d", issue, issue, j))
			bys, _ := ioutil.ReadAll(res.Body)
			ioutil.WriteFile(pth, bys, os.ModePerm)
		}
		//
		setRowText(row, fmt.Sprintf("[%s] Packing archive..", issue))
		outf, _ := os.Create(finp)
		outz := zip.NewWriter(outf)
		files, _ := ioutil.ReadDir(dir)
		for _, item := range files {
			zw, _ := outz.Create(item.Name())
			bs, _ := ioutil.ReadFile(dir + item.Name())
			zw.Write(bs)
		}
		outz.Close()
	}
	setRowText(row, fmt.Sprintf("[x] Completed Issue %s.", issue))
	count--
	waitgroup.Done()
}
