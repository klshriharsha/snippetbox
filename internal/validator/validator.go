package validator

import (
	"strings"
	"unicode/utf8"
)

// Validator struct handles form validation
type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if there are no field errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds a new error to FieldErrors map
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds a new error to FieldErrors map if `ok` is true
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank checks if `value` is empty
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars checks if the length of `value` is under `n`
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// ValidInt checks if `value` is one of `values`
func ValidInt(value int, values ...int) bool {
	for i := range values {
		if value == values[i] {
			return true
		}
	}

	return false
}