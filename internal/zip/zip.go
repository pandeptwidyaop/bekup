package zip

import (
	Z "archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

				if info == nil {
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

	typ, err := checkPathType(f.TempPath)
	if err != nil {
		f.Err = err
		return f
	}	

	if typ == "file" {
		return doZipSingleFile(f)
	}

	return doZipDirectory(f)
}

func doZipSingleFile(f *models.BackupFileInfo) *models.BackupFileInfo {
	file, err := os.Create(f.ZipPath)
	if err != nil {
		f.Err = err
		return f
	}
	defer file.Close()

	zw := Z.NewWriter(file)
	defer zw.Close()

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

func doZipDirectory(f *models.BackupFileInfo) *models.BackupFileInfo {
	zipfile, err := os.Create(f.ZipPath)
	if err != nil {
		f.Err = err
		return f
	}
	defer zipfile.Close()

	archive := Z.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(f.TempPath)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(f.TempPath)
	}

	err = filepath.Walk(f.TempPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := Z.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, f.TempPath))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = Z.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return err
	})

	if err != nil {
		f.Err = err
		return f
	}

	return f
}

func checkPathType(path string) (string, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", errors.New("path does not exists")
	}

	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "dir", nil
	} else {
		return "file", nil
	}
}
