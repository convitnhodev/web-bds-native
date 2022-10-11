package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	*viper.Viper
}

func New() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// Select the .env file
	config.SetConfigName(".env")
	config.SetConfigType("dotenv")
	config.AddConfigPath(".")

	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}

	// TODO: Logger (Maybe a different zap object)

	// TODO: Add APP_KEY generation

	// TODO: Write changes to configuration file
	return config
}

func (c *Config) setDefaults() {
	// Set default App configuration
	c.SetDefault("APP_ADDR", ":3000")
	c.SetDefault("APP_ENV", "local")
	// Set default database configuration
	c.SetDefault("DB_URI", "postgresql://postgres:postgres@localhost/postgres")
}
