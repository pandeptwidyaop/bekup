package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func Run(ctx context.Context, in <-chan models.BackupFileInfo, worker int) <-chan models.BackupFileInfo {
	out := make(chan models.BackupFileInfo)

	return out
}

func run(ctx context.Context, in <-chan models.BackupFileInfo) <-chan models.BackupFileInfo {
	out := make(chan models.BackupFileInfo)

	go func() {
		defer close(out)

		for f := range in {
			select {
			case out <- doUpload(f):
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func doUpload(f models.BackupFileInfo) models.BackupFileInfo {
	return f
}

func newSession(src config.ConfigDestination) (*session.Session, error) {
	session.NewSession()
}
