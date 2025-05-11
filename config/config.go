package config

import "time"

type Config struct {
	App      App
	Database Database
	Redis    Redis
}

type App struct {
	Service  string
	Port     string
	Env      string
	LogLevel string
}

type Database struct {
	MySQLURI string
	Env      string
}

type Redis struct {
	Host           string
	Port           string
	Password       string
	DB             int
	Timeout        time.Duration
	ConnectTimeout time.Duration
}
