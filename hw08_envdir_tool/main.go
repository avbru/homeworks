package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Wrong arguments. Usage: go-envdir /path/to/env/dir command arg1 arg2")
	}

	dir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatalf("error reading dir: %s\n", err)
	}

	RunCmd(command, env)
	os.Exit(0)
}
