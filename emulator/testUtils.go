package emulator

import (
	"bytes"
	"io"
	"os"
)

func FileCompare(file1, file2 string) (bool, error) {
	const chunkSize = 64000

	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else if err1 != nil {
				return false, err1
			} else {
				return false, err2
			}
		}

		if !bytes.Equal(b1, b2) {
			return false, nil
		}
	}
}
