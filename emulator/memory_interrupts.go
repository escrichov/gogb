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

func (e *Memory) SetIF(value uint8) {
	e.io[271] = value
}

func (e *Memory) GetIF() uint8 {
	return e.io[271]
}

func (e *Memory) requestInterruptVBlank() {
	e.SetIF(SetBit8(e.GetIF(), InterruptVBlankBit, true))
}

func (e *Memory) requestInterruptLCDStat() {
	e.SetIF(SetBit8(e.GetIF(), InterruptLCDStatBit, true))
}

func (e *Memory) requestInterruptTimer() {
	e.SetIF(SetBit8(e.GetIF(), InterruptTimerBit, true))
}

func (e *Memory) requestInterruptSerial() {
	e.SetIF(SetBit8(e.GetIF(), InterruptSerialBit, true))
}

func (e *Memory) requestInterruptJoypad() {
	e.SetIF(SetBit8(e.GetIF(), InterruptJoypadBit, true))
}

func (e *Memory) SetIE(value uint8) {
	e.io[511] = value
}

func (e *Memory) GetIE() uint8 {
	return e.io[511]
}

func (e *Memory) hasPendingInterrupts() bool {
	interruptFlag := e.GetIF()
	interruptEnable := e.GetIE()
	pendingInterrupts := interruptFlag&interruptEnable != 0
	return pendingInterrupts
}

func (e *Memory) getInterruptType() InterruptType {
	interruptFlag := e.GetIF() & e.GetIE()
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
