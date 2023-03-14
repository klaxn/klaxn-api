package config

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App            CoreConfig        `yaml:"app"`
	OutboundConfig []*OutboundConfig `yaml:"outbounds"`
	Aws            map[string]string `yaml:"aws"`
	DatabaseConfig *DatabaseConfig   `yaml:"database-config"`
}

type CoreConfig struct {
	Name string `yaml:"name"`
}
type OutboundConfig struct {
	Name    string                 `yaml:"name"`
	Enabled bool                   `yaml:"enabled"`
	Config  map[string]interface{} `yaml:"config"`
}

type DatabaseConfig struct {
	Hostname     string `yaml:"hostname"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database-name"`
}

func (dc *DatabaseConfig) SetFromEnvOrDefault() {
	if dc.Hostname == "" {
		dc.Hostname = getEnvOrDefault("DB_HOSTNAME", "localhost")
	}

	if dc.Port == "" {
		dc.Port = getEnvOrDefault("DB_PORT", "5432")
	}

	// TODO: set better database defaults
	if dc.User == "" {
		dc.User = getEnvOrDefault("DB_USER", "gorm")
	}

	if dc.Password == "" {
		dc.Password = getEnvOrDefault("DB_PASSWORD", "gorm")
	}

	if dc.DatabaseName == "" {
		dc.DatabaseName = getEnvOrDefault("DB_NAME", "postgres")
	}
}

func getEnvOrDefault(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}

func New() (*Config, error) {
	c := &Config{}

	yamlFile, err := ioutil.ReadFile("klaxn.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}

	if c.DatabaseConfig == nil {
		c.DatabaseConfig = &DatabaseConfig{}
	}

	c.DatabaseConfig.SetFromEnvOrDefault()

	return c, nil
}

func (c *Config) GetAWSConfig() (aws.Config, error) {
	if c.Aws["profile"] != "" {
		return config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
			o.SharedConfigProfile = c.Aws["profile"]
			o.Region = "eu-west-1"
			return nil
		})
	}

	return config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "eu-west-1"
		return nil
	})
}
