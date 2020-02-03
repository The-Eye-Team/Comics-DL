package sites

import (
	"os"
	"strconv"
	"strings"

	"github.com/The-Eye-Team/Comics-DL/pkg/idata"
	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"
	"github.com/The-Eye-Team/Comics-DL/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

func init() {
	idata.Hosts["myreadingmanga.info"] = itypes.HostVal{1, func(host string, id string, path string, outputDir string) func(*mbpp.BarProxy) {

		savePage := func(p int, d *goquery.Document, b *mbpp.BarProxy, dir string) {
			s := d.Find("div.entry-content img.img-myreadingmanga")
			b.AddToTotal(int64(s.Length()))

			s.Each(func(i int, el *goquery.Selection) {
				pn := iutil.PadPgNum(p) + "_" + iutil.PadPgNum(i)
				//
				urlS, _ := el.Attr("data-lazy-src")
				mbpp.CreateDownloadJob(urlS, dir+"/"+pn+".jpg", b)
			})
		}

		return func(bar *mbpp.BarProxy) {

			d := iutil.GetDoc("https://" + host + "/" + id + "/")
			t := strings.TrimSuffix(d.Find("title").Text(), " - MyReadingManga")

			dir := outputDir + "/" + iutil.FixTitleForFilename(t)
			out := dir + ".cbz"
			if util.DoesFileExist(out) {
				return
			}
			os.MkdirAll(dir, os.ModePerm)

			savePage(1, d, bar, dir)
			c := d.Find("a.post-page-numbers")
			if c.Length() > 0 {
				end := iutil.ParseInt(c.Eq(c.Length() - 2).Text())
				for i := 2; i < end; i++ {
					is := strconv.Itoa(i)
					savePage(i, iutil.GetDoc("https://"+host+"/"+id+"/"+is+"/"), bar, dir)
				}
			}

			bar.Wait()
			iutil.PackCbzArchive(dir, host+"/"+id, bar)
		}
	}}
}
