package main

import (
	"WBTech0/config"
	"WBTech0/internal/app"
	"log"
	"os"
)

func main() {
	configPath := findConfigPath()

	cfg, err := config.Parse(configPath)
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}

func findConfigPath() string {
	const (
		devConfig  = "config/dev.config.toml"
		prodConfig = "config/config.toml"
	)

	if os.Getenv("CFG") == "DEV" {
		return devConfig
	}

	return prodConfig
}
