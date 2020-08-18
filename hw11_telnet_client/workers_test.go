package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"
)

func TestLostConnection(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	tc := NewTelnetClient(l.Addr().String(), time.Second, os.Stdin, os.Stdout)
	require.NoError(t, tc.Connect())
	ctx, cancel := context.WithCancel(context.Background())

	stderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	go listen(ctx, tc, cancel)

	require.NoError(t, l.Close())
	time.Sleep(time.Second)
	require.NoError(t, tc.Close())
	require.NoError(t, w.Close())

	out, _ := ioutil.ReadAll(r)
	os.Stderr = stderr
	require.Equal(t, "...Connection closed by peer\n", string(out))

	_, ok := <-ctx.Done()
	require.False(t, ok)
}

func TestEOF(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	tc := NewTelnetClient(l.Addr().String(), time.Second, os.Stdin, os.Stdout)
	require.NoError(t, tc.Connect())

	stderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	serve(context.Background(), tc)

	require.NoError(t, w.Close())

	out, _ := ioutil.ReadAll(r)
	os.Stderr = stderr
	require.Equal(t, "EOF\n", string(out))
}
