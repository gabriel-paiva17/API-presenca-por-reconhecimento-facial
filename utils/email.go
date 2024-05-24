package utils

import "regexp"

func IsValidEmail(email string) bool {
    // Regex para validar um email
    const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(emailRegexPattern)
    return re.MatchString(email)
}