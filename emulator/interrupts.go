package emulator

func (e *Emulator) hastToManageInterrupts() bool {
	pendingInterrupts := e.mem.hasPendingInterrupts()
	return e.IME == 1 && pendingInterrupts
}

func (e *Emulator) manageInterrupts() {
	interruptType := e.mem.getInterruptType()
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
		interruptType = e.mem.getInterruptType()
		e.cpu.PC = interruptType.address
	}

	// Clear interrupt IF bit if interrupt not cancelled
	if interruptType.address != InterruptCancelledAddress {
		// The IF bit corresponding to this interrupt is reset by the CPU
		e.mem.SetIF(SetBit8(e.mem.GetIF(), interruptType.bit, false))
	}

	if interruptType.address == InterruptVBlankAddress {
		e.loadGameSharkCodes()
	}
}
