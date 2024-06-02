package siutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

func HmacSha256(message []byte, hmacKey []byte) hash.Hash {
	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(message)
	return mac
}

func HmacSha256HexStr(message []byte, hmacKey []byte) string {
	mac := HmacSha256(message, hmacKey)
	return hex.EncodeToString(mac.Sum(nil))
}
