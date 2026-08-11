package utils

func IsValidSize(size string) bool {
	switch size {
	case "35", "36", "37", "38", "39",
		"40", "41", "42", "43", "44",
		"45", "46", "47", "48":
		return true
	default:
		return false
	}
}
