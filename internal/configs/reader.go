package configs

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func ReadConfig() Config {
	log.Println("Loading corvette.toml configuration file.")
	configContent := loadConfigFile()

	var conf Config
	_, err := toml.Decode(configContent, &conf)
	if err != nil {
		log.Panic(err)
	}

	return conf
}

func loadConfigFile() string {
	filePath := "corvette.toml"
	dat, err := os.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	return string(dat)
}
