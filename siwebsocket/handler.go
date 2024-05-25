package siwebsocket

import (
	"io"
	"log"

	"github.com/wonksing/si/v2/internal/sio"
)

// MessageHandler handles data read from r with ReaderOption opts.
// Returning error from Handle method stops client's reading and writing. For siwebsocket.Client,
// Handle method should return error if reading returns an error other than EOF.
type MessageHandler interface {
	Handle(r io.Reader, opts ...sio.ReaderOption) error
}

type NopMessageHandler struct{}

func (o *NopMessageHandler) Handle(r io.Reader, opts ...sio.ReaderOption) error {
	// discard reader data
	_, err := io.Copy(io.Discard, r)
	return err
}

type DefaultMessageHandler struct{}

func (o *DefaultMessageHandler) Handle(r io.Reader, opts ...sio.ReaderOption) error {
	sr := sio.GetReader(r, opts...)
	defer sio.PutReader(sr)

	_, err := sr.ReadAll()
	if err != nil {
		return err
	}

	return nil
}

type DefaultMessageLogHandler struct{}

func (o *DefaultMessageLogHandler) Handle(r io.Reader, opts ...sio.ReaderOption) error {
	sr := sio.GetReader(r, opts...)
	defer sio.PutReader(sr)

	b, err := sr.ReadAll()
	if err != nil {
		return err
	}
	log.Println(string(b))
	return nil
}
