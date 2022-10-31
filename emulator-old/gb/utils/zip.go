package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
)

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func UnzipBytes(bs []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(bs), int64(len(bs)))
	if err != nil {
		return nil, err
	}

	// Read first file from zip archive
	for _, zipFile := range zipReader.File {
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			return nil, err
		}

		return unzippedFileBytes, nil
	}

	return nil, fmt.Errorf("no files in zip archive")
}
