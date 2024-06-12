package util

import "errors"

func IsSlug(s string) bool {
	for _, c := range s {
		if isUppercase(c) {
			return false
		}
		if !isLetter(c) && !isNumber(c) && c != '-' {
			return false
		}
	}
	return true
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isNumber(c rune) bool {
	return (c >= '0' && c <= '9')
}

func isUppercase(c rune) bool {
	return (c >= 'A' && c <= 'Z')
}

var SlugValidateFunc = func(s string) error {
	if len(s) < 3 {
		return errors.New("application name must be at least 3 characters")
	}
	if len(s) > 50 {
		return errors.New("application name must be at most 50 characters")
	}
	if !IsSlug(s) {
		return errors.New("application name must be a valid slug")
	}
	return nil
}
