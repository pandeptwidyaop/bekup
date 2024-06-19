package dump

import (
	"context"
	"sync"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/exception"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func Run(ctx context.Context, worker int, config config.Config, sources ...config.ConfigSource) (<-chan *models.BackupFileInfo, error) {
	return databaseManager(ctx, worker, config, sources...)
}

func databaseManager(ctx context.Context, worker int, config config.Config, sources ...config.ConfigSource) (<-chan *models.BackupFileInfo, error) {

	var chans []<-chan *models.BackupFileInfo

	for _, source := range sources {
		switch source.Driver {
		case "mysql":
			chans = append(chans, MysqlRun(ctx, config, source, worker))
		case "postgres":
			chans = append(chans, PostgresRun(ctx, config, source, worker))
		case "redis", "redis-clusters":
			chans = append(chans, RedisRun(ctx, config, source, worker))
		case "mongodb":
			chans = append(chans, MongoRun(ctx, config, source, worker))
		default:
			return nil, exception.ErrConfigSourceDriverNotAvailable
		}
	}

	return mergeChannel(chans), nil
}

func mergeChannel(chans []<-chan *models.BackupFileInfo) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)
	wg := sync.WaitGroup{}

	wg.Add(len(chans))

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ch := range chans {
		go func(c <-chan *models.BackupFileInfo) {
			for c := range ch {
				out <- c
			}

			wg.Done()
		}(ch)
	}

	return out
}
