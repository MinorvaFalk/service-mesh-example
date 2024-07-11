package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type config struct {
	Env  string `mapstructure:"ENV"`
	Port string `mapstructure:"port"`
	Nsq  nsq    `mapstructure:"nsq"`
}

type nsq struct {
	Host    string  `mapstructure:"host"`
	Port    string  `mapstructure:"port"`
	Lookupd lookupd `mapstructure:"lookupd"`
	Topic   string  `mapstructure:"topic"`
	Channel string  `mapstructure:"channel"`
}

type lookupd struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func (l lookupd) Address() string {
	return l.Host + ":" + l.Port
}

func (n nsq) Address() string {
	return n.Host + ":" + n.Port
}

var c *config

func InitConfig(configPath ...string) {
	vi := viper.New()

	if len(configPath) != 0 && configPath[0] != "" {
		vi.AddConfigPath(configPath[0])
	} else {
		if _, err := os.Stat("./config.yaml"); err == nil {
			vi.SetConfigFile("./config.yaml")
		}
	}

	setDefault(vi)
	vi.AutomaticEnv()

	if vi.ConfigFileUsed() != "" {
		if err := vi.ReadInConfig(); err != nil {
			log.Fatalf("failed to read config file: %v", err)
		}
	}

	if err := vi.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}

func ReadConfig(configPath ...string) *config {
	if c == nil {
		InitConfig(configPath...)
	}

	return c
}

func setDefault(vi *viper.Viper) {
	vi.SetDefault("ENV", "development")

	vi.SetDefault("port", "8080")
	vi.SetDefault("nsq.host", "127.0.0.1")
	vi.SetDefault("nsq.port", "4150")
	vi.SetDefault("nsq.topic", "notification")
	vi.SetDefault("nsq.channel", "urgent")
	vi.SetDefault("nsq.lookupd.host", "127.0.0.1")
	vi.SetDefault("nsq.lookupd.port", "4161")
}
