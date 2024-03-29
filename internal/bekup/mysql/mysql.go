package mysql

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/pandeptwidyaop/bekup/internal/bekup"
	"github.com/pandeptwidyaop/bekup/internal/config"
)

func Register(ctx context.Context, source config.ConfigSource) <-chan bekup.BackupFileInfo {
	channel := make(chan bekup.BackupFileInfo)

	go func() {
		defer close(channel)

		for _, db := range source.Databases {
			fileName := fmt.Sprintf("mysql-%s-%s.sql", time.Now().Format("2006-01-02-15-04-05-00"), db)

			select {
			case channel <- bekup.BackupFileInfo{
				FileName: fileName,
				Config:   source,
			}:
			case <-ctx.Done():
				return
			}

		}

	}()

	return channel
}

func BackupWithWorker(ctx context.Context, in <-chan bekup.BackupFileInfo, worker int) <-chan bekup.BackupFileInfo {
	wg := sync.WaitGroup{}

	out := make(chan bekup.BackupFileInfo)
	var ins []<-chan bekup.BackupFileInfo

	wg.Add(worker)

	for i := 0; i < worker; i++ {
		ins = append(ins, Backup(ctx, in))
	}

	for _, ch := range ins {
		go func(c <-chan bekup.BackupFileInfo) {
			for cc := range c {
				out <- cc
			}

			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Backup(ctx context.Context, in <-chan bekup.BackupFileInfo) <-chan bekup.BackupFileInfo {
	out := make(chan bekup.BackupFileInfo)

	go func() {
		defer close(out)

		for info := range in {
			select {
			case out <- doBackup(info):
			case <-ctx.Done():
				return
			}
		}

	}()

	return out
}

func doBackup(f bekup.BackupFileInfo) bekup.BackupFileInfo {

	f.TempPath = path.Join(config.GetTempPath(), f.FileName)

	file, err := os.Create(f.TempPath)
	if err != nil {
		f.Err = err
		return f
	}

	defer file.Close()

	command := exec.Command("mysqldump", "-h", f.Config.Host, "-P", f.Config.Port, "-u", f.Config.Username, "-p"+f.Config.Password, f.TempPath)

	command.Stdout = file

	err = command.Run()
	if err != nil {
		f.Err = err
		return f
	}

	return f
}
