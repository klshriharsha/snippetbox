package config

import (
	"text/template"
	"time"
)

func standardDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"standardDate": standardDate,
}
