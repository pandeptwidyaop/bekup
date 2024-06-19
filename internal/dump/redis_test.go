package dump_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/dump"
)

func Test_Redis(t *testing.T) {
	config := config.ConfigSource{
		Driver:   "redis",
		Password: "",
		Host:     "localhost",
		Port:     "6379",
		Databases: []string{
			"all",
		},
	}

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	ch := dump.RedisRun(ctx, config, 2)

	for c := range ch {
		if c.Err != nil {
			cancel()
			fmt.Println(c.Err.Error())
			return
		}
		fmt.Println(c)
	}
}
