package config

import (
	"flag"
	"os"
)

type Config struct {
	DBurl             string
	SrvAddr           string
	RemoteServiceAddr string
}

func LoadConfig() *Config {
	var (
		envSrvAddr = os.Getenv("RUN_ADDRESS")
		envDBurl   = os.Getenv("DATABASE_URI")
		envRSAddr  = os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	)

	config := &Config{
		DBurl:             envDBurl,
		SrvAddr:           envSrvAddr,
		RemoteServiceAddr: envRSAddr,
	}

	var flagSrvAddr, flagDBurl, flagRSAddr string
	flag.StringVar(&flagSrvAddr, "-a", "127.0.0.1:8000", "address for http server")
	flag.StringVar(&flagDBurl, "-d", "127.0.0.1:5432", "address for DB postgres")
	flag.StringVar(&flagRSAddr, "-r", "127.0.0.1:8001", "address for remote accrual system")
	flag.Parse()

	if config.SrvAddr == "" {
		config.SrvAddr = flagSrvAddr
	}

	if config.DBurl == "" {
		config.DBurl = flagDBurl
	}

	if config.RemoteServiceAddr == "" {
		config.RemoteServiceAddr = flagRSAddr
	}

	return config
}
