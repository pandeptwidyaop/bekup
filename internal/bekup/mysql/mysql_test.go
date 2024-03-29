package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pandeptwidyaop/bekup/internal/bekup/mysql"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"golang.org/x/sync/errgroup"
)

func Test_run(t *testing.T) {
	source := config.ConfigSource{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "root",
		Password: "",
	}

	for i := 0; i <= 100; i++ {
		source.Databases = append(source.Databases, fmt.Sprintf("%d", i))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	ch := mysql.Register(ctx, source)

	ch1 := mysql.BackupWithWorker(ctx, ch, 10)

	g.Go(func() error {
		for f := range ch1 {
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

}
