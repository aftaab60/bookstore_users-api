package crypto_utils

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

func GetMd5Hash(input string) string {
	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetBcryptHash(input string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareBcryptHashWithPassword(hash string, input string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input)); err != nil {
		return err
	}
	return nil
}