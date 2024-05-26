package sio_test

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/sio"
)

func testCreateFileToRead(fileName, data string) error {
	fr, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer fr.Close()

	_, err = fr.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

func TestReader_File_Read(t *testing.T) {
	fileName := "./data/TestReader_File_Read.txt"
	require.Nil(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	require.Nil(t, err)
	defer f.Close()

	s := sio.GetReader(f)
	defer sio.PutReader(s)

	expected := testDataFile[:10]

	byt := make([]byte, 10)
	n, err := s.Read(byt)
	require.Nil(t, err)

	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)

	fileName2 := "./data/TestReader_Read_2.txt"
	require.Nil(t, testCreateFileToRead(fileName2, testDataFile2))

	f2, err := os.OpenFile(fileName2, os.O_RDONLY, 0644)
	require.Nil(t, err)
	defer f2.Close()

	s.Reset(f2, sio.SetDefaultEOFChecker())

	expected = testDataFile2[:10]

	// byt := make([]byte, 10)
	n, err = s.Read(byt)
	require.Nil(t, err)

	assert.Equal(t, expected, string(byt))
	assert.Equal(t, 10, n)
}

func TestReader_File_ReadAll(t *testing.T) {
	fileName := "./data/TestReader_File_ReadAll.txt"
	require.Nil(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	require.Nil(t, err)
	defer f.Close()

	s := sio.GetReader(f)
	defer sio.PutReader(s)

	b, err := s.ReadAll()
	require.Nil(t, err)

	str := strings.ReplaceAll(string(b), "\r\n", "\n")
	assert.Equal(t, testDataFile, str)
}

func TestReader_File_ReadSmall(t *testing.T) {
	fileName := "./data/TestReader_File_ReadSmall.txt"
	require.Nil(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	require.Nil(t, err)
	defer f.Close()

	s := sio.GetReader(f)
	defer sio.PutReader(s)
	b := make([]byte, 1)
	n, err := s.Read(b)
	require.Nil(t, err)

	assert.EqualValues(t, "{", string(b))
	assert.Equal(t, 1, n)
}

func TestReader_File_ReadZeroCase1(t *testing.T) {
	fileName := "./data/TestReader_File_ReadZeroCase1.txt"
	require.Nil(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	require.Nil(t, err)
	defer f.Close()

	s := sio.GetReader(f)
	defer sio.PutReader(s)

	var b []byte
	n, err := s.Read(b)
	require.Nil(t, err)

	assert.Equal(t, 0, n)
}

func TestReader_File_ReadZeroCase2(t *testing.T) {
	fileName := "./data/TestReader_File_ReadZeroCase2.txt"
	require.Nil(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	require.Nil(t, err)
	defer f.Close()

	s := sio.GetReader(f)
	defer sio.PutReader(s)

	b := make([]byte, 0, len(testDataFile))
	n, err := s.Read(b)
	require.Nil(t, err)
	assert.Equal(t, 0, n)

}

func TestReader_File_Decode(t *testing.T) {
	fileName := "./data/TestReader_File_Decode.txt"
	require.Nil(t, testCreateFileToRead(fileName, testDataFile))

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	require.Nil(t, err)
	defer f.Close()

	r := sio.GetReader(f, sio.SetJsonDecoder())
	defer sio.PutReader(r)

	var p Person
	require.Nil(t, r.Decode(&p))
	assert.EqualValues(t, Person{Name: "wonk", Age: 20, Email: "wonk@wonk.org"}, p)
}

func TestWriter_File_Write(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_Write.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	require.Nil(t, err)

	s := sio.GetWriter(f)
	defer sio.PutWriter(s)

	expected := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	expected += "\n"
	n, err := s.Write([]byte(expected))
	require.Nil(t, err)
	require.Nil(t, s.Flush())

	assert.EqualValues(t, len(expected), n)
}

func TestWriter_File_WriteMany(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_WriteMany.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	require.Nil(t, err)
	defer f.Close()

	s := sio.GetWriter(f)
	defer sio.PutWriter(s)
	line := `{"name":"wonk","age":20,"email":"wonk@wonk.org"}`
	line += "\n"
	expected := bytes.Repeat([]byte(line), 1000)
	n, err := s.Write(expected)
	require.Nil(t, err)
	require.Nil(t, s.Flush())
	assert.Equal(t, len(line)*1000, n)

}

type Person struct {
	Name           string `json:"name"`
	Age            uint8  `json:"age"`
	Email          string `json:"email"`
	Gender         string `json:"gender"`
	MarriageStatus string `json:"marriage_status"`
	NumChildren    uint8  `json:"num_children"`
}

func TestWriter_File_EncodeDefaultEncoderByte(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_EncodeDefaultEncoderByte.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	require.Nil(t, err)
	defer f.Close()

	s := sio.GetWriter(f, sio.SetDefaultEncoder())
	defer sio.PutWriter(s)
	byt := []byte(`{"name":"wonk","age":20,"email":"wonk@wonk.wonk","gender":"M","marriage_status":"Yes","num_children":10}`)

	err = s.Encode(byt)
	require.Nil(t, err)
	require.Nil(t, s.Flush())

	err = s.Encode(&byt)
	require.Nil(t, err)
	require.Nil(t, s.Flush())

}
func TestWriter_File_EncodeDefaultEncoderString(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_EncodeDefaultEncoderString.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	require.Nil(t, err)
	defer f.Close()

	s := sio.GetWriter(f, sio.SetDefaultEncoder())
	defer sio.PutWriter(s)
	str := `{"name":"wonk","age":20,"email":"wonk@wonk.wonk","gender":"M","marriage_status":"Yes","num_children":10}`

	err = s.Encode(str)
	require.Nil(t, err)
	require.Nil(t, s.Flush())

	err = s.Encode(&str)
	require.Nil(t, err)
	require.Nil(t, s.Flush())
}
func TestWriter_File_WriteAnyStruct(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_WriteAnyStruct.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	require.Nil(t, err)
	defer f.Close()

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sio.GetWriter(f, sio.SetJsonEncoder())
	defer sio.PutWriter(s)

	err = s.Encode(p)
	require.Nil(t, err)
	require.Nil(t, s.Flush())
}

func TestWriter_File_EncodeJsonEncodeStruct(t *testing.T) {

	t.Run("succeed-1", func(t *testing.T) {
		f, err := os.OpenFile("./data/TestWriter_File_EncodeJsonEncodeStruct.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		require.Nil(t, err)
		defer f.Close()

		// bufio readwriter wrap around f
		bw := bufio.NewWriter(f)
		s := sio.GetWriter(bw, sio.SetJsonEncoder())
		defer sio.PutWriter(s)

		p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}
		require.Nil(t, s.Encode(p))
		require.Nil(t, s.Flush())
	})

	t.Run("succeed-2", func(t *testing.T) {
		f, err := os.OpenFile("./data/TestWriter_File_EncodeJsonEncodeStruct.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		require.Nil(t, err)
		defer f.Close()

		// bufio readwriter wrap around f
		bw := bufio.NewWriter(f)
		s := sio.GetWriter(bw, sio.SetJsonEncoder())
		defer sio.PutWriter(s)

		p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}
		require.Nil(t, s.Encode(p))
		require.Nil(t, s.Flush())
	})
}

func TestWriter_File_EncodeNoEncoderFail(t *testing.T) {
	f, err := os.OpenFile("./data/TestWriter_File_EncodeNoEncoderFail.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	require.Nil(t, err)
	defer f.Close()

	p := &Person{"wonk", 20, "wonk@wonk.wonk", "M", "Yes", 10}

	s := sio.GetWriter(f)
	err = s.Encode(p)
	require.NotNil(t, err)
	require.Nil(t, s.Flush())

}
