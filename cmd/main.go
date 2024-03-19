package main

import (
	"flag"
	"log"

	"vk_quests/config"
	"vk_quests/internal/app"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
