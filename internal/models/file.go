package models

import "github.com/pandeptwidyaop/bekup/internal/config"

type BackupFileInfo struct {
	DatabaseName string
	FileName     string
	TempPath     string
	ZipName      string
	ZipPath      string
	Config       config.ConfigSource
	Err          error
}
