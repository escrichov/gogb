package emulator

const (
	InterruptVBlankAddress    uint16 = 0x40
	InterruptLCDStatAddress          = 0x48
	InterruptTimerAddress            = 0x50
	InterruptSerialAddress           = 0x58
	InterruptJoypadAddress           = 0x60
	InterruptCancelledAddress        = 0x00
)

const (
	InterruptVBlankBit    uint8 = 0
	InterruptLCDStatBit         = 1
	InterruptTimerBit           = 2
	InterruptSerialBit          = 3
	InterruptJoypadBit          = 4
	InterruptCancelledBit       = 255
)

type InterruptType struct {
	address uint16
	bit     uint8
}

func (mem *Memory) SetIF(value uint8) {
	mem.io[271] = value
}

func (mem *Memory) GetIF() uint8 {
	return mem.io[271]
}

func (mem *Memory) requestInterruptVBlank() {
	mem.SetIF(SetBit8(mem.GetIF(), InterruptVBlankBit, true))
}

func (mem *Memory) requestInterruptLCDStat() {
	mem.SetIF(SetBit8(mem.GetIF(), InterruptLCDStatBit, true))
}

func (mem *Memory) requestInterruptTimer() {
	mem.SetIF(SetBit8(mem.GetIF(), InterruptTimerBit, true))
}

func (mem *Memory) requestInterruptSerial() {
	mem.SetIF(SetBit8(mem.GetIF(), InterruptSerialBit, true))
}

func (mem *Memory) requestInterruptJoypad() {
	mem.SetIF(SetBit8(mem.GetIF(), InterruptJoypadBit, true))
}

func (mem *Memory) SetIE(value uint8) {
	mem.io[511] = value
}

func (mem *Memory) GetIE() uint8 {
	return mem.io[511]
}

func (mem *Memory) hasPendingInterrupts() bool {
	interruptFlag := mem.GetIF()
	interruptEnable := mem.GetIE()
	pendingInterrupts := interruptFlag&interruptEnable != 0
	return pendingInterrupts
}

func (mem *Memory) getInterruptType() InterruptType {
	interruptFlag := mem.GetIF() & mem.GetIE()
	if GetBit(interruptFlag, 0) { // VBlank
		return InterruptType{address: InterruptVBlankAddress, bit: InterruptVBlankBit}
	} else if GetBit(interruptFlag, 1) { // LCD STAT
		return InterruptType{address: InterruptLCDStatAddress, bit: InterruptLCDStatBit}
	} else if GetBit(interruptFlag, 2) { // Timer
		return InterruptType{address: InterruptTimerAddress, bit: InterruptTimerBit}
	} else if GetBit(interruptFlag, 3) { // Serial
		return InterruptType{address: InterruptSerialAddress, bit: InterruptSerialBit}
	} else if GetBit(interruptFlag, 4) { // Joypad
		return InterruptType{address: InterruptJoypadAddress, bit: InterruptJoypadBit}
	} else {
		return InterruptType{address: InterruptCancelledAddress, bit: InterruptCancelledBit}
	}
}
