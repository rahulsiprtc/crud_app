package config

import (
	"log"

	"github.com/caarlos0/env/v7"
)

type MongoConfig struct {
	MONGO_URI string `env:"MONGO_URI,required" envDefault:"mongodb://localhost:27017"`
	MONGO_DB  string `env:"MONGO_DB,required" envDefault:"crud_app"`
}

type AppConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
}

type BridgeConfig struct {
	Mongo MongoConfig `env:"MONGO_CONFIG"`
	App   AppConfig   `env:"APP_CONFIG"`
}

const ()

var Config *BridgeConfig

func InitializeConfig() {
	cfg := &BridgeConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	Config = cfg
	log.Println("Config loaded successfully")
}
