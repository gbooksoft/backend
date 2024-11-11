package password

import (
	"unicode"
	"unicode/utf8"
)

type PasswordValidationResult struct {
	LengthError      bool
	DigitError       bool
	UppercaseError   bool
	LowercaseError   bool
	SpecialCharError bool
	Errors           []string
}

func (p *PasswordValidationResult) Error() string {
	if len(p.Errors) > 0 {
		return p.Errors[0]
	}

	return ""
}

func ValidatePassword(password string) *PasswordValidationResult {
	result := &PasswordValidationResult{}

	if utf8.RuneCountInString(password) < 8 {
		result.LengthError = true
		result.Errors = append(result.Errors, "Пароль должен содержать не менее 8 символов")
	}

	hasDigit := false
	hasUpper := false
	hasLower := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// Проверка на наличие цифры
	if !hasDigit {
		result.DigitError = true
		result.Errors = append(result.Errors, "Пароль должен содержать хотя бы одну цифру")
	}

	// Проверка на наличие заглавной буквы
	if !hasUpper {
		result.UppercaseError = true
		result.Errors = append(result.Errors, "Пароль должен содержать хотя бы одну заглавную букву")
	}

	// Проверка на наличие строчной буквы
	if !hasLower {
		result.LowercaseError = true
		result.Errors = append(result.Errors, "Пароль должен содержать хотя бы одну строчную букву")
	}

	// Проверка на наличие специального символа
	if !hasSpecial {
		result.SpecialCharError = true
		result.Errors = append(result.Errors, "Пароль должен содержать хотя бы один специальный символ")
	}

	return result
}
