package internal_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wonksing/si/v2/internal"
)

func TestGetHmacSha256Pool(t *testing.T) {
	secret := "1234"

	hm := internal.GetHmacSha256Hash(secret)
	hm.Write([]byte("asdf"))
	hashed := hex.EncodeToString(hm.Sum(nil))
	internal.PutHmacSha256Hash(secret, hm)
	assert.EqualValues(t, "e5e9f44b2dcbe23988aa01743748a5fe64f494d7c5eea29ea94ae4e34878868e", hashed)

	hm = internal.GetHmacSha256Hash(secret)
	hm.Write([]byte("qwer"))
	hashed = hex.EncodeToString(hm.Sum(nil))
	internal.PutHmacSha256Hash(secret, hm)
	assert.EqualValues(t, "685f4fdb529e85b9e8fab7f9daaf550b5534e956d5c5f0f7a33c1ade0d8d67ea", hashed)

	secret = "asdf"
	hm = internal.GetHmacSha256Hash(secret)
	hm.Write([]byte("asdf"))
	hashed = hex.EncodeToString(hm.Sum(nil))
	internal.PutHmacSha256Hash(secret, hm)
	assert.EqualValues(t, "8a8423ba78c8f3da60a602493663c1cdc248a89541b12980e292399c0f0cad21", hashed)

	secret = "1234"
	hm = internal.GetHmacSha256Hash(secret)
	hm.Write([]byte("qwer"))
	hashed = hex.EncodeToString(hm.Sum(nil))
	internal.PutHmacSha256Hash(secret, hm)
	assert.EqualValues(t, "685f4fdb529e85b9e8fab7f9daaf550b5534e956d5c5f0f7a33c1ade0d8d67ea", hashed)
}
