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

	flag.Parse()

	if *configPath == "" {
		fmt.Println("Usage : bekup --config /path/to/file.json")
		os.Exit(1)
	}

	conf, err := config.LoadConfigFromPath(*configPath)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	ctx := context.Background()

	worker := conf.Worker

	err = bekup.Run(ctx, conf, worker)
	if err != nil {
		fmt.Println("Error : ", err)
	}
}
