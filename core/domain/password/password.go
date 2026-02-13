package password

import (
	"fmt"
	_errors "password-validator/core/errors"
	constants "password-validator/core/utils"
	"unicode"
)

type (
	Password struct {
		password string
		isValid  bool
	}

	PasswordParams func(p *Password)
)

func New(params ...PasswordParams) (*Password, error) {
	p := &Password{}
	for _, param := range params {
		param(p)
	}

	err := p.validate()
	if err != nil {
		p.isValid = false
		return p, err
	}
	p.isValid = true
	return p, nil
}

func (p *Password) validate() error {
	trimmed := ""
	for _, c := range p.password {
		if !unicode.IsSpace(c) {
			trimmed += string(c)
		}
	}

	if len(trimmed) < 9 {
		return _errors.InvalidField{
			Field: "password",
			AsIs:  "Must have at least 9 characters (excluding spaces)",
		}
	}

	digit := false
	lower := false
	upper := false
	special := false
	seen := make(map[rune]bool)

	for _, c := range trimmed {
		if seen[c] {
			return _errors.InvalidField{
				Field: "password",
				AsIs:  "Must not contain repeated characters (excluding spaces)",
			}
		}
		seen[c] = true

		switch {
		case unicode.IsDigit(c):
			digit = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case containsRune(constants.SPECIAL_CHARS, c):
			special = true
		}
	}

	if !digit {
		return _errors.InvalidField{
			Field: "password",
			AsIs:  "Must contain at least one digit (excluding spaces)",
		}
	}
	if !lower {
		return _errors.InvalidField{
			Field: "password",
			AsIs:  "Must contain at least one lowercase letter (excluding spaces)",
		}
	}
	if !upper {
		return _errors.InvalidField{
			Field: "password",
			AsIs:  "Must contain at least one uppercase letter (excluding spaces)",
		}
	}
	if !special {
		return _errors.InvalidField{
			Field: "password",
			AsIs:  fmt.Sprintf("Must contain at least one special character (%s, excluding spaces)", constants.SPECIAL_CHARS),
		}
	}

	return nil
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

func WithPassword(password string) PasswordParams {
	return func(p *Password) {
		p.password = password
	}
}

func (p *Password) Password() string {
	return p.password
}

func (p *Password) IsValid() bool {
	return p.isValid
}
