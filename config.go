package main

import (
	"log"

	"github.com/jinzhu/configor"
)

// Config contains global app's configuration
var Config AppConfig

// AppConfig contains app's configuration
type AppConfig struct {
	Server struct {
		Port string `default:":80"`
	}
	DB struct {
		Path string `default:"./data.db"`
	}
}

// LoadFromFile method loads and parses config file
func (c *AppConfig) LoadFromFile(url string) {
	err := configor.Load(&Config, url)
	if err != nil {
		log.Fatalf("Can't load the config file: %s", err)
	}
}

func (c *AppConfig) DataSourceName() string {
	return c.DB.Path
}
