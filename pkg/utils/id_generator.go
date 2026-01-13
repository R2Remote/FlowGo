package utils

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// GenerateID 生成唯一ID（使用时间戳+随机数）
func GenerateID() uint64 {
	now := time.Now().UnixNano()
	b := make([]byte, 8)
	rand.Read(b)
	random := binary.BigEndian.Uint64(b)
	return uint64(now) ^ random
}
