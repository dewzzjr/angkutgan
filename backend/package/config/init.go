package config

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var config *Config

func init() {
	if config == nil {
		load(
			"./config.yaml",
			"../config.yaml",
			"../../config.yaml",
			"./configs/angkutgan.yaml",
			"../configs/angkutgan.yaml",
			"/etc/conf/angkutgan.yaml",
		)
	}
}

func load(path ...string) {
	for _, p := range path {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		if err = yaml.NewDecoder(f).Decode(&config); err != nil {
			continue
		}
		validate()
	}
	if config == nil {
		log.Println("no config loaded")
	}
}

func validate() {
	if config.View.Port == 0 {
		config.View.Port = config.Delivery.Port
	}
	config.View.Path = check(config.View.Path, 0)
}

func check(path string, tries int) string {
	if _, err := os.Stat(path); os.IsNotExist(err) && tries < 3 {
		path = "../" + strings.TrimPrefix(path, "/")
		return check(path, tries+1)
	}
	return path
}

// Load config from specific
func Load(path string) (err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	if err = yaml.NewDecoder(f).Decode(&config); err != nil {
		return
	}
	return
}

// Get configuration
func Get() *Config {
	if config == nil {
		log.Fatal("no config found")
	}
	return config
}
