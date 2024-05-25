package internal

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
)

// EncodeJson encode src into json bytes then write to dst.
func EncodeJson(dst io.Writer, src any) error {
	sw := GetWriter(dst, SetJsonEncoder())
	defer PutWriter(sw)
	return sw.EncodeFlush(src)
}

// EncodeJsonCopied encode src into json bytes then write to dst.
// It also write encoded bytes of src to a bytes.Buffer then returns it.
func EncodeJsonCopied(dst io.Writer, src any) (*bytes.Buffer, error) {
	bb := GetBytesBuffer(nil)
	mw := io.MultiWriter(dst, bb)
	sw := GetWriter(mw, SetJsonEncoder())
	defer PutWriter(sw)
	return bb, sw.EncodeFlush(src)
}

// HmacSha256HexEncoded creates an hmac sha256 hash from secret and mesage.
func HmacSha256HexEncoded(secret string, message []byte) (string, error) {
	hm := GetHmacSha256Hash(secret)
	defer PutHmacSha256Hash(secret, hm)
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

// Encoder encode input parameter and write to a writer.
type Encoder interface {
	Encode(v any) error
}

// DefaultEncoder is to write string or []byte type to the underlying Writer
type DefaultEncoder struct {
	w io.Writer
}

// Reset resets underyling Writer
func (de *DefaultEncoder) Reset(w io.Writer) {
	de.w = w
}

// Encode writes v to underyling Writer only when its type is []byte, string or pointer to these two.
func (de *DefaultEncoder) Encode(v any) error {
	if v == nil {
		return nil
	}

	switch c := v.(type) {
	case []byte:
		_, err := de.w.Write(c)
		return err
	case *[]byte:
		_, err := de.w.Write(*c)
		return err
	case string:
		_, err := de.w.Write([]byte(c))
		return err
	case *string:
		_, err := de.w.Write([]byte(*c))
		return err
	default:
		return errors.New("unable to encode v")
	}

}
