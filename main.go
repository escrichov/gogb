package main

import (
	"emulator-go/emulator"
	"log"
)

func main() {

	emu, err := emulator.NewEmulator(
		"./assets/roms/games/supermarioland.gb",
		//"./assets/roms/bootroms/dmg_boot.bin",
		"",
		"./assets/fonts/arial.ttf",
		true,
		true)
	if err != nil {
		log.Fatal(err)
	}
	defer emu.Destroy()

	//err = emulator.BessLoad("rom.s3")
	//if err != nil {
	//	log.Fatal(err)
	//}

	emu.Run(0)
}
