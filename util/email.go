package util

import (
	"errors"
	"regexp"
	"strings"
)

// Validate if email address is formally correct
func Validate(addr string) bool {
	emailValidator := regexp.MustCompile(`[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]{2,}`)
	valid := emailValidator.MatchString(addr)
	return valid
}

// GetEmailDomain extracts domain host from a given email address
func GetEmailDomain(addr string) (string, error) {
	_, domain, err := SplitEmail(addr)
	return domain, err
}

// SplitEmail extracts name and domain host from email
func SplitEmail(addr string) (local, domain string, err error) {
	if !Validate(addr) {
		return "", "", errors.New("mta: invalid mail address")
	}
	parts := strings.SplitN(addr, "@", 2)

	if len(parts) != 2 {
		return "", "", errors.New("mta: invalid mail address")
	}
	return parts[0], parts[1], nil
}
