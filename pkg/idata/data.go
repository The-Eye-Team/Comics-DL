package idata

import (
	"os"

	"github.com/The-Eye-Team/Comics-DL/pkg/itypes"
)

var (
	Hosts   = map[string]itypes.HostVal{}
	KeepJpg bool
	Log     *os.File
)
