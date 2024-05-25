package codec

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/wonksing/si/v2/internal"
)

// EncodeJson encode src into json bytes then write to dst.
func EncodeJson(dst io.Writer, src any) error {
	sw := internal.GetWriter(dst, internal.SetJsonEncoder())
	defer internal.PutWriter(sw)
	return sw.EncodeFlush(src)
}

// EncodeJsonCopied encode src into json bytes then write to dst.
// It also write encoded bytes of src to a bytes.Buffer then returns it.
func EncodeJsonCopied(dst io.Writer, src any) (*bytes.Buffer, error) {
	bb := internal.GetBytesBuffer(nil)
	mw := io.MultiWriter(dst, bb)
	sw := internal.GetWriter(mw, internal.SetJsonEncoder())
	defer internal.PutWriter(sw)
	return bb, sw.EncodeFlush(src)
}

// HmacSha256HexEncoded creates an hmac sha256 hash from secret and mesage.
func HmacSha256HexEncoded(secret string, message []byte) (string, error) {
	hm := internal.GetHmacSha256Hash(secret)
	defer internal.PutHmacSha256Hash(secret, hm)
	_, err := hm.Write(message)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hm.Sum(nil)), nil
}

// HmacSha256HexEncodedWithReader creates an hmac sha256 hash from secret and r.
func HmacSha256HexEncodedWithReader(secret string, r io.Reader) (string, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return HmacSha256HexEncoded(secret, body)
}

// DecodeJson read src with json bytes then decode it into dst.
func DecodeJson(dst any, src io.Reader) error {
	sr := internal.GetReader(src, internal.SetJsonDecoder())
	defer internal.PutReader(sr)
	return sr.Decode(dst)
}

// DecodeJsonCopied read src with json bytes then decode it into dst.
// It also write the data read from src into a bytes.Buffer then returns it.
func DecodeJsonCopied(dst any, src io.Reader) (*bytes.Buffer, error) {
	bb := internal.GetBytesBuffer(nil)
	r := io.TeeReader(src, bb)
	sr := internal.GetReader(r, internal.SetJsonDecoder())
	defer internal.PutReader(sr)
	return bb, sr.Decode(dst)
}
