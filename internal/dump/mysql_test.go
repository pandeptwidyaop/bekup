package dump_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/dump"
	"github.com/pandeptwidyaop/bekup/internal/zip"
	"golang.org/x/sync/errgroup"
)

func Test_run(t *testing.T) {
	source := config.ConfigSource{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "33061",
		Username: "root",
		Password: "root",
	}

	for i := 0; i < 1000; i++ {
		source.Databases = append(source.Databases, "classicmodels")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	startAt := time.Now()

	chInit := dump.MysqlRun(ctx, source, 10)

	chZip := zip.ZipWithWorker(ctx, chInit, 10)

	g.Go(func() error {
		for f := range chZip {
			if f.Err != nil {
				return f.Err
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
		cancel()
	}

	fmt.Println(time.Since(startAt))

}
