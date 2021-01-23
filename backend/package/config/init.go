package config

import (
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var config *Config

func init() {
	if config == nil {
		env := os.Getenv("ENVIRONMENT")
		if env == "" {
			loadFile()
		} else {
			loadEnv(env)
		}
		if config == nil {
			log.Fatalln("no config loaded")
		}
		validate()
	}
}

func loadFile() {
	path := []string{
		"./config.yaml",
		"../config.yaml",
		"../../config.yaml",
		"./configs/angkutgan.yaml",
		"../configs/angkutgan.yaml",
		"../../configs/angkutgan.yaml",
		"/etc/conf/angkutgan.yaml",
	}
	for _, p := range path {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		if err = yaml.NewDecoder(f).Decode(&config); err != nil {
			continue
		}
	}
}

func loadEnv(env string) {
	if err := envconfig.Process("APP", &config); err != nil {
		log.Fatalln("envconfig.Process", env, err)
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
