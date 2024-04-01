package models

import "github.com/pandeptwidyaop/bekup/internal/config"

type BackupFileInfo struct {
	DatabaseName string
	FileName     string
	TempPath     string
	ZipPath      string
	Config       config.ConfigSource
	Uploads      map[string]string
	Err          error
}


