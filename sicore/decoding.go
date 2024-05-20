package sicore

import (
	"bytes"
	"errors"
	"io"
)

// DecodeJson read src with json bytes then decode it into dst.
func DecodeJson(dst any, src io.Reader) error {
	sr := GetReader(src, SetJsonDecoder())
	defer PutReader(sr)
	return sr.Decode(dst)
}

// DecodeJsonCopied read src with json bytes then decode it into dst.
// It also write the data read from src into a bytes.Buffer then returns it.
func DecodeJsonCopied(dst any, src io.Reader) (*bytes.Buffer, error) {
	bb := GetBytesBuffer(nil)
	r := io.TeeReader(src, bb)
	sr := GetReader(r, SetJsonDecoder())
	defer PutReader(sr)
	return bb, sr.Decode(dst)
}

// Decoder is an interface that has Decode method.
type Decoder interface {
	Decode(v any) error
}

// DefaultDecoder just read underlying r.
type DefaultDecoder struct {
	r io.Reader
}

// NewDefaultDecoder returns a new DefaultDecoder
func NewDefaultDecoder(r io.Reader) *DefaultDecoder {
	return &DefaultDecoder{r}
}

// Reset resets underyling Reader
func (d *DefaultDecoder) Reset(r io.Reader) {
	d.r = r
}

// Decode read underlying r, and pass the data to v.
// v must be either type of *[]byte or *string.
func (d *DefaultDecoder) Decode(v any) error {
	switch t := v.(type) {
	case *[]byte:
		b, err := io.ReadAll(d.r)
		if err != nil {
			return err
		}
		*t = b
		return nil
	case *string:
		b, err := io.ReadAll(d.r)
		if err != nil {
			return err
		}
		*t = string(b)
		return nil
	}

	return errors.New("not supported")
}
