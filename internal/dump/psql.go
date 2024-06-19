package dump

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/log"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func PostgresRun(ctx context.Context, config config.Config, source config.ConfigSource, worker int) <-chan *models.BackupFileInfo {
	ch := postgresRegister(ctx, config, source)

	return postgresBackupWithWorker(ctx, ch, worker)
}

func postgresRegister(ctx context.Context, config config.Config, source config.ConfigSource) <-chan *models.BackupFileInfo {
	log.GetInstance().Info("postgres: preparing backup")

	out := make(chan *models.BackupFileInfo)

	go func() {
		defer close(out)

		for _, db := range source.Databases {

			select {
			case <-ctx.Done():
				return
			default:
				id := uuid.New().String()

				fileName := fmt.Sprintf("postgres-%s-%s-%s.sql", time.Now().Format("2006-01-02-15-04-05-00"), db, id)

				log.GetInstance().Info("postgres: registering db ", db)

				out <- &models.BackupFileInfo{
					Driver:       source.Driver,
					DatabaseName: db,
					FileName:     fileName,
					Config:       source,
					TempPath:     config.TempPath,
					ZipPassword:  config.ZipPassword,
				}
			}

		}
	}()

	return out
}

func postgresBackupWithWorker(ctx context.Context, in <-chan *models.BackupFileInfo, worker int) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	wg := sync.WaitGroup{}

	var chans []<-chan *models.BackupFileInfo

	wg.Add(worker)

	for i := 0; i < worker; i++ {
		chans = append(chans, postgresBackup(ctx, in))
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ch := range chans {
		go func(c <-chan *models.BackupFileInfo) {
			for cc := range c {
				out <- cc
			}

			wg.Done()
		}(ch)
	}

	return out
}

func postgresBackup(ctx context.Context, in <-chan *models.BackupFileInfo) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
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

				out <- postgresDoBackup(info)
			}
		}
	}()

	return out
}

func postgresDoBackup(f *models.BackupFileInfo) *models.BackupFileInfo {
	log.GetInstance().Info("postgres: processing ", f.FileName)

	var stderr bytes.Buffer

	f.TempPath = path.Join(f.TempPath, f.FileName)

	fmt.Println(f.Config.Host)

	command := exec.Command("pg_dump", "-h", f.Config.Host, "-p", f.Config.Port, "-U", f.Config.Username, "-d", f.DatabaseName, "-f", f.TempPath)

	command.Stderr = &stderr
	command.Env = append(command.Env, "PGPASSWORD="+f.Config.Password)

	err := command.Run()
	if err != nil {
		f.Err = errors.New(stderr.String())
		return f
	}

	log.GetInstance().Info("postgres: done processing ", f.FileName)

	return f
}
