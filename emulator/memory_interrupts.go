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

func (m *Memory) SetIF(value uint8) {
	m.io[271] = value
}

func (m *Memory) GetIF() uint8 {
	return m.io[271]
}

func (m *Memory) requestInterruptVBlank() {
	m.SetIF(SetBit8(m.GetIF(), InterruptVBlankBit, true))
}

func (m *Memory) requestInterruptLCDStat() {
	m.SetIF(SetBit8(m.GetIF(), InterruptLCDStatBit, true))
}

func (m *Memory) requestInterruptTimer() {
	m.SetIF(SetBit8(m.GetIF(), InterruptTimerBit, true))
}

func (m *Memory) requestInterruptSerial() {
	m.SetIF(SetBit8(m.GetIF(), InterruptSerialBit, true))
}

func (m *Memory) requestInterruptJoypad() {
	m.SetIF(SetBit8(m.GetIF(), InterruptJoypadBit, true))
}

func (m *Memory) SetIE(value uint8) {
	m.io[511] = value
}

func (m *Memory) GetIE() uint8 {
	return m.io[511]
}

func (m *Memory) hasPendingInterrupts() bool {
	interruptFlag := m.GetIF()
	interruptEnable := m.GetIE()
	pendingInterrupts := interruptFlag&interruptEnable != 0
	return pendingInterrupts
}

func (m *Memory) getInterruptType() InterruptType {
	interruptFlag := m.GetIF() & m.GetIE()
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
