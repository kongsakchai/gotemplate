package config

type Config struct {
	App      App
	Database Database
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
