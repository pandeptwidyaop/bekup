package bekup

import (
	"github.com/pandeptwidyaop/bekup/internal/config"
)

type BackupFileInfo struct {
	FileName string
	TempPath string
	ZipPath  string
	Config   config.ConfigSource
	Err      error
}
