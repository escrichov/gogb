package main

import (
	"emulator-go/emulator"
	"log"
)

func main() {

	// Failed
	// 01-special.gb
	// 02-interrupts.gb
	// 03-op sp,hl.gb
	// 07-jr,jp,call,ret,rst.gb
	// 08-misc instrs.gb
	// 11-op a,(hl).gb

	// Passed
	// 04-op r,imm.gb
	// 05-op rp.gb
	// 06-ld r,r.gb
	// 09-op r,r.gb
	// 10-bit ops.gb
	emulator, err := emulator.NewEmulator("./assets/roms/pokeblue.bin", "save.bin", "")
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
