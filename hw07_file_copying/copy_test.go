package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

const in = "test.file"
const out = in + ".copy"

func TestCopy(t *testing.T) {
	defer deleteFiles(t)

	err := Copy("strange\nfilename", "", 0, 0)
	require.NotNil(t, err, "cannot find file")

	err = Copy("/dev/urandom", in+".copy", 0, 0)
	require.Equal(t, ErrUnsupportedFile, err, "unsupported file")
	createFile(t, "")
	err = Copy(in, out, 1, 0)
	require.Equal(t, ErrOffsetExceedsFileSize, err, "file size less then offset")

	err = Copy(in, out, 0, 0)
	require.Equal(t, "", readFile(t), "zero sized file")

	createFile(t, "0123456789")
	err = Copy(in, out, 0, 0)
	require.Equal(t, "0123456789", readFile(t), "limit 0, offset 0")

	err = Copy(in, out, 5, 0)
	require.Equal(t, "56789", readFile(t), "offset 5, offset 0")

	err = Copy(in, out, 0, 5)
	require.Equal(t, "01234", readFile(t), "offset 0, limit 5")

	err = Copy(in, out, 9, 1)
	require.Equal(t, "9", readFile(t), "offset 9, limit 1")

}

func createFile(t *testing.T, data string) {
	err := ioutil.WriteFile(in, []byte(data), 0664)
	if err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T) string {
	data, err := ioutil.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func deleteFiles(t *testing.T) {
	err := os.Remove(in)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(out)
	if err != nil {
		t.Fatal(err)
	}
}
