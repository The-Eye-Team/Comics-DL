package sites

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/The-Eye-Team/Comics-DL/pkg/idata"
	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"
	"github.com/The-Eye-Team/Comics-DL/pkg/iutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/nektro/go-util/mbpp"
)

func init() {
	idata.Hosts["e-hentai.org"] = itypes.HostVal{2, func(host string, id string, path string, outputDir string) func(*mbpp.BarProxy, *sync.WaitGroup) {

		savePage := func(p int, d *goquery.Document, b *mbpp.BarProxy, dir string) {
			s := d.Find("div.gdtm a img")
			b.AddToTotal(int64(s.Length()))

			s.Each(func(i int, el *goquery.Selection) {
				pn := iutil.PadPgNum(p) + "_" + iutil.PadPgNum(i)
				//
				url1, _ := el.Parent().Attr("href")
				url2, _ := iutil.GetDoc(url1).Find("#img").Attr("src")
				mbpp.CreateDownloadJob(url2, dir+"/"+pn+".jpg", mbpp.BlankWaitGroup(), b)
			})
		}

		return func(bar *mbpp.BarProxy, _ *sync.WaitGroup) {

			d := iutil.GetDoc("https://" + host + path + "?p=0")
			t := strings.TrimSuffix(d.Find("title").Text(), " - E-Hentai Galleries")
			bar.AddToTotal(1)

			dir := outputDir + "/" + id + "." + iutil.FixTitleForFilename(t)
			os.MkdirAll(dir, os.ModePerm)

			savePage(0, d, bar, dir)
			c := d.Find("table.ptt td")
			if c.Length() > 3 {
				pu, _ := c.Eq(c.Length() - 2).Children().Eq(0).Attr("href")
				pl, _ := strconv.ParseInt(strings.Split(pu, "=")[1], 10, 32)

				pli := int(pl)
				for i := 1; i <= pli; i++ {
					savePage(i, iutil.GetDoc("https://"+host+path+"?p="+strconv.Itoa(pli)), bar, dir)
				}
			}

			iutil.PackCbzArchive(dir, host+"/"+id, bar)
		}
	}}
}
