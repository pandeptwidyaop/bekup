package cleanup

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/pandeptwidyaop/bekup/internal/log"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func Run(ctx context.Context, in <-chan *models.BackupFileInfo, worker int) <-chan *models.BackupFileInfo {
	return cleanupWithWorker(ctx, in, worker)
}

func cleanupWithWorker(ctx context.Context, in <-chan *models.BackupFileInfo, worker int) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	wg := sync.WaitGroup{}
	wg.Add(worker)

	var chans []<-chan *models.BackupFileInfo

	for i := 0; i < worker; i++ {
		chans = append(chans, cleanup(ctx, in))
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

func cleanup(ctx context.Context, in <-chan *models.BackupFileInfo) <-chan *models.BackupFileInfo {
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

				out <- doCleanup(info)
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func doCleanup(f *models.BackupFileInfo) *models.BackupFileInfo {
	if f.TempPath != "" {

		log.GetInstance().Info("cleanup: removing ", f.TempPath)
		if checkFileExist(f.TempPath) {
			err := os.RemoveAll(f.TempPath)
			if err != nil {
				f.Err = err
			}
		} else {
			log.GetInstance().Warn("file ", f.TempPath, " is not exist")
		}
	}

	if f.ZipPath != "" {
		log.GetInstance().Info("cleanup: removing ", f.ZipPath)

		if checkFileExist(f.ZipPath) {
			err := os.RemoveAll(f.ZipPath)
			if err != nil {
				f.Err = err
			}
		} else {
			log.GetInstance().Warn("file ", f.TempPath, " is not exist")
		}

	}

	return f
}

func checkFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func ForceCleanUp() ([]string, error) {
	all := []string{}

	files, err := os.ReadDir(".")
	if err != nil {
		return all, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.HasSuffix(f.Name(), ".sql") || strings.HasSuffix(f.Name(), ".zip") {
			err := os.Remove(f.Name())
			if err != nil {
				return all, err
			}

			all = append(all, f.Name())
		}
	}

	return all, nil
}
