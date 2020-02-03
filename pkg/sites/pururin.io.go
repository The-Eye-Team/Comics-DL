package sites

import (
	"os"
	"strconv"

	"github.com/The-Eye-Team/Comics-DL/pkg/idata"
	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"
	"github.com/The-Eye-Team/Comics-DL/pkg/iutil"

	"github.com/nektro/go-util/mbpp"
	"github.com/nektro/go-util/util"
)

func init() {
	idata.Hosts["pururin.io"] = itypes.HostVal{2, func(host string, id string, path string, outputDir string) func(*mbpp.BarProxy) {
		return func(bar *mbpp.BarProxy) {

			d := iutil.GetDoc("https://" + host + path)
			t := d.Find("div.content-wrapper div.title h1").Text()

			lr, _ := d.Find("gallery-thumbnails").Attr(":total")
			l := iutil.ParseInt(lr)

			dir := outputDir + "/" + iutil.FixTitleForFilename(t)
			out := dir + ".cbz"
			if util.DoesFileExist(out) {
				return
			}
			os.MkdirAll(dir, os.ModePerm)

			bar.AddToTotal(int64(l))
			for i := 1; i <= l; i++ {
				f := strconv.Itoa(i) + ".jpg"
				g := iutil.PadPgNum(i) + ".jpg"
				go mbpp.CreateDownloadJob("https://cdn.pururin.io/assets/images/data/"+id+"/"+f, dir+"/"+g, bar)
			}

			bar.Wait()
			iutil.PackCbzArchive(dir, host+"/"+id, bar)
		}
	}}
}
