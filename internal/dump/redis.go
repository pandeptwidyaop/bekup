package dump

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/log"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func RedisRun(ctx context.Context, config config.Config, source config.ConfigSource, worker int) <-chan *models.BackupFileInfo {
	ch := redisRegister(ctx, config, source)

	return redisBackupWithWorker(ctx, ch, worker)
}

func redisRegister(ctx context.Context, config config.Config, source config.ConfigSource) <-chan *models.BackupFileInfo {
	log.GetInstance().Info("redis: preparing backup")

	out := make(chan *models.BackupFileInfo)

	go func() {
		defer close(out)
		db := "all"
		select {
		case <-ctx.Done():
			return
		default:
			id := uuid.New().String()

			fileName := fmt.Sprintf("redis-%s-%s-%s.rdb", time.Now().Format("2006-01-02-15-04-05-00"), db, id)
			if source.Driver == "redis-clusters" {
				if len(source.Databases) > 0 {
					db = source.Databases[0]
				} else {
					db = "slot1"
				}
				fileName = fmt.Sprintf("redis-cluster-%s-%s-%s.rdb", time.Now().Format("2006-01-02-15-04-05-00"), db, id)
			}
			log.GetInstance().Info("redis: registering db ", db)

			out <- &models.BackupFileInfo{
				Driver:       source.Driver,
				DatabaseName: db,
				FileName:     fileName,
				Config:       source,
				TempPath:     config.TempPath,
				ZipPassword:  config.ZipPassword,
			}
		}

	}()

	return out
}

func redisBackupWithWorker(ctx context.Context, in <-chan *models.BackupFileInfo, worker int) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	wg := sync.WaitGroup{}

	var chans []<-chan *models.BackupFileInfo

	wg.Add(worker)

	for i := 0; i < worker; i++ {
		chans = append(chans, redisBackup(ctx, in))
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

func redisBackup(ctx context.Context, in <-chan *models.BackupFileInfo) <-chan *models.BackupFileInfo {
	out := make(chan *models.BackupFileInfo)

	go func() {
		defer close(out)

		for {
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

				out <- redisDoBackup(info)
			}
		}
	}()

	return out
}

func redisDoBackup(f *models.BackupFileInfo) *models.BackupFileInfo {
	log.GetInstance().Info("redis: processing ", f.FileName)

	var stdErrC1 bytes.Buffer
	var stdErrC2 bytes.Buffer
	var stdErrC3 bytes.Buffer

	f.TempPath = path.Join(f.TempPath, f.FileName)

	fmt.Println(f.Config.Host)

	fmt.Println("initiating background saving")
	// Initiate BG Save
	c1 := exec.Command("redis-cli", "-h", f.Config.Host, "-p", f.Config.Port, "-a", f.Config.Password, "BGSAVE")

	var stdOutC1 bytes.Buffer
	c1.Stdout = &stdOutC1
	c1.Stderr = &stdErrC1

	err := c1.Run()
	if err != nil {
		f.Err = errors.New(stdErrC1.String())
		return f
	}

	if !strings.Contains(strings.ToLower(stdOutC1.String()), "background saving started") &&
		!strings.Contains(strings.ToLower(stdOutC1.String()), "already in progress") {
		f.Err = errors.New("failed to execute background saving, either not allowed or already running")
		return f
	}

	fmt.Println("waiting background saving to complete")
	// Wait BG Save Progress
	for {
		c2 := exec.Command("redis-cli", "-h", f.Config.Host, "-p", f.Config.Port, "-a", f.Config.Password, "INFO", "persistence")
		var outC2 bytes.Buffer
		c2.Stdout = &outC2
		c2.Stderr = &stdErrC2

		err = c2.Run()
		if err != nil {
			f.Err = errors.New(stdErrC2.String())
			return f
		}

		if strings.Contains(outC2.String(), "bgsave_in_progress:1") {
			time.Sleep(60 * time.Second)
			continue
		} else if strings.Contains(outC2.String(), "bgsave_in_progress:0") {
			break
		} else {
			f.Err = errors.New("unexpected output from Redis INFO persistence")
			return f
		}
	}

	fmt.Println("exporting dump file")
	// Export RDB File
	c3 := exec.Command("redis-cli", "-h", f.Config.Host, "-p", f.Config.Port, "-a", f.Config.Password, "--rdb", f.TempPath)
	c3.Stderr = &stdErrC3

	err = c3.Run()
	if err != nil {
		f.Err = errors.New(stdErrC3.String())
		return f
	}

	log.GetInstance().Info("redis: done processing ", f.FileName)

	return f
}
