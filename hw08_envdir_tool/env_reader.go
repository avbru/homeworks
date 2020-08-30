package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]string

var ErrUnsupportedFileName = errors.New("unsupported filename")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, f := range files {
		if strings.Contains(f.Name(), "=") {
			return nil, ErrUnsupportedFileName
		}

		file, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			return nil, fmt.Errorf("envdir cannot open file: %w", err)
		}

		reader := bufio.NewReader(file)
		str, _ := reader.ReadString('\n')
		str = strings.TrimRight(str, "\r\n")
		str = strings.Replace(str, "\x00", "\n", -1)
		env[f.Name()] = str
	}

	return env, nil
}
