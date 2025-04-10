package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncryptPassword(password string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}
