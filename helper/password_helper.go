package helper

import (
    "golang.org/x/crypto/bcrypt"
    "regexp"
)

// HashPassword hashes the given password with bcrypt and returns the hashed password or an error.
func HashPassword(password string) (string, error) {
    // Generate the hashed password with a salt (default cost of 10).
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plaintext password.
func CheckPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

func ValidateEmail(email string) bool{
    regexPattern:=`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re:= regexp.MustCompile(regexPattern)
    return re.MatchString(email)
}
