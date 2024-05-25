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

func WithSqlColumnType(name string, columnType SqlColType) RowScannerOption {
	return RowScannerOptionFunc(func(rs *RowScanner) {
		switch columnType {
		case SqlColTypeBool:
			rs.SetSqlColumn(name, sqlNullBoolTypeValue)
		case SqlColTypeByte:
			rs.SetSqlColumn(name, sqlNullByteTypeValue)
		case SqlColTypeBytes:
			rs.SetSqlColumn(name, sqlBytesTypeValue)
		case SqlColTypeString:
			rs.SetSqlColumn(name, sqlNullStringTypeValue)
		case SqlColTypeInt:
			rs.SetSqlColumn(name, sqlNullIntTypeValue)
		case SqlColTypeInt8:
			rs.SetSqlColumn(name, sqlNullInt8TypeValue)
		case SqlColTypeInt16:
			rs.SetSqlColumn(name, sqlNullInt16TypeValue)
		case SqlColTypeInt32:
			rs.SetSqlColumn(name, sqlNullInt32TypeValue)
		case SqlColTypeInt64:
			rs.SetSqlColumn(name, sqlNullInt64TypeValue)
		case SqlColTypeUint:
			rs.SetSqlColumn(name, sqlNullUintTypeValue)
		case SqlColTypeUint8:
			rs.SetSqlColumn(name, sqlNullUint8TypeValue)
		case SqlColTypeUint16:
			rs.SetSqlColumn(name, sqlNullUint16TypeValue)
		case SqlColTypeUint32:
			rs.SetSqlColumn(name, sqlNullUint32TypeValue)
		case SqlColTypeUint64:
			rs.SetSqlColumn(name, sqlNullUint64TypeValue)
		case SqlColTypeFloat32:
			rs.SetSqlColumn(name, sqlNullFloat32TypeValue)
		case SqlColTypeFloat64:
			rs.SetSqlColumn(name, sqlNullFloat64TypeValue)
		case SqlColTypeTime:
			rs.SetSqlColumn(name, sqlNullTimeTypeValue)
		case SqlColTypeints:
			rs.SetSqlColumn(name, intsTypeValue)
		case SqlColTypeints8:
			rs.SetSqlColumn(name, ints8TypeValue)
		case SqlColTypeints16:
			rs.SetSqlColumn(name, ints16TypeValue)
		case SqlColTypeints32:
			rs.SetSqlColumn(name, ints32TypeValue)
		case SqlColTypeints64:
			rs.SetSqlColumn(name, ints64TypeValue)
		case SqlColTypeUints:
			rs.SetSqlColumn(name, uintsTypeValue)
		case SqlColTypeUints8:
			rs.SetSqlColumn(name, uints8TypeValue)
		case SqlColTypeUints16:
			rs.SetSqlColumn(name, uints16TypeValue)
		case SqlColTypeUints32:
			rs.SetSqlColumn(name, uints32TypeValue)
		case SqlColTypeUints64:
			rs.SetSqlColumn(name, uints64TypeValue)
		}
	})
}
