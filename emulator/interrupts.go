package emulator

const (
	InterruptVBlank  uint16 = 0x40
	InterruptLCDStat        = 0x48
	InterruptTimer          = 0x50
	InterruptSerial         = 0x58
	InterruptJoypad         = 0x60
)

func (e *Emulator) SetIF(value uint8) {
	e.io[271] = value
}

func (e *Emulator) GetIF() uint8 {
	return e.io[271]
}

func (e *Emulator) SetInterruptVBlank() {
	e.SetIF(e.GetIF() | 1)
}

func (e *Emulator) SetInterruptLCDStat() {
	e.SetIF(e.GetIF() | 2)
}

func (e *Emulator) SetInterruptTimer() {
	e.SetIF(e.GetIF() | 4)
}

func (e *Emulator) SetInterruptSerial() {
	e.SetIF(e.GetIF() | 8)
}

func (e *Emulator) SetInterruptJoypad() {
	e.SetIF(e.GetIF() | 16)
}

func (e *Emulator) SetIE(value uint8) {
	e.io[511] = value
}

func (e *Emulator) GetIE() uint8 {
	return e.io[511]
}

func (e *Emulator) hasPendingInterrupts() bool {
	interruptFlag := e.GetIF()
	interruptEnable := e.GetIE()
	pendingInterrupts := interruptFlag&interruptEnable != 0
	return pendingInterrupts
}

func (e *Emulator) hastToManageInterrupts() bool {
	pendingInterrupts := e.hasPendingInterrupts()
	return e.IME == 1 && pendingInterrupts
}

func (e *Emulator) getInterruptType() uint16 {
	interruptFlag := e.GetIF()
	if GetBit(interruptFlag, 0) { // VBlank
		return InterruptVBlank
	} else if GetBit(interruptFlag, 1) { // LCD STAT
		return InterruptLCDStat
	} else if GetBit(interruptFlag, 2) { // Timer
		return InterruptTimer
	} else if GetBit(interruptFlag, 3) { // Serial
		return InterruptSerial
	} else if GetBit(interruptFlag, 4) { // Joypad
		return InterruptJoypad
	} else {
		return InterruptVBlank
	}
}

func (e *Emulator) manageInterrupts() {
	interruptType := e.getInterruptType()
	e.halt = 0

	// The IF bit corresponding to this interrupt and the IME flag are reset by the CPU
	e.SetIF(0)
	e.IME = 0

	// Two wait states are executed (2 M-cycles pass while nothing happens; presumably the CPU is executing nops during this time).
	e.tick()
	e.tick()

	// The current value of the PC register is pushed onto the stack, consuming 2 more M-cycles.
	e.push(e.cpu.PC)

	// The PC register is set to the address of the handler (one of: $40, $48, $50, $58, $60). This consumes one last M-cycle.
	e.cpu.PC = interruptType

	if interruptType == InterruptVBlank {
		e.loadGamesharkCodes()
	}
}
