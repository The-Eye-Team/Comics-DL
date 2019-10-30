package idata

import (
	"os"

	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"

	"golang.org/x/sync/semaphore"
)

var (
	Hosts   = map[string]itypes.HostVal{}
	KeepJpg bool
	Guard   *semaphore.Weighted
	Log     *os.File
)
