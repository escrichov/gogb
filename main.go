package main

import (
	"emulator-go/emulator"
	"log"
)

func main() {

	// Passed
	// 01-special.gb
	// 02-interrupts.gb
	// 03-op sp,hl.gb
	// 04-op r,imm.gb
	// 05-op rp.gb
	// 06-ld r,r.gb
	// 07-jr,jp,call,ret,rst.gb
	// 08-misc instrs.gb
	// 09-op r,r.gb
	// 10-bit ops.gb
	// 11-op a,(hl).gb
	emulator, err := emulator.NewEmulator(
		"./assets/roms/gb-test-roms/instr_timing/instr_timing.gb",
		//"./assets/roms/pokeblue.gb",
		"rom.sav",
		//"./assets/roms/bootroms/DMG_ROM.bin",
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
