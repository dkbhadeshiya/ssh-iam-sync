package main

import (
	"fmt"

	"github.com/kkyr/fig"
)

// Config is the configuration for the application
type Config struct {
	// Aws configuration
	Aws struct {
		Method    string   `fig:"method" default:"profile"`
		Region    string   `fig:"region" default:"eu-west-1"`
		AccessKey string   `fig:"accessKey"`
		SecretKey string   `fig:"secretKey"`
		Profile   string   `fig:"profile" default:"default"`
		Groups    []string `fig:"groups"`
	} `fig:"aws"`
	// AuthorizedKeys file path to write
	AuthorizedKeys string `fig:"authorizedKeys" validate:"required"`
	// Over write existing keys
	Overwrite bool `fig:"overwrite"`
}

// Parse the configurations of the application
func GetAppConfig() Config {
	var cfg Config
	if err := fig.Load(&cfg, fig.File("config.yaml"), fig.Dirs("./", "/etc/ssh-iam-sync")); err != nil {
		panic(fmt.Errorf("error loading config: %v", err))
	}
	return cfg
}
