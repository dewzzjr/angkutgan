package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var config *Config

func init() {
	if config == nil {
		load(
			"./config.yaml",
			"../config.yaml",
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
	}
	if config == nil {
		log.Println("no config loaded")
	}
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
