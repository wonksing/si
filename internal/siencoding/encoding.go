package siencoding

import (
	"errors"
	"io"
)

// Encoder encode input parameter and write to a writer.
type Encoder interface {
	Encode(v any) error
}

// NewDefaultEncoder returns an instance of DefaultEncoder
func NewDefaultEncoder(w io.Writer) *DefaultEncoder {
	return &DefaultEncoder{w}
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
