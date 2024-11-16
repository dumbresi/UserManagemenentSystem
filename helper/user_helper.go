package helper

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func Encode(email, token string) string {
	data := email + ":" + token
	encoded := base64.URLEncoding.EncodeToString([]byte(data))
	return encoded
}

func Decode(encoded string) (string, string, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}

	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid encoded string")
	}

	return parts[0], parts[1], nil
}