package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var env = Environment{
	"BAR":   "bar",
	"FOO":   "   foo\nwith new line",
	"HELLO": "\"hello\"",
	"UNSET": "",
}

func TestReadDir(t *testing.T) {
	res, err := ReadDir("")
	require.NotNil(t, err)
	require.Nil(t, res)

	res, err = ReadDir("testdata/env")
	require.Nil(t, err)
	require.Equal(t, env, res)
}
