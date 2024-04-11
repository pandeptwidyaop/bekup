package tests

import (
	"context"
	"testing"

	"github.com/pandeptwidyaop/bekup/internal/bekup"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {

	ctx := context.Background()

	cf := config.Config{
		Sources: []config.ConfigSource{
			{
				Driver:     "mongodb",
				MongoDBURI: "mongodb://test:test@mongo:27017?authSource=admin",
				Databases:  []string{"test"},
			},
		},
	}

	err := bekup.Run(ctx, cf, 2)
	if err != nil {
		assert.Nil(t, err)
	}
}
