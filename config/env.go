package config

import "os"

var Env string

func init() {
	Env = os.Getenv("ENV")
}

const (
	Local string = "LOCAL"
	Dev   string = "DEV"
	Prod  string = "PROD"
)

func IsLocal() bool {
	return Env == Local
}

func IsDev() bool {
	return Env == Dev
}

func IsProd() bool {
	return Env == Prod
}
