package config

import (
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
	UploadToB2At           string
	ServerIP               string
	ATPNotifyUrl           string
	ATPRedirectUrl         string
	CDNUrl                 string
}

// Default config
func (c *Config) Default() error {
	if c.DB == "" {
		return errors.New("please config postgres db (config.json)")
	}

	if c.B2BucketName == "" {
		return errors.New("please config B2BucketName backblaze(config.json)")
	}

	if c.B2AccountId == "" {
		return errors.New("please config B2AccountId backblaze(config.json)")
	}

	if c.B2AccountKey == "" {
		return errors.New("please config B2AccountKey backblaze(config.json)")
	}

	/*	TODO: Enable when ready
		if c.ServerIP == "" {
			return errors.New("please config ServerIP (config.json)")
		}

		if c.ATPNotifyUrl == "" {
			return errors.New("please config ATPNotifyUrl (config.json)")
		}

		if c.ATPRedirectUrl == "" {
			return errors.New("please config ATPRedirectUrl (config.json)")
		}
	*/

	if c.UploadToB2At == "" {
		c.UploadToB2At = "19:00"
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
}
