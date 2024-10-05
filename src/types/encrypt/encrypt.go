package encrypt

import (
	"crypto/subtle"
	"encoding/base64"

	"github.com/delyr1c/dechoric/src/types/cerr"
	"golang.org/x/crypto/scrypt"
)

// 使用固定的盐值
var fixedSalt = []byte("fixed_salt") // 请根据需要设置固定盐值

func Encrypt(password string) (string, error) {
	hash, err := scrypt.Key([]byte(password), fixedSalt, 16384, 8, 1, 32)
	if err != nil {
		return "", cerr.LogError(err) // 记录并返回错误
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

func DeCode(password string, hash string) (bool, error) {
	retrievedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, cerr.LogError(err)
	}
	// 使用固定盐值进行哈希计算
	newHash, err := scrypt.Key([]byte(password), fixedSalt, 16384, 8, 1, 32)
	if err != nil {
		return false, cerr.LogError(err)
	}
	return subtle.ConstantTimeCompare(newHash, retrievedHash) == 1, nil
}
