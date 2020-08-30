package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const chunks = 20

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	stat, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("copy cannot get source file stats: %w", err)
	}

	if !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}

	file, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("copy cannot open source file: %w", err)
	}

	defer file.Close()
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("copy cannot rewind file to offset: %w", err)
	}

	size := stat.Size()

	total := size - offset
	if total >= limit && limit != 0 {
		total = limit
	}
	chunk := maxInt64(total/chunks, 1)

	outfile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("copy cannot open destination file: %w", err)
	}
	defer outfile.Close()

	for written := int64(0); written < total; {
		wb, err := io.CopyN(outfile, file, chunk)

		written += wb
		progressBar(float32(written) / float32(total))

		switch err {
		case io.EOF:
			return nil
		case nil:
		default:
			return fmt.Errorf("copy cannot copy bytes: %w", err)
		}
	}
	return nil
}

//progressBar prints progressbar for 0 < p < 1. Step=10%.
func progressBar(p float32) {
	n := int(p * 10)                                                                                                    //nolint:gomnd
	s := "\u001b[1000D copying [" + strings.Repeat("=", n) + strings.Repeat(" ", 10-n) + "]" + strconv.Itoa(n*10) + "%" //nolint:gomnd
	_, _ = os.Stdout.WriteString(s)
}

func maxInt64(a, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}
