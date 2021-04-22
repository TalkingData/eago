package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

// GenSha256HashCode 生成SHA256哈希值
func GenSha256HashCode(message string) string {
	hash := sha256.New()
	hash.Write([]byte(message))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}
