package dump_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/dump"
)

func Test_psql(t *testing.T) {
	config := config.ConfigSource{
		Driver:   "postgres",
		Username: "test",
		Password: "test",
		Host:     "postgres",
		Port:     "5432",
		Databases: []string{
			"test",
		},
	}

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	ch := dump.PostgresRun(ctx, config, 2)

	for c := range ch {
		if c.Err != nil {
			cancel()
			fmt.Println(c.Err.Error())
			return
		}
		fmt.Println(c)
	}
}
