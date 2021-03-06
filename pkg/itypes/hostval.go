package itypes

import (
	"github.com/nektro/go-util/mbpp"
)

type HostVal struct {
	IDPathIndex  int
	DownloadFunc func(string, string, string, string) func(*mbpp.BarProxy)
}
