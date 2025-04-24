package configs

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Storage struct {
	Type   string `yaml:"Type"`
	Source string `yaml:"Source"`
	Host   string `yaml:"Host"`
	Port   string `yaml:"Port"`
	User   string `yaml:"User"`
	Pass   string `yaml:"Pass"`
}

type Config struct {
	Storage *Storage `yaml:"Storage"`
}

var conf *Config

func Load(filename string) (*Config, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("config file read error: %v", err)
	}
	if err := yaml.Unmarshal(yamlFile, &conf); err != nil {
		return nil, fmt.Errorf("config file unmarshal error: %v", err)
	}
	return conf, nil
}
