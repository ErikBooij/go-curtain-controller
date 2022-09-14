package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	API     APIConfig  `yaml:"api"`
	MQTT    MQTTConfig `yaml:"mqtt"`
	Devices struct {
		AqaraShutters map[string]AqaraShutterConfig `yaml:"aqara-shutters"`
		SlideCurtains map[string]SlideCurtainConfig `yaml:"slide-curtains"`
	} `yaml:"devices"`
}

type APIConfig struct {
	Bind string `yaml:"bind"`
	Port uint16 `yaml:"port"`
}

type MQTTConfig struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ClientId string `yaml:"client-id"`
}

type AqaraShutterConfig struct {
	Topic string `yaml:"topic"`
}

type SlideCurtainConfig struct {
	DeviceID string `yaml:"device-id"`
	IP       string `yaml:"ip"`
}

func LoadConfig(path string) AppConfig {
	content, err := os.ReadFile(path)

	if err != nil {
		panic(fmt.Errorf("unable to read config file: %s", err))
	}

	appConfig := AppConfig{}

	err = yaml.Unmarshal(content, &appConfig)

	if err != nil {
		panic(fmt.Errorf("unable to parse config file: %s", err))
	}

	return appConfig
}
