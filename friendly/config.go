package friendly

import (
	"github.com/spf13/viper"
)

type Config struct {
	configPath string
	configType string
}

func NewConfig(configPath string, configType string) *Config {
	return &Config{
		configPath: configPath,
		configType: configType,
	}
}

func (c *Config) Init() error {
	viper.SetConfigFile(c.configPath)
	viper.SetConfigType(c.configType)
	return viper.ReadInConfig()
}
