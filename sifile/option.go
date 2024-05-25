package sifile

import (
	"github.com/wonksing/si/v2/internal/sio"
)

// FileOption is an interface with apply method.
type FileOption interface {
	apply(f *File)
}

// FileOptionFunc wraps a function to conforms to FileOption interface
type FileOptionFunc func(f *File)

// apply implements FileOption's apply method.
func (o FileOptionFunc) apply(f *File) {
	o(f)
}

func WithReaderOpt(opt sio.ReaderOption) FileOptionFunc {
	return FileOptionFunc(func(f *File) {
		f.appendReaderOpt(opt)
	})
}

func WithWriterOpt(opt sio.WriterOption) FileOptionFunc {
	return FileOptionFunc(func(f *File) {
		f.appendWriterOpt(opt)
	})
}
