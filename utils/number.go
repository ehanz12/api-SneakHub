package utils

import "strings"

func IsValidPhone(phone string) bool {
	return strings.HasPrefix(phone, "08")
}
