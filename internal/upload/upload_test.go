package upload_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/models"
	"github.com/pandeptwidyaop/bekup/internal/upload"
	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan *models.BackupFileInfo)

	go func() {
		defer close(ch)

		for i := 0; i < 100; i++ {
			ch <- &models.BackupFileInfo{
				ZipPath:      "/Users/pande/Projects/bekup/tests/temp/zip.zip",
				ZipName:      fmt.Sprintf("mysql-2024-04-01-10-00-00-%d.sql.zip", i),
				DatabaseName: fmt.Sprintf("database-%d", i),
			}
		}
	}()

	cf := []config.ConfigDestination{
		{
			AWSRegion:     os.Getenv("AWS_REGION"),
			AWSAccessKey:  os.Getenv("AWS_ACCESS_KEY"),
			AWSSecretKey:  os.Getenv("AWS_SECRET_KEY"),
			AWSBucket:     os.Getenv("AWS_BUCKET"),
			AWSUrl:        os.Getenv("AWS_URL"),
			RootDirectory: "backup1",
			Driver:        "s3",
		},
		{
			AWSRegion:     os.Getenv("AWS_REGION"),
			AWSAccessKey:  os.Getenv("AWS_ACCESS_KEY"),
			AWSSecretKey:  os.Getenv("AWS_SECRET_KEY"),
			AWSBucket:     os.Getenv("AWS_BUCKET"),
			AWSUrl:        os.Getenv("AWS_URL"),
			RootDirectory: "backup2",
			Driver:        "s3",
		},
	}

	ch2 := upload.Run(ctx, ch, 10, cf...)

	var dt []*models.BackupFileInfo

	for m := range ch2 {
		dt = append(dt, m)
		assert.Nil(t, m.Err)
		if m.Err != nil {
			cancel()
		}
	}

	fmt.Println("PANJANG DATA:", len(dt))

}
