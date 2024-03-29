package config_test

import (
	"bytes"
	"testing"

	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/exception"
	"github.com/stretchr/testify/assert"
)

func Test_LoadConf(t *testing.T) {
	type testCase struct {
		name     string
		conf     string
		isError  bool
		expected error
	}

	cases := []testCase{
		{
			name:     "check is config json is valid",
			conf:     ``,
			isError:  true,
			expected: exception.ErrConfigNotValid,
		},
		{
			name:     "check is config sources not valid array",
			conf:     `{"sources": ""}`,
			isError:  true,
			expected: exception.ErrConfigNotValid,
		},
		{
			name:     "check is config sources valid array but lengt is < 1",
			conf:     `{"sources": []}`,
			isError:  true,
			expected: exception.ErrConfigSourceNotExist,
		},
		{
			name:     "check is config has destinations",
			conf:     `{"sources": [{}], "destinations": ""}`,
			isError:  true,
			expected: exception.ErrConfigNotValid,
		},
		{
			name:     "check driver is not available",
			conf:     `{"sources": [{}], "destinations": [{"driver": ""}]}`,
			isError:  true,
			expected: exception.ErrConfigSourceDriverNotAvailable,
		},
		{
			name:     "check if destination driver is not available yet",
			conf:     `{"sources": [{"driver": "mysql"}], "destinations": [{"driver": "sftp"}]}`,
			isError:  true,
			expected: exception.ErrConfigDestinationDriverNotAvailable,
		},
		{
			name:     "check mysql driver config error",
			conf:     `{"sources": [{"driver": "mysql"}], "destinations": [{"driver": "ftp"}]}`,
			isError:  true,
			expected: exception.ErrConfigSourceError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cf := bytes.NewReader([]byte(c.conf))

			_, err := config.LoadConfig(cf)

			if c.isError {
				assert.EqualError(t, err, c.expected.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}

}
