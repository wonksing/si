package sicore

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
)

// Encoder encode input parameter and write to a writer.
type Encoder interface {
	Encode(v any) error
}

// SetJsonEncoder is a WriterOption to encode w's data in json format
func SetJsonEncoder() WriterOption {
	return WriterOptionFunc(func(w *Writer) {
		w.enc = json.NewEncoder(w)
	})
}

// DefaultEncoder is to write string or []byte type to the underlying Writer
type DefaultEncoder struct {
	w io.Writer
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

// SetDefaultEncoder sets DefaultEncoder to w
func SetDefaultEncoder() WriterOption {
	return WriterOptionFunc(func(w *Writer) {
		w.enc = &DefaultEncoder{w}
	})
}

// Decoder is an interface that has Decode method.
type Decoder interface {
	Decode(v any) error
}

// SetJsonDecoder sets json.Decoder to r.
func SetJsonDecoder() ReaderOption {
	return ReaderOptionFunc(func(r *Reader) {
		r.dec = json.NewDecoder(r)
	})
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
	body, err := ReadAll(r)
	if err != nil {
		return "", err
	}
	return HmacSha256HexEncoded(secret, body)
}
