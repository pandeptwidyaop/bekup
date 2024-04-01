package bekup

import (
	"context"
	"sync"

	"github.com/pandeptwidyaop/bekup/internal/cleanup"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/exception"
	"github.com/pandeptwidyaop/bekup/internal/models"
	"github.com/pandeptwidyaop/bekup/internal/mysql"
	"github.com/pandeptwidyaop/bekup/internal/s3"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, config config.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	var chans []<-chan models.BackupFileInfo

	for _, source := range config.Sources {
		switch source.Driver {
		case "mysql":
			chans = append(chans, mysql.Run(ctx, source, 10))
		case "postgres":
		case "mongodb":
		default:
			cancel()
			return exception.ErrConfigSourceDriverNotAvailable
		}
	}

	backupChan := mergeChannel(chans)

	//Upload here
	var uploadChans []<-chan models.BackupFileInfo

	for _, dst := range config.Destinations {
		switch dst.Driver {
		case "s3":
		case "ftp":
		default:
			cancel()
			return exception.ErrConfigDestinationDriverNotAvailable
		}
	}

	for _, dst := range config.Destinations {
		switch dst.Driver {
		case "s3":
			uploadChans = append(uploadChans, s3.Run(ctx, backupChan, dst, 10))
		case "ftp":
	
		}
	}

	uploadCh := mergeChannel(uploadChans)

	cleanupCh := cleanup.Run(ctx, uploadCh)

	g.Go(func() error {
		for m := range cleanupCh {
			if m.Err != nil {
				return m.Err
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		cancel()
		return err
	}

	return nil
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

func mergeSequentialChannel(chans []<-chan models.BackupFileInfo) <-chan models.BackupFileInfo {
	out := make(chan models.BackupFileInfo)

	wg := sync.WaitGroup{}
	wg.Add(len(chans))

	go func() {
		wg.Wait()
		close(out)
	}()

	go func() {
		for _, ch := range chans {
			for c := range ch {
				out <- c
			}

			wg.Done()
		}
	}()

	return out
}
