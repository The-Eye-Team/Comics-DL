package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	flag "github.com/spf13/pflag"
)

const (
	domain = "https://readcomicsonline.ru"
)

var (
	worker sync.WaitGroup
	count  = 0
)

func main() {
	flagComic := flag.String("comic-id", "", "")
	flagConcur := flag.Int("concurrency", 4, "The number of files to download simultaneously.")
	flag.Parse()

	id := *flagComic
	if len(id) == 0 {
		log("Must send a valid comic ID")
		log(">If you'd like to download https://readcomicsonline.ru/comic/justice-league-2016")
		log(">then pass --comic-id justice-league-2016")
		return
	}
	log("Saving comic:", id)

	d := getDoc(domain + "/comic/" + id).Find("ul.chapters li")
	log("Found", d.Length(), "issues")

	d.Each(func(i int, el *goquery.Selection) {
	waitgroup := sync.WaitGroup{}
		is0, _ := el.Children().First().Children().First().Attr("href")
		is1 := strings.Split(is0, "/")
		is2 := is1[len(is1)-1]
		waitgroup.Add(1)
		count++
		go getIssue(id, is2, &waitgroup)
		if count == *flagConcur {
			waitgroup.Wait()
		}
	})
	waitgroup.Wait()
	log("Done!")
}

func getIssue(id string, issue string, wtgrp *sync.WaitGroup) {
	dir := fmt.Sprintf("./results-jpg/%s/Issue %s/", id, issue)
	os.MkdirAll(dir, os.ModePerm)
	for j := 0; true; j++ {
		pth := fmt.Sprintf("%s%02d.jpg", dir, j+1)
		if doesFileExist(pth) {
			continue
		}
		u := fmt.Sprintf("https://readcomicsonline.ru/uploads/manga/%s/chapters/%s/%02d.jpg", id, issue, j+1)
		res := doRequest(fmt.Sprintf("https://readcomicsonline.ru/uploads/manga/%s/chapters/%s/%02d.jpg", id, issue, j+1))
		if res.StatusCode >= 400 {
			break
		}
		log(u)
		bys, _ := ioutil.ReadAll(res.Body)
		ioutil.WriteFile(pth, bys, os.ModePerm)
	}
	log("Completed download of Issue", issue)
	//
	// //
	count--
	wtgrp.Done()
}

func getDoc(lru string) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(doRequest(lru).Body)
	return doc
}

func doesFileExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func log(message ...interface{}) {
	fmt.Print("[" + time.Now().UTC().String()[0:19] + "] ")
	fmt.Println(message...)
}

func doRequest(urlAsText string) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, urlAsText, strings.NewReader(""))
	res, _ := http.DefaultClient.Do(req)
	return res
}
