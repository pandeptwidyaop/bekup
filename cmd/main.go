package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pandeptwidyaop/bekup/internal/bekup"
	"github.com/pandeptwidyaop/bekup/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage : bekup </path/to/file/config.json>")
		os.Exit(1)
	}

	configPath := os.Args[1]

	conf, err := config.LoadConfigFromPath(configPath)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	ctx := context.Background()

	worker := 2

	err = bekup.Run(ctx, conf, worker)
	if err != nil {
		fmt.Println("Error : ", err)
	}
}
