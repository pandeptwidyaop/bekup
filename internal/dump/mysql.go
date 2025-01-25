package dump

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/log"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func MysqlRun(ctx context.Context, config config.Config, source config.ConfigSource, worker int) <-chan *models.BackupFileInfo {
	ch := mysqlRegister(ctx, config, source)

	return mysqlBackupWithWorker(ctx, ch, worker)
}

func mysqlRegister(ctx context.Context, config config.Config, source config.ConfigSource) <-chan *models.BackupFileInfo {
	log.GetInstance().Info("mysql: preparing backup")
	channel := make(chan *models.BackupFileInfo)

	go func() {
		defer close(channel)

		for _, db := range source.Databases {
			id := uuid.New().String()
			fileName := fmt.Sprintf("mysql-%s-%s-%s.sql", time.Now().Format("2006-01-02-15-04-05-00"), db, id)
			log.GetInstance().Info("mysql: registering db ", db)
			select {
			case channel <- &models.BackupFileInfo{
				Driver:       source.Driver,
				DatabaseName: db,
				FileName:     fileName,
				Config:       source,
				TempPath:     config.TempPath,
				ZipPassword:  config.ZipPassword,
			}:
			case <-ctx.Done():
				return
			}

		}

	}()

	return channel
}

func mysqlBackupWithWorker(ctx context.Context, in <-chan *models.BackupFileInfo, worker int) <-chan *models.BackupFileInfo {
	wg := sync.WaitGroup{}

	out := make(chan *models.BackupFileInfo)
	var ins []<-chan *models.BackupFileInfo

	wg.Add(worker)

	for i := 0; i < worker; i++ {
		ins = append(ins, mysqlBackup(ctx, in))
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ch := range ins {
		go func(c <-chan *models.BackupFileInfo) {
			for cc := range c {
				out <- cc
			}

			wg.Done()
		}(ch)
	}

	return out
}

func mysqlBackup(ctx context.Context, in <-chan *models.BackupFileInfo) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	go func() {
		defer close(out)

		for {
			select {
			case info, ok := <-in:
				if !ok {
					return
				}

				if info == nil {
					continue
				}

				if info.Err != nil {
					out <- info
					continue
				}

				out <- mysqlDoBackup(info)
			case <-ctx.Done():
				return
			}
		}

	}()

	return out
}

func mysqlDoBackup(f *models.BackupFileInfo) *models.BackupFileInfo {
	log.GetInstance().Info("mysql: Processing ", f.FileName)

	var stderr bytes.Buffer
	f.TempPath = path.Join(f.TempPath, f.FileName)

	file, err := os.Create(f.TempPath)
	if err != nil {
		f.Err = err
		return f
	}

	defer file.Close()

	command := exec.Command("mysqldump", "-h", f.Config.Host, "-P", f.Config.Port, "-u", f.Config.Username, "-p"+f.Config.Password, f.DatabaseName)

	command.Stdout = file
	command.Stderr = &stderr

	err = command.Run()
	if err != nil {
		// fmt.Println(err)
		f.Err = errors.New(stderr.String())
		return f
	}

	log.GetInstance().Info("mysql: done processing ", f.FileName)

	return f
}
