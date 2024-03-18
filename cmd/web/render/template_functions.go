package render

import (
	"text/template"
	"time"
)

// standardDate formats date into a human-readable format
func standardDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// map of custom template functions to register before parsing any template
var functions = template.FuncMap{
	"standardDate": standardDate,
}
