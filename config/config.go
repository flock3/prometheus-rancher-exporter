package config

import (
	"os"
	"strconv"

	cfg "github.com/infinityworksltd/go-common/config"
)

type Config struct {
	*cfg.BaseConfig
	rancherURL string
	accessKey  string
	secretKey  string
	hideSys    string
}

func (c Config) RancherURL() string {
	return c.rancherURL
}

func (c Config) AccessKey() string {
	return c.accessKey
}

func (c Config) SecretKey() string {
	return c.secretKey
}

func (c Config) HideSys() bool {
	s, _ := strconv.ParseBool(c.hideSys)
	return s
}

func Init() Config {
	ac := cfg.Init()

	appConfig := Config{
		&ac,
		os.Getenv("CATTLE_URL"),        // URL of Rancher Server API e.g. http://192.168.0.1:8080/v2-beta
		os.Getenv("CATTLE_ACCESS_KEY"), // Optional - Access Key for Rancher API
		os.Getenv("CATTLE_SECRET_KEY"), // Optional - Secret Key for Rancher API
		cfg.GetEnv("HIDE_SYS", "true"), // hideSys - Optional - Flag that indicates if the environment variable `HIDE_SYS` is set to a boolean true value
	}

	return appConfig
}
