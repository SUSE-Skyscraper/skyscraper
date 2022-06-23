package application

import (
	"github.com/spf13/viper"
)

type OktaConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Issuer   string `mapstructure:"issuer"`
	ClientID string `mapstructure:"client_id"`
}

type FrontendConfig struct {
	URL string `mapstructure:"url"`
}

type DBConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Port     int64  `mapstructure:"port"`
	Host     string `mapstructure:"host"`
}

type Config struct {
	Okta     OktaConfig     `mapstructure:"okta"`
	DB       DBConfig       `mapstructure:"db"`
	Frontend FrontendConfig `mapstructure:"frontend"`
}

func Configuration() (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	// set defaults
	conf := Config{
		Okta: OktaConfig{
			Enabled: false,
		},
	}

	err = v.Unmarshal(&conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
