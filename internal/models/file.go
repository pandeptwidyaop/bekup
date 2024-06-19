package models

import "github.com/pandeptwidyaop/bekup/internal/config"

type BackupFileInfo struct {
	Driver       string
	DatabaseName string
	FileName     string
	TempPath     string
	ZipName      string
	ZipPath      string
	ZipPassword  string
	Config       config.ConfigSource
	Err          error
}
