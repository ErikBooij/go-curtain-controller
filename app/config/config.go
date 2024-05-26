package config

import (
	"fmt"
	"os"
	"regexp"

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
	Auth     bool   `yaml:"auth"`
}

func LoadConfig(path string) AppConfig {
	content, err := os.ReadFile(path)

	if err != nil {
		panic(fmt.Errorf("unable to read config file: %s", err))
	}

	placeholderRegex := regexp.MustCompile(`env\((?:(.+?)(?::(.+))?)\)`)
	// specRegex := regexp.MustCompile()

	content = placeholderRegex.ReplaceAllFunc(content, func(bytes []byte) []byte {
		subMatches := placeholderRegex.FindStringSubmatch(string(bytes))

		envValue, ok := os.LookupEnv(subMatches[1])

		if ok {
			return []byte(envValue)
		}

		return []byte(subMatches[2])
	})

	appConfig := AppConfig{}

	err = yaml.Unmarshal(content, &appConfig)

	if err != nil {
		panic(fmt.Errorf("unable to parse config file: %s", err))
	}

	return appConfig
}
