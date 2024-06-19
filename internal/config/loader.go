package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/pandeptwidyaop/bekup/internal/exception"
	"github.com/pandeptwidyaop/bekup/internal/log"
)

var (
	AVAILABLE_SOURCE_DRIVERS      []string = []string{"mysql", "postgres", "mongodb", "redis", "redis-clusters"}
	AVAILABLE_DESTINATION_DRIVERS []string = []string{"s3", "ftp"}
)

type ConfigSource struct {
	Driver     string   `json:"driver" validate:"required"`
	Host       string   `json:"host"`
	Port       string   `json:"port"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	Databases  []string `json:"databases"`
	MongoDBURI string   `json:"mongodb_uri"`
}

type ConfigDestination struct {
	Driver        string `json:"driver"`
	AWSAccessKey  string `json:"aws_access_key"`
	AWSSecretKey  string `json:"aws_secret_access_key"`
	AWSBucket     string `json:"aws_bucket"`
	AWSRegion     string `json:"aws_region"`
	AWSUrl        string `json:"aws_url"`
	RootDirectory string `json:"root_directory"`
	Host          string `json:"host"`
	FTPPort       string `json:"port"`
	FTPUsername   string `json:"username"`
	FTPPassword   string `json:"password"`
}

type Config struct {
	Sources      []ConfigSource      `json:"sources"`
	Destinations []ConfigDestination `json:"destinations"`
	ZipPassword  string              `json:"zip_password"`
	Worker       int                 `json:"worker"`
	TempPath     string              `json:"temp_path"`
}

func LoadConfigFromPath(path string) (Config, error) {
	conf := Config{}

	if !isFileExist(path) {
		return conf, exception.ErrFileNotExists
	}

	file, err := os.Open(path)
	if err != nil {
		return conf, err
	}

	defer file.Close()

	return LoadConfig(file)
}

func LoadConfig(file io.Reader) (Config, error) {
	conf := Config{}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.GetInstance().Error(err)
		return conf, err
	}

	err = json.Unmarshal(jsonData, &conf)
	if err != nil {
		log.GetInstance().Error(err)
		return conf, exception.ErrConfigNotValid
	}

	if !isConfigHasSource(conf) {
		return conf, exception.ErrConfigSourceNotExist
	}

	if !isConfigHasDestinations(conf) {
		return conf, exception.ErrConfigDestinationNotExist
	}

	if err = checkSourcesDriver(conf); err != nil {
		log.GetInstance().Error("driver: error defined source 'driver' is not exist, for now only support: ", strings.Join(AVAILABLE_SOURCE_DRIVERS, ","))
		return conf, err
	}

	if err = checkDestinationDriver(conf); err != nil {
		log.GetInstance().Error("driver: error defined destination 'driver' is not exist, for now only support: ", strings.Join(AVAILABLE_DESTINATION_DRIVERS, ","))
		return conf, err
	}

	err = validateSourceConfig(conf)
	if err != nil {
		log.GetInstance().Error(err)
		return conf, exception.ErrConfigSourceError
	}

	if conf.TempPath == "" {
		return conf, exception.ErrConfigTempPathNotExists
	}

	if conf.Worker == 0 {
		conf.Worker = 2
	}

	return conf, nil
}

func isFileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func isConfigHasSource(conf Config) bool {
	return len(conf.Sources) > 0
}

func isConfigHasDestinations(conf Config) bool {
	return len(conf.Destinations) > 0
}

func checkSourcesDriver(conf Config) error {
	count := 0
	for _, c := range conf.Sources {
		for _, d := range AVAILABLE_SOURCE_DRIVERS {
			if c.Driver == d {
				count++
			}
		}
	}

	if count != len(conf.Sources) {
		return exception.ErrConfigSourceDriverNotAvailable
	}

	return nil
}

func checkDestinationDriver(conf Config) error {
	// count := 0
	isAvail := false
	for _, c := range conf.Destinations {
		for _, d := range AVAILABLE_DESTINATION_DRIVERS {
			if c.Driver == d {
				isAvail = true
				// 			count++
			}
		}
	}

	if !isAvail {

		return exception.ErrConfigDestinationDriverNotAvailable
	}

	// if count != len(conf.Sources) {
	// 	return exception.ErrConfigDestinationDriverNotAvailable
	// }

	return nil
}

func validateSourceConfig(conf Config) error {
	for _, sc := range conf.Sources {
		switch sc.Driver {
		case "mysql":
			return checkSourceMysqlDriver(sc)
		case "postgres":
			return checkSourcePostgresDriver(sc)
		case "redis", "redis-clusters":
			return checkSourceRedisDriver(sc)
		case "mongodb":
			return checkSourceMongodbDriver(sc)
		}
	}

	return nil
}

func checkSourceMysqlDriver(source ConfigSource) error {
	msg := []string{}

	if source.Host == "" {
		msg = append(msg, "host")
	}

	if source.Port == "" {
		msg = append(msg, "port")
	}

	if source.Username == "" {
		msg = append(msg, "username")
	}

	if source.Password == "" {
		msg = append(msg, "password")
	}

	for _, db := range source.Databases {
		if db == "" {
			msg = append(msg, "database name")
		}
	}

	if len(msg) > 0 {
		return errors.New("config mysql: some field empty: " + strings.Join(msg, ","))
	}

	return nil
}

func checkSourcePostgresDriver(source ConfigSource) error {
	msg := []string{}

	if source.Host == "" {
		msg = append(msg, "host")
	}

	if source.Port == "" {
		msg = append(msg, "port")
	}

	if source.Username == "" {
		msg = append(msg, "username")
	}

	if source.Password == "" {
		msg = append(msg, "password")
	}

	for _, db := range source.Databases {
		if db == "" {
			msg = append(msg, "database name")
		}
	}

	if len(msg) < 1 {
		return nil
	}

	return errors.New("config postgres: some field empty: " + strings.Join(msg, ","))
}

func checkSourceMongodbDriver(source ConfigSource) error {
	msg := []string{}

	if source.MongoDBURI == "" {
		msg = append(msg, "mongodb_uri")
	}

	for _, db := range source.Databases {
		if db == "" {
			msg = append(msg, "database name")
		}
	}

	if len(msg) < 1 {
		return nil
	}

	return errors.New("config mongodb: some field empty: " + strings.Join(msg, ","))
}

func checkSourceRedisDriver(source ConfigSource) error {
	msg := []string{}

	if source.Host == "" {
		msg = append(msg, "host")
	}

	if source.Port == "" {
		msg = append(msg, "port")
	}

	if len(msg) < 1 {
		return nil
	}

	return errors.New("config redis: some field empty: " + strings.Join(msg, ","))
}
