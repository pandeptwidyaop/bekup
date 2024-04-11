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

func MongoRun(ctx context.Context, source config.ConfigSource, worker int) <-chan *models.BackupFileInfo {
	ch := mongoRegister(ctx, source)

	return mongoBackupWithWorker(ctx, ch, worker)
}

func mongoRegister(ctx context.Context, source config.ConfigSource) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	go func() {
		defer close(out)

		for _, db := range source.Databases {

			select {
			case <-ctx.Done():
				return
			default:
				id := uuid.New().String()

				fileName := fmt.Sprintf("mongo-%s-%s-%s", time.Now().Format("2006-01-02-15-04-05-00"), db, id)

				log.GetInstance().Info("mongo: registering db ", db)

				out <- &models.BackupFileInfo{
					DatabaseName: db,
					FileName:     fileName,
					Config:       source,
				}
			}
		}
	}()

	return out
}

func mongoBackupWithWorker(ctx context.Context, in <-chan *models.BackupFileInfo, worker int) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	var chans []<-chan *models.BackupFileInfo

	wg := sync.WaitGroup{}

	wg.Add(worker)

	for i := 0; i < worker; i++ {
		chans = append(chans, mongoBackup(ctx, in))
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

func mongoBackup(ctx context.Context, in <-chan *models.BackupFileInfo) <-chan *models.BackupFileInfo {
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

				out <- mongoDoBackup(info)

			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func mongoDoBackup(f *models.BackupFileInfo) *models.BackupFileInfo {
	log.GetInstance().Info("mongo: processing ", f.FileName)

	var stderr bytes.Buffer
	f.TempPath = path.Join(f.TempPath, f.FileName)

	err := os.MkdirAll(f.TempPath, 0775)
	if err != nil {
		f.Err = err
		return f
	}

	var uri string

	if f.Config.MongoDBURI != "" {
		uri = f.Config.MongoDBURI
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", f.Config.Username, f.Config.Password, f.Config.Host, f.Config.Port)
	}

	command := exec.Command("mongodump", "--uri", uri, "--db", f.DatabaseName, "--out", f.TempPath)

	command.Stderr = &stderr

	err = command.Run()
	if err != nil {
		f.Err = errors.New(stderr.String())
		return f
	}

	return f
}
