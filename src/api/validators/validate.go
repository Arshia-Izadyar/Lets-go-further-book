package validators

import (
	"regexp"
)

var (
	EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	EmailRX    = regexp.MustCompile(`^[\w-\.]+@(\w+\.)\w{2,4}$`)
)

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) AddError(key, value string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = value
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Check(ok bool, key, value string) {
	if !ok {
		v.AddError(key, value)
	}
}

func (v *Validator) In(value string, list ...string) bool {
	for i := range list {
		if list[i] == value {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique(values []string) bool {
	set := make(map[string]struct{})
	for _, i := range values {
		set[i] = struct{}{}
	}
	return len(set) == len(values)
}
