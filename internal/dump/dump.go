package dump

import (
	"context"
	"sync"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/exception"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func Run(ctx context.Context, worker int, sources ...config.ConfigSource) (<-chan models.BackupFileInfo, error) {
	return databaseManager(ctx, worker, sources...)
}

func databaseManager(ctx context.Context, worker int, sources ...config.ConfigSource) (<-chan models.BackupFileInfo, error) {

	var chans []<-chan models.BackupFileInfo

	for _, source := range sources {
		switch source.Driver {
		case "mysql":
			chans = append(chans, MysqlRun(ctx, source, worker))
		case "postgres":
			chans = append(chans, PostgresRun(ctx, source, worker))
		default:
			return nil, exception.ErrConfigSourceDriverNotAvailable
		}
	}

	return mergeChannel(chans), nil
}

func mergeChannel(chans []<-chan models.BackupFileInfo) <-chan models.BackupFileInfo {
	out := make(chan models.BackupFileInfo)
	wg := sync.WaitGroup{}

	wg.Add(len(chans))

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ch := range chans {
		go func(c <-chan models.BackupFileInfo) {
			for c := range ch {
				out <- c
			}

			wg.Done()
		}(ch)
	}

	return out
}
