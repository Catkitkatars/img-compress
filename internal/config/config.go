package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath string = "config/local.yaml"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage-path" env-required:"true"`
	LogPath     string `yaml:"log-path" env-required:"true"`
	HTTP        HTTP   `yaml:"http"`
	Image       Image  `yaml:"image"`
}

type HTTP struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type Image struct {
	SavePath        string `yaml:"save-path" env-required:"true"`
	WmPath          string `yaml:"wm-path" env-required:"true"`
	Expansion       string `yaml:"expansion" env-default:".jpg"`
	CompressQuality int    `yaml:"compress-quality" env-default:"100"`
}

var Cfg Config

func New() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file - %s does not exist", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &Cfg); err != nil {

		log.Fatalf("Failed to read config file: %v", err)
	}
}
