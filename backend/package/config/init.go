package config

import (
	"log"
	"os"
	"strings"

	"github.com/vrischmann/envconfig"
	"gopkg.in/yaml.v2"
)

var config *Config

func init() {
	if config == nil {
		env := os.Getenv("ENVIRONMENT")
		switch env {
		case "PRODUCTION", "STAGING":
			loadEnv(env)
		default:
			loadFile(
				"./configs/angkutgan.yaml",
				"../configs/angkutgan.yaml",
			)
		}
		if config == nil {
			log.Fatalln("no config loaded")
		}
	}
}

func loadFile(path ...string) {
	for _, p := range path {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		if err = yaml.NewDecoder(f).Decode(&config); err != nil {
			continue
		}
		config.ViewPort = config.Port
		config.Path = pathChecker(config.Path)
	}
}

func loadEnv(env string) {
	if err := envconfig.InitWithOptions(&config, envconfig.Options{AllOptional: true}); err != nil {
		log.Fatalln(env, err)
	}
	const (
		UseBuildIn   = "USE_BUILD_IN"
		BuildInDBKey = "JAWSDB_MARIA_URL"
	)
	if _, ok := os.LookupEnv(UseBuildIn); ok {
		db := os.Getenv(BuildInDBKey)
		db = strings.Trim(db, "mysql://")
		db = strings.Replace(db, "@", "@tcp(", 1)
		db = strings.Replace(db, "/", ")/", 1)
		config.DatabaseURL = db
	}
}

func pathChecker(path string) string {
	return check(path, 3)
}

func check(path string, tries int) string {
	if _, err := os.Stat(path); os.IsNotExist(err) && tries > 0 {
		path = "../" + strings.TrimPrefix(path, "./")
		return check(path, tries-1)
	}
	return path
}

// Get configuration
func Get() *Config {
	if config == nil {
		log.Fatal("no config found")
	}
	return config
}
