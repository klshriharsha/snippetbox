package config

import "log"

// application is used for dependency injection throughout the `web` application
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
