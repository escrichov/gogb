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

func (e *Emulator) SetIF(value uint8) {
	e.io[271] = value
}

func (e *Emulator) GetIF() uint8 {
	return e.io[271]
}

func (e *Emulator) requestInterruptVBlank() {
	e.SetIF(SetBit8(e.GetIF(), InterruptVBlankBit, true))
}

func (e *Emulator) requestInterruptLCDStat() {
	e.SetIF(SetBit8(e.GetIF(), InterruptLCDStatBit, true))
}

func (e *Emulator) requestInterruptTimer() {
	e.SetIF(SetBit8(e.GetIF(), InterruptTimerBit, true))
}

func (e *Emulator) requestInterruptSerial() {
	e.SetIF(SetBit8(e.GetIF(), InterruptSerialBit, true))
}

func (e *Emulator) requestInterruptJoypad() {
	e.SetIF(SetBit8(e.GetIF(), InterruptJoypadBit, true))
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

func (e *Emulator) getInterruptType() InterruptType {
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

func (e *Emulator) manageInterrupts() {
	interruptType := e.getInterruptType()
	e.halt = 0
	e.IME = 0 // IME should be 0 after a cancellation

	// Two wait states are executed (2 M-cycles pass while nothing happens; presumably the CPU is executing nops during this time).
	e.tick()
	e.tick()

	// The current value of the PC register is pushed onto the stack, consuming 2 more M-cycles.
	// The PC register is set to the address of the handler (one of: $40, $48, $50, $58, $60). This consumes one last M-cycle.
	// This is a regular call, exactly like what would be performed by a call <address> instruction
	// (the current PC is pushed onto the stack and then set to the address of the interrupt handler).
	e.instCall(interruptType.address)

	// IE register can be the target for one of the PC pushes during interrupt dispatch.
	// Only during upper byte push, SP after push must be 0xFFFE
	if e.cpu.SP.Get() == 0xFFFE {
		// Cancel execution of interrupt or execute a different interrupt
		interruptType = e.getInterruptType()
		e.cpu.PC = interruptType.address
	}

	// Clear interrupt IF bit if interrupt not cancelled
	if interruptType.address != InterruptCancelledAddress {
		// The IF bit corresponding to this interrupt is reset by the CPU
		e.SetIF(SetBit8(e.GetIF(), interruptType.bit, false))
	}

	if interruptType.address == InterruptVBlankAddress {
		e.loadGamesharkCodes()
	}
}
