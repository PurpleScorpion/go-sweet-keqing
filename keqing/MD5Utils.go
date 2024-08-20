package keqing

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	encryptedPassword := md5.Sum([]byte(str))
	hashedPassword := hex.EncodeToString(encryptedPassword[:])
	return hashedPassword
}

func MD5Salt(str string, salt string) string {
	return MD5(str + salt)
}
