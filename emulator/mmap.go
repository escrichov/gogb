package emulator

import (
	"log"
	"os"
	"syscall"
)

func createFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
}

func resizeFile(file *os.File, sizeBytes int) error {
	_, err := file.Seek(int64(sizeBytes-1), 0)
	if err != nil {
		log.Println("Error opening file: ", err)
		return err
	}
	_, err = file.Write([]byte(" "))
	if err != nil {
		log.Println("Error writing file: ", err)
		return err
	}

	return nil
}

func createMMAP(fileName string, sizeBytes int) ([]byte, error) {
	file, err := createFile(fileName)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() < int64(sizeBytes) {
		err := resizeFile(file, sizeBytes)
		if err != nil {
			return nil, err
		}
	}

	mmapData, err := syscall.Mmap(
		int(file.Fd()),
		0,
		sizeBytes,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		return nil, err
	}

	return mmapData, nil
}
