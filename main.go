package main

import (
	"emulator-go/emulator"
	"log"
)

func main() {

	emulator, err := emulator.NewEmulator(
		"./assets/roms/roms/supermarioland.gb",
		//"./assets/roms/bootroms/dmg_boot.bin",
		"",
		"./assets/fonts/arial.ttf",
		true)
	if err != nil {
		log.Fatal(err)
	}
	defer emulator.Destroy()

	//err = emulator.BessLoad("rom.s3")
	//if err != nil {
	//	log.Fatal(err)
	//}

	emulator.Run()
}
