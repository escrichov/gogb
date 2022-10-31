package emulator

import (
	"errors"
	"log"
	"os"
	"syscall"
	"unsafe"
)

func (e *Emulator) initializeSaveFile(fileName string) error {

	t := int(unsafe.Sizeof(uint8(8))) * 32768
	var mapFile *os.File

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		mapFile, err = os.Create(fileName)
		if err != nil {
			log.Println("Error opening file: ", err)
			return err
		}
		_, err = mapFile.Seek(int64(t-1), 0)
		if err != nil {
			log.Println("Error opening file: ", err)
			return err
		}
		_, err = mapFile.Write([]byte(" "))
		if err != nil {
			log.Println("Error writing file: ", err)
			return err
		}
	} else {
		mapFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Error opening file: ", err)
			return err
		}
	}

	mmap, err := syscall.Mmap(int(mapFile.Fd()), 0, int(t), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		log.Println(err)
		return err
	}

	e.extrambank = (*[32768]uint8)(unsafe.Pointer(&mmap[0]))

	return nil
}
