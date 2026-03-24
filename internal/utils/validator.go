package utils

import "regexp"

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ValidateUsername 验证用户名格式
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}
