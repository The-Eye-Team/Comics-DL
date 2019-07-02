package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/valyala/fastjson"
)

func s02GetComic(id string) {
	log("Saving comic: tsumino.com /", id)

	d := getDoc("https://" + s02Host + "/Book/Info/" + id)

	n0 := trim(d.Find("#Title").Text())
	n1 := strings.Split(n0, " / ")[0]
	n2 := strings.Replace(n1, " | ", " -- ", -1)
	n := strings.Split(n2, " ---- ")[0]

	setupUIList(n, id)
	setRowText(0, "Preparing...")

	dir2 := outputDir + "/cbz/"
	os.MkdirAll(dir2, os.ModePerm)

	finp := dir2 + n + ".cbz"
	if !doesFileExist(finp) {
		images := s02GetPageURLs(id)
		ln := len(images)
		setRowText(0, F("Found %d pages of %s", ln, n))

		if ln > 0 {
			dir := F("%s/jpg/%s/", outputDir, n)
			os.MkdirAll(dir, os.ModePerm)
			for i, item := range images {
				pth := F("%s/jpg/%s/%03d.jpg", outputDir, n, i)
				if doesFileExist(pth) {
					continue
				}
				itm := url.Values{}
				itm.Add("v", item)
				res := doRequest("https://www.tsumino.com/Image/Object?name=" + itm.Encode()[2:])
				setRowText(0, F("Downloading Page %d/%d", i+1, ln))
				bys, _ := ioutil.ReadAll(res.Body)
				ioutil.WriteFile(pth, bys, os.ModePerm)
			}

			setRowText(0, F("Packing archive.."))
			packCbzArchive(dir, finp)
		}
	}
	setRowText(0, "Completed.")
}

func s02GetPageURLs(id string) []string {
	prm := url.Values{}
	prm.Add("q", id)
	req, _ := http.NewRequest(http.MethodPost, "https://www.tsumino.com/Read/Load", strings.NewReader(prm.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://www.tsumino.com/Read/View/"+id)
	res, _ := http.DefaultClient.Do(req)
	bys, _ := ioutil.ReadAll(res.Body)
	fjv, _ := fastjson.Parse(string(bys))
	arr := fjv.GetArray("reader_page_urls")

	val := []string{}
	for _, item := range arr {
		str := strings.TrimSuffix(strings.TrimPrefix(item.String(), "\""), "\"")
		val = append(val, str)
	}
	return val
}
