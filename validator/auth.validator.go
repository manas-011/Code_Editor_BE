package validator

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsStrongPassword(p string) bool {
	if len(p) < 8 {
		return false
	}

	var u, l, n, s bool
	for _, c := range p {
		switch {
		case unicode.IsUpper(c):
			u = true
		case unicode.IsLower(c):
			l = true
		case unicode.IsDigit(c):
			n = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			s = true
		}
	}
	return u && l && n && s
}
