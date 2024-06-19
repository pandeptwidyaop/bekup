package dump_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/dump"
	"github.com/stretchr/testify/assert"
)

func TestMongo(t *testing.T) {
	conf := config.ConfigSource{
		Driver:     "mongodb",
		MongoDBURI: "mongodb://test:test@mongo:27017",
		Databases: []string{
			"test",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := dump.MongoRun(ctx, conf, 2)

	for c := range ch {
		if c.Err != nil {
			assert.Nil(t, c.Err)
			cancel()
		}

		fmt.Println(c)
	}
}
