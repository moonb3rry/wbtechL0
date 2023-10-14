package config

import (
	"github.com/BurntSushi/toml"
	"time"
)

type (
	Config struct {
		App        App        `toml:"Application"`
		Db         Db         `toml:"DB"`
		HttpServer HttpServer `toml:"HttpServer"`
		Nats       Nats       `toml:"Nats"`
	}

	App struct {
		Name    string `toml:"Name"`
		Version string `toml:"Version"`
	}

	Db struct {
		Name        string `toml:"Name"`
		Host        string `toml:"Host"`
		Port        string `toml:"Port"`
		Schema      string `toml:"Schema"`
		MaxPoolSize int    `toml:"MaxPoolSize"`
		User        string `env:"DBUSER"`
		Password    string `env:"DBPASSWORD"`
	}

	HttpServer struct {
		ReadTimeout     *time.Duration `toml:"ReadTimeout"`
		WriteTimeout    *time.Duration `toml:"WriteTimeout"`
		Addr            string         `toml:"Addr"`
		ShutdownTimeout *time.Duration `toml:"ShutdownTimeout"`
	}

	Nats struct {
		URL string `toml:"URL"`
	}
)

func Parse(path string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
