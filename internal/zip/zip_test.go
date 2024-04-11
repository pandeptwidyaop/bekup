package zip_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pandeptwidyaop/bekup/internal/models"
	"github.com/pandeptwidyaop/bekup/internal/zip"
)

func TestZipDir(t *testing.T) {

	ch := make(chan *models.BackupFileInfo)

	go func() {
		cout := zip.Run(context.Background(), ch, 1)

		fmt.Println(<-cout)
	}()

	ch <- &models.BackupFileInfo{
		TempPath: "/home/devops/tests",
	}

	close(ch)

	time.Sleep(15 * time.Second)

}
