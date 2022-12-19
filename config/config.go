package config

import (
<<<<<<< HEAD
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

// Config for this project
type Config struct {
	DB                     string
	ProductTypes           []string
	TelegramChatID         string
	TelegramToken          string
	PostmarkApiToken       string // token để gởi email của postmark.com
	ESMS_APIKEY            string // apikey cua esms.vn
	ESMS_SECRET            string // secret cua esms.vn
	B2Prefix               string
	UploadingRoot          string
	MappingUploadLocalLink string
	B2BucketName           string
	B2AccountId            string
	B2AccountKey           string
	ServerIP               string
	APTPaymentNotifyUrl    string
	APTPaymentRedirectUrl  string
	APTPartnerCode         string
	APTApiKey              string
	APTSecretKey           string
	APTPaymentHost         string
	APTEbillHost           string
	APTEbillNotifyUrl      string
	CDNUrl                 string
	MinCostPerSlot         int
}

// Default config
func (c *Config) Default() error {
	if c.DB == "" {
		return errors.New("please config postgres db (config.json)")
	}

	if c.B2Prefix == "" {
		c.B2Prefix = ""
	}

	if c.MappingUploadLocalLink == "" {
		c.MappingUploadLocalLink = "/"
	}

	if c.CDNUrl == "" {
		c.CDNUrl = "https://cdn.deein.com"
	}

	if c.MinCostPerSlot == 0 {
		c.MinCostPerSlot = 10000000
	}

	if len(c.ProductTypes) == 0 {
		c.ProductTypes = []string{
			"Căn hộ chung cư",
			"Nhà riêng",
			"Nhà biệt thự, liền kề",
			"Nhà mặt phố",
			"Shop house, nhà phố thương mại",
			"Đất nền dự án",
			"Đất",
			"Trang trại, khu nghỉ dưỡng",
			"Condotel",
			"Kho, nhà xưởng",
			"Các loại hình bất động sản khác",
		}
	}

	return nil
}

func (c *Config) Save() error {
	log.Println("CONFIG: save config.json")
	b, _ := json.MarshalIndent(c, "", "  ")
	return ioutil.WriteFile("config.json", b, 0644)
}

// Parse config from json byte
func Parse(b []byte, path string) (*Config, error) {
	c := &Config{}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	if err := c.Default(); err != nil {
		return nil, err
	}

	return c, nil
}

// Read the config from `path`
func Read(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b, path)
}

// Create the sample config file
func Create() error {
	log.Println("CONFIG: create config.json")
	b, _ := json.MarshalIndent(Config{}, "", "  ")
	return ioutil.WriteFile("config.json", b, 0644)
=======
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
>>>>>>> 5ed448758ec77912d8f15b1cd516eb78d5d5ec71
}
