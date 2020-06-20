package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var e = Environment{
	"BAR":   "bar",
	"UNSET": "",
}

func TestRunCmd(t *testing.T) {
	os.Setenv("UNSET", "remove me")
	require.Equal(t, "remove me", os.Getenv("UNSET"))

	code := RunCmd([]string{"echo", "something"}, e)
	require.Equal(t, 0, code)
	require.Equal(t, "bar", os.Getenv("BAR"))
	require.Equal(t, "", os.Getenv("UNSET"))

	code = RunCmd([]string{"cat", "nofile"}, e)
	require.Equal(t, 1, code)

	code = RunCmd([]string{"nocommand"}, e)
	require.Equal(t, -1, code)
}
