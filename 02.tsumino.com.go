package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/valyala/fastjson"
)

func init() {
	hosts["tsumino.com"] = HostVal{3, s02GetComic}
}

func s02GetComic(wg *sync.WaitGroup, b *BarProxy, host string, id string, path string, outputDir string) {
	defer wg.Done()

	d := getDoc("https://" + host + "/Book/Info/" + id)

	n0 := trim(d.Find("#Title").Text())
	n1 := strings.Split(n0, " / ")[0]
	n2 := strings.Replace(n1, " | ", " -- ", -1)
	n := strings.Split(n2, " ---- ")[0]
	n = id + " -- " + n

	dir2 := outputDir + "/cbz/"
	os.MkdirAll(dir2, os.ModePerm)

	finp := dir2 + n + ".cbz"
	if !doesFileExist(finp) {
		images := s02GetPageURLs(id)
		ln := len(images)

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
				bys, _ := ioutil.ReadAll(res.Body)
				ioutil.WriteFile(pth, bys, os.ModePerm)
			}

			packCbzArchive(dir, finp, b)
		}
	}
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
