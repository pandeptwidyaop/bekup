package zip

import (
	Z "archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pandeptwidyaop/bekup/internal/log"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func Run(ctx context.Context, in <-chan *models.BackupFileInfo, worker int) <-chan *models.BackupFileInfo {
	return ZipWithWorker(ctx, in, worker)
}

func Zip(ctx context.Context, in <-chan *models.BackupFileInfo) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	go func() {
		defer close(out)

		for {
			select {
			case info, ok := <-in:
				if !ok {
					return
				}

				if info != nil {
					continue
				}

				if info.Err != nil {
					out <- info
					continue
				}

				out <- doZip(info)
			case <-ctx.Done():
				return
			}

		}
	}()

	return out
}

func ZipWithWorker(ctx context.Context, in <-chan *models.BackupFileInfo, worker int) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)
	var ins []<-chan *models.BackupFileInfo

	wg := sync.WaitGroup{}
	wg.Add(worker)

	for i := 0; i < worker; i++ {
		ins = append(ins, Zip(ctx, in))
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ch := range ins {
		go func(c <-chan *models.BackupFileInfo) {
			for cc := range c {
				out <- cc
			}

			wg.Done()
		}(ch)

	}

	return out
}

func doZip(f *models.BackupFileInfo) *models.BackupFileInfo {
	log.GetInstance().Info("zip: zipping ", f.TempPath)
	f.ZipPath = fmt.Sprintf("%s.zip", f.TempPath)
	f.ZipName = fmt.Sprintf("%s.zip", f.FileName)

	file, err := os.Create(f.ZipPath)
	if err != nil {
		f.Err = err
		return f
	}
	defer file.Close()

	zw := Z.NewWriter(file)
	defer zw.Close()

	//TODO: Must check the f.TempPath is file or just a directory
	fileToZip, err := os.Open(f.TempPath)
	if err != nil {
		f.Err = err
		return f
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		f.Err = err
		return f
	}

	head, err := Z.FileInfoHeader(info)
	if err != nil {
		f.Err = err
		return f
	}

	head.Method = Z.Deflate

	wr, err := zw.CreateHeader(head)
	if err != nil {
		f.Err = err
		return f
	}

	_, err = io.Copy(wr, fileToZip)
	if err != nil {
		f.Err = err
		return f
	}

	log.GetInstance().Info("zip: success zip", f.TempPath, " to ", f.ZipPath)

	return f
}
