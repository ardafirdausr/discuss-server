package usecase

import (
	"crypto/sha1"
	"fmt"
)

func hashString(param string) string {
	hash := sha1.New()
	hash.Write([]byte(param))
	hashedPass := hash.Sum(nil)
	return fmt.Sprintf("%x", hashedPass)
}
