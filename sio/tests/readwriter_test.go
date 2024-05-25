package sio_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/sio"
)

func TestReader_Buffer_Read(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.Write([]byte(testDataFile))

	r := sio.GetReader(buf)
	defer sio.PutReader(r)

	expected := testDataFile[:10]
	byt := make([]byte, 10)
	n, err := r.Read(byt)
	require.Nil(t, err)
	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)
}

func TestReader_Buffer_ReadBufio(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.Write([]byte(testDataFile))

	br := bufio.NewReader(buf)

	r := sio.GetReader(br)
	defer sio.PutReader(r)

	expected := testDataFile[:10]
	byt := make([]byte, 10)
	n, err := r.Read(byt)
	require.Nil(t, err)
	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)
}

func TestReader_Buffer_ReadAll(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.Write([]byte(testDataFile))

	r := sio.GetReader(buf)
	defer sio.PutReader(r)

	expected := testDataFile

	byt, err := r.ReadAll()
	require.Nil(t, err)
	assert.Equal(t, expected, string(byt))
}
