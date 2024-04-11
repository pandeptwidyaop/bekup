package upload

import (
	"context"
	"sync"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/exception"
	"github.com/pandeptwidyaop/bekup/internal/log"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func Run(ctx context.Context, in <-chan *models.BackupFileInfo, worker int, destinations ...config.ConfigDestination) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	wg := sync.WaitGroup{}

	var lists []<-chan *models.BackupFileInfo

	wg.Add(worker)

	for i := 0; i < worker; i++ {
		lists = append(lists, run(ctx, in, destinations...))
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ls := range lists {
		go func(c <-chan *models.BackupFileInfo) {
			for cc := range ls {
				out <- cc
			}
			wg.Done()
		}(ls)
	}

	return out
}

func run(ctx context.Context, in <-chan *models.BackupFileInfo, destinations ...config.ConfigDestination) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	go func() {
		defer close(out)

		for {
			for _, destination := range destinations {
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

					out <- uploadManager(ctx, info, destination)
				}
			}
		}
	}()

	return out
}

func uploadManager(ctx context.Context, f *models.BackupFileInfo, d config.ConfigDestination) *models.BackupFileInfo {
	switch d.Driver {
	case "s3":
		return S3Upload(ctx, f, d)
	default:
		log.GetInstance().Error("your config driver ", d.Driver, " is not available yet")
		f.Err = exception.ErrConfigDestinationDriverNotAvailable
		return f
	}
}
