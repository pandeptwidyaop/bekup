package cleanup

import (
	"context"
	"os"

	"github.com/pandeptwidyaop/bekup/internal/log"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func Run(ctx context.Context, in <-chan models.BackupFileInfo) <-chan models.BackupFileInfo {
	return Cleanup(ctx, in)
}

func Cleanup(ctx context.Context, in <-chan models.BackupFileInfo) <-chan models.BackupFileInfo {
	out := make(chan models.BackupFileInfo)

	go func() {
		defer close(out)

		for file := range in {
			select {
			case out <- doCleanup(file):
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func doCleanup(f models.BackupFileInfo) models.BackupFileInfo {
	if f.TempPath != "" {
		log.GetInstance().Info("cleanup: removing ", f.TempPath)

		err := os.Remove(f.TempPath)
		if err != nil {
			f.Err = err
		}
	}

	if f.ZipPath != "" {
		log.GetInstance().Info("cleanup: removing ", f.ZipPath)

		err := os.Remove(f.ZipPath)
		if err != nil {
			f.Err = err
		}
	}

	return f
}
