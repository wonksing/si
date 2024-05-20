package sicore

import "encoding/json"

// WriterOption is an interface that has apply method.
type WriterOption interface {
	apply(w *Writer)
}

// WriterOptionFunc wraps a function to conforms to WtierOption interface
type WriterOptionFunc func(*Writer)

// apply implements WriterOption's apply method.
func (s WriterOptionFunc) apply(w *Writer) {
	s(w)
}

// SetJsonEncoder is a WriterOption to encode w's data in json format
func SetJsonEncoder() WriterOption {
	return WriterOptionFunc(func(w *Writer) {
		w.SetEncoder(json.NewEncoder(w))
	})
}

// SetDefaultEncoder sets DefaultEncoder to w
func SetDefaultEncoder() WriterOption {
	return WriterOptionFunc(func(w *Writer) {
		w.SetEncoder(&DefaultEncoder{w})
	})
}

// ReaderOption is an interface that wraps an apply method.
type ReaderOption interface {
	apply(r *Reader)
}

// ReaderOptionFunc wraps a function to conforms to ReaderOption's apply method.
type ReaderOptionFunc func(*Reader)

// apply implements ReaderOption's apply method.
func (o ReaderOptionFunc) apply(r *Reader) {
	o(r)
}

// SetEofChecker sets chk to r.
func SetEofChecker(chk EofChecker) ReaderOption {
	return ReaderOptionFunc(func(r *Reader) {
		r.SetEofChecker(chk)
	})
}

func SetDefaultEOFChecker() ReaderOption {
	return ReaderOptionFunc(func(r *Reader) {
		r.SetEofChecker(&defaultEofChecker{})
	})
}

// SetJsonDecoder sets json.Decoder to r.
func SetJsonDecoder() ReaderOption {
	return ReaderOptionFunc(func(r *Reader) {
		r.SetDecoder(json.NewDecoder(r))
	})
}

// RowScannerOption is an interface that wraps an apply method.
type RowScannerOption interface {
	apply(rs *RowScanner)
}

// RowScannerOptionFunc wraps a function to conforms to RowScannerOption's apply method.
type RowScannerOptionFunc func(rs *RowScanner)

// apply implements RowScannerOption's apply method.
func (o RowScannerOptionFunc) apply(rs *RowScanner) {
	o(rs)
}

// WithTagKey sets key to rs's tagKey.
func WithTagKey(key string) RowScannerOption {
	return RowScannerOptionFunc(func(rs *RowScanner) {
		rs.SetTagKey(key)
	})
}
