package gb

import (
	apu "emulator-go/emulator/gb/audio"
	"emulator-go/emulator/gb/cpu"
	ppu "emulator-go/emulator/gb/graphics"
	mmu "emulator-go/emulator/gb/memory"
	"math"
)

type GB struct {
	clockHz        int
	fps            float64
	cyclesPerFrame int

	CPU cpu.CPU
	MMU mmu.MMU
	APU apu.APU
	PPU ppu.PPU

	totalCycles int64
}

func (gb *GB) Initialize(bootRom []byte, gameRom []byte) error {
	if err := gb.MMU.Init(bootRom, gameRom); err != nil {
		return err
	}
	gb.CPU.Init(true, &gb.MMU)
	gb.APU.Init(true)
	gb.PPU.Init(&gb.MMU)
	gb.clockHz = 1024 * 1024 * 4
	gb.fps = 59.727500569606
	gb.cyclesPerFrame = int(math.Round(float64(gb.clockHz) / gb.fps))

	return nil
}

func (gb *GB) Run() {
	var cycles = 0

	for cycles < gb.cyclesPerFrame {
		cycles += cpu.ExecuteNextOpcode(&gb.CPU)
	}
	gb.totalCycles += int64(cycles)
	gb.PPU.Update()
}
