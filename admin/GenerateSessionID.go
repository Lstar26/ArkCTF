package admin

import (
	"crypto/rand"
	"encoding/hex"
)

// 生成加密的session_id
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16) // 生成16字节的随机数
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil // 转换为十六进制字符串
}
