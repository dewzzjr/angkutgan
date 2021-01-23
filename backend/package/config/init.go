package config

import (
	"log"
	"os"
	"strconv"
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
	config = &Config{}
	if err := envconfig.Process("APP", &config.Delivery); err != nil {
		log.Fatalln("config.Delivery", env, err)
	}
	if err := envconfig.Process("APP", &config.Users); err != nil {
		log.Fatalln("config.Users", env, err)
	}
	if err := envconfig.Process("APP", &config.View); err != nil {
		log.Fatalln("config.View", env, err)
	}
	if err := envconfig.Process("APP", &config.Repository); err != nil {
		log.Fatalln("config.Repository", env, err)
	}
	config.Delivery.Port, _ = strconv.Atoi(os.Getenv("PORT"))
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
