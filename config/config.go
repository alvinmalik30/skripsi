package config

import (
	"fmt"
	"os"
	"polen/utils/common"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// db config
type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Driver   string
}

// api config
type ApiConfig struct {
	ApiHost string
	ApiPort string
}

// file config (logger)
type FileConfig struct {
	FilePath string
}

type TokenConfig struct {
	ApplicationName  string
	JwtSignatureKey  []byte
	JwtSigningMethod *jwt.SigningMethodHMAC
	ExpirationToken  int
}

type Config struct {
	DbConfig
	ApiConfig
	FileConfig
	TokenConfig
}

func (c *Config) ReadConfig() error {
	// env
	err := common.LoadEnv()
	if err != nil {
		return err
	}

	// environtment Variable
	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	// api config env
	c.ApiConfig = ApiConfig{
		ApiHost: os.Getenv("API_HOST"),
		ApiPort: os.Getenv("API_PORT"),
	}

	// logged file path
	c.FileConfig = FileConfig{
		FilePath: os.Getenv("FILE_PATH"),
	}

	expiration, err := strconv.Atoi(os.Getenv("APP_EXPIRATION_TOKEN"))
	if err != nil {
		return err
	}

	c.TokenConfig = TokenConfig{
		ApplicationName:  os.Getenv("APP_TOKEN_NAME"),
		JwtSignatureKey:  []byte(os.Getenv("APP_TOKEN_KEY")),
		JwtSigningMethod: jwt.SigningMethodHS256,
		ExpirationToken:  expiration,
	}

	if c.DbConfig.Host == "" ||
		c.DbConfig.Port == "" ||
		c.DbConfig.Name == "" ||
		c.DbConfig.User == "" ||
		c.DbConfig.Password == "" ||
		c.DbConfig.Driver == "" ||
		c.ApiConfig.ApiHost == "" ||
		c.ApiConfig.ApiPort == "" ||
		c.FileConfig.FilePath == "" {
		return fmt.Errorf("missing requirenment variable")
	}

	return nil
}

// constructor
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.ReadConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
