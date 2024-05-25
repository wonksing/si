package sisql

import (
	"github.com/wonksing/si/v2/internal/sio"
)

// SqlOption is an interface with apply method.
type SqlOption interface {
	apply(db *SqlDB)
}

// SqlOptionFunc wraps a function to conforms to SqlOption interface.
type SqlOptionFunc func(db *SqlDB)

// apply implements SqlOption's apply method.
func (o SqlOptionFunc) apply(db *SqlDB) {
	o(db)
}

func WithRowScannerOpt(opt sio.RowScannerOption) SqlOptionFunc {
	return SqlOptionFunc(func(db *SqlDB) {
		db.appendRowScannerOpt(opt)
	})
}

func WithTagKey(key string) SqlOptionFunc {
	return SqlOptionFunc(func(db *SqlDB) {
		db.appendRowScannerOpt(sio.WithTagKey(key))
	})
}

func WithType(name string, typ sio.SqlColType) SqlOptionFunc {
	return SqlOptionFunc(func(db *SqlDB) {
		db.appendRowScannerOpt(sio.WithSqlColumnType(name, typ))
	})
}

// SqlTxOption is an interface with apply method.
type SqlTxOption interface {
	apply(db *SqlTx)
}

// SqlTxOptionFunc wraps a function to conforms to SqlTxOption interface.
type SqlTxOptionFunc func(db *SqlTx)

// apply implements SqlOption's apply method.
func (o SqlTxOptionFunc) apply(db *SqlTx) {
	o(db)
}

func WithTxRowScannerOpt(opt sio.RowScannerOption) SqlTxOptionFunc {
	return SqlTxOptionFunc(func(db *SqlTx) {
		db.appendRowScannerOpt(opt)
	})
}

func WithTxTagKey(key string) SqlTxOptionFunc {
	return SqlTxOptionFunc(func(db *SqlTx) {
		db.appendRowScannerOpt(sio.WithTagKey(key))
	})
}

func WithTxType(name string, typ sio.SqlColType) SqlTxOptionFunc {
	return SqlTxOptionFunc(func(db *SqlTx) {
		db.appendRowScannerOpt(sio.WithSqlColumnType(name, typ))
	})
}
