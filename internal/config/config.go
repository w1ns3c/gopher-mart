package config

import (
	"flag"
	"github.com/joho/godotenv"
	"gopher-mart/internal/domain"
	"os"
	"strconv"
	"time"
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

func LoadEnvfileConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	secret, exists := os.LookupEnv("SECRET")
	if exists {
		domain.Secret = secret
	}

	// tables
	orders, exists := os.LookupEnv("TableOrders")
	if exists {
		domain.TableOrders = orders
	}
	users, exists := os.LookupEnv("TableUsers")
	if exists {
		domain.TableUsers = users
	}
	balance, exists := os.LookupEnv("TableBalance")
	if exists {
		domain.TableBalance = balance
	}

	CookieName, exists := os.LookupEnv("CookieName")
	if exists {
		domain.CookieName = CookieName
	}
	CookieHoursLifeTime, exists := os.LookupEnv("CookieHoursLifeTime")
	if exists {
		val, err := strconv.ParseUint(CookieHoursLifeTime, 10, 64)
		if err != nil {
			return err
		}
		domain.CookieHoursLifeTime = time.Duration(val) * time.Hour
	}
	return nil
}
