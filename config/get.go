package config

import "os"

func getAppConfig() App {
	return App{
		Name:    os.Getenv("APP_NAME"),
		Port:    os.Getenv("APP_PORT"),
		Version: os.Getenv("APP_VERSION"),
	}
}

func getHeaderConfig() Header {
	return Header{
		RefIDKey: os.Getenv("HEADER_REF_ID_KEY"),
	}
}
