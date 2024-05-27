package utils

import (
    "net"
    "regexp"
    "strings"
)

func IsValidEmailFormat(email string) bool {
    const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(emailRegexPattern)
    return re.MatchString(email)
}

func CheckMXRecords(domain string) bool {
    mxRecords, err := net.LookupMX(domain)
    return err == nil && len(mxRecords) > 0
}

func IsValidEmail(email string) bool {
    if !IsValidEmailFormat(email) {
        return false
    }

    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }

    domain := parts[1]
    return CheckMXRecords(domain)
}

// considerar envio de email de verificacao
// ex inicial: 

/*

func SendConfirmationEmail(email string, confirmationLink string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", "no-reply@example.com")
    m.SetHeader("To", email)
    m.SetHeader("Subject", "Please confirm your email address")
    m.SetBody("text/html", fmt.Sprintf("Click <a href=\"%s\">here</a> to confirm your email address.", confirmationLink))

    d := gomail.NewDialer("smtp.example.com", 587, "user", "password")
    return d.DialAndSend(m)
}

*/
