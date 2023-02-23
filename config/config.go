package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Db struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbName"`
		Port     string `yaml:"port"`
	} `yaml:"db"`
	Aws struct {
		InstanceId string
	}
	Sns struct {
		ServerEndpoint      string `yaml:"serverEndpoint"`
		Region             string `yaml:"region"`
		AlbumUpdateTopicArn string `yaml:"albumUpdateTopicArn"`
		Credentials        struct {
			Id     string `yaml:"id"`
			Secret string `yaml:"secret"`
			Token  string `yaml:"token"`
		} `yaml:"credentials"`
	}
	Sqs struct {
		ServerEndpoint     string `yaml:"serverEndpoint"`
		Region             string `yaml:"region"`
		Credentials        struct {
			Id     string `yaml:"id"`
			Secret string `yaml:"secret"`
			Token  string `yaml:"token"`
		} `yaml:"credentials"`
	} `yaml:"sqs"`
}

func NewConfig(configEnvironment string) *Config {
	// Create config structure
	config := &Config{}

	configPath := fmt.Sprintf("./config/%v/config.yaml", configEnvironment)
	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		panic(err)
	}

	// env var overrides. use viper someday.
	config.Aws.InstanceId = "i-12345ABCDEF12345"

	return config
}
