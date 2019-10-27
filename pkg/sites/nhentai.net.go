package sites

import (
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/The-Eye-Team/Comics-DL/pkg/idata"
	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"
	"github.com/The-Eye-Team/Comics-DL/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

func init() {
	idata.Hosts["nhentai.net"] = itypes.HostVal{2, func(host string, id string, path string, outputDir string) func(*mbpp.BarProxy, *sync.WaitGroup) {

		saveImage := func(i int, img *goquery.Selection, d string, b *mbpp.BarProxy) {
			urlS, _ := img.Attr("data-src")
			urlS = strings.ReplaceAll(urlS, "t.jpg", ".jpg")
			urlS = strings.ReplaceAll(urlS, "/t.", "/i.")
			pth := d + "/" + iutil.PadPgNum(i) + ".jpg"
			go mbpp.CreateDownloadJob(urlS, pth, mbpp.BlankWaitGroup(), b)
		}

		return func(bar *mbpp.BarProxy, _ *sync.WaitGroup) {

			d := iutil.GetDoc("https://" + host + path)
			t := d.Find("div#info h1").Text()

			dir := outputDir + "/" + iutil.FixTitleForFilename(t)
			out := dir + ".cbz"
			if util.DoesFileExist(out) {
				return
			}
			os.MkdirAll(dir, os.ModePerm)

			s := d.Find("div#thumbnail-container img[is=lazyload-image]")
			bar.AddToTotal(int64(s.Length()))
			s.Each(func(i int, el *goquery.Selection) {
				saveImage(i+1, el, dir, bar)
			})

			bar.Wait()
			bar.AddToTotal(1)
			iutil.PackCbzArchive(dir, out, bar)
		}
	}}
}
