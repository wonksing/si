package sitcp

import (
	"time"

	"github.com/wonksing/si/v2/sio"
)

type TcpOption interface {
	apply(c *Conn)
}

type TcpOptionFunc func(*Conn)

func (s TcpOptionFunc) apply(c *Conn) {
	s(c)
}

func WithReaderOpt(opt sio.ReaderOption) TcpOptionFunc {
	return TcpOptionFunc(func(c *Conn) {
		c.appendReaderOption(opt)
	})
}

func WithWriterOpt(opt sio.WriterOption) TcpOptionFunc {
	return TcpOptionFunc(func(c *Conn) {
		c.appendWriterOption(opt)
	})
}

func WithEofChecker(chk sio.EofChecker) TcpOptionFunc {
	return TcpOptionFunc(func(c *Conn) {
		c.appendReaderOption(sio.SetEofChecker(chk))
	})
}

func WithWriteTimeout(writeTimeout time.Duration) TcpOptionFunc {
	return TcpOptionFunc(func(c *Conn) {
		c.writeTimeout = writeTimeout
	})
}

func WithReadTimeout(readTimeout time.Duration) TcpOptionFunc {
	return TcpOptionFunc(func(c *Conn) {
		c.readTimeout = readTimeout
	})
}

func WithWriteBufferSize(writeBufferSize int) TcpOptionFunc {
	return TcpOptionFunc(func(c *Conn) {
		c.writeBufferSize = writeBufferSize
	})
}

func WithReadBufferSize(readBufferSize int) TcpOptionFunc {
	return TcpOptionFunc(func(c *Conn) {
		c.readBufferSize = readBufferSize
	})
}
