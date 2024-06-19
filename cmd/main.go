package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pandeptwidyaop/bekup/internal/bekup"
	"github.com/pandeptwidyaop/bekup/internal/config"
)

func main() {

	configPath := flag.String("config", "", "Path to configuration file")
	tempPath := flag.String("temp", "", "Path to temporary path")

	flag.Parse()

	if *configPath == "" {
		fmt.Println("Usage : bekup --config /path/to/file.json")
		fmt.Printf("Optional:\n		--temp <path>\n			Path to temporary path")
		os.Exit(1)
	}

	conf, err := config.LoadConfigFromPath(*configPath)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	if tempPath != nil && *tempPath != "" {
		fmt.Println("using argument temporary path")
		conf.TempPath = *tempPath
	}

	ctx := context.Background()

	worker := conf.Worker

	err = bekup.Run(ctx, conf, worker)
	if err != nil {
		fmt.Println("Error : ", err)
	}
}
