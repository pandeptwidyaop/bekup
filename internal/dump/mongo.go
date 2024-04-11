package dump

import (
	"context"
	"fmt"
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

				fileName := fmt.Sprintf("mongo-%s-%s-%s", time.Now(), db, id)

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
					return
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

	return f
}
