package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"

	. "github.com/nektro/go-util/alias"
)

func init() {
	hosts["readcomicsonline.ru"] = HostVal{2, func(host string, id string, path string, outputDir string) func(*mbpp.BarProxy, *sync.WaitGroup) {
		return func(mbar *mbpp.BarProxy, _ *sync.WaitGroup) {
			//
			d := getDoc("https://" + host + "/comic/" + id)
			s := d.Find("ul.chapters li")
			n := fixTitleForFilename(trim(d.Find("h2.listmanga-header").Eq(0).Text()))
			mbar.AddToTotal(int64(s.Length()))
			s.Each(func(i int, el *goquery.Selection) {
				is0, _ := el.Children().First().Children().First().Attr("href")
				is1 := strings.Split(is0, "/")
				is2 := is1[len(is1)-1]
				is3, _ := url.ParseQuery("x=" + is2)
				//
				name := n
				issue := is3["x"][0]
				go mbpp.CreateJob(name+" / "+issue, func(jbar *mbpp.BarProxy, wg *sync.WaitGroup) {
					defer mbar.Increment(1)
					//
					dir2 := F(outputDir+"/cbz/%s/", name)
					os.MkdirAll(dir2, os.ModePerm)
					finp := F("%sIssue %s.cbz", dir2, issue)
					dir := F(outputDir+"/jpg/%s/Issue %s/", name, issue)
					if !util.DoesFileExist(finp) {
						os.MkdirAll(dir, os.ModePerm)
						jbar.AddToTotal(1)
						for j := 1; true; j++ {
							pth := F("%s%03d.jpg", dir, j)
							if util.DoesFileExist(pth) {
								continue
							}
							u := F("https://readcomicsonline.ru/uploads/manga/%s/chapters/%s/%02d.jpg", id, issue, j)
							res, _ := http.Head(u)
							if res.StatusCode != 200 {
								break
							}
							jbar.AddToTotal(1)
							wg.Add(1)
							go mbpp.CreateDownloadJob(u, pth, wg, jbar)
						}
						wg.Wait()
						packCbzArchive(dir, finp, jbar)
					}
				})
			})
		}
	}}
}
