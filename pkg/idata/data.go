package idata

import (
	"sync"

	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"
)

var (
	Hosts   = map[string]itypes.HostVal{}
	KeepJpg bool
	Wg      *sync.WaitGroup
	C       int
)
