package emulator

func (e *Emulator) getSerialTransferData() uint8 { // 0xFF01
	return e.mem.io[257]
}

func (e *Emulator) setSerialTransferData(value uint8) { // 0xFF01
	e.mem.io[257] = value
}

func (e *Emulator) getSerialTransferControl() uint8 { // 0xFF02
	return e.mem.io[258]
}

func (e *Emulator) setSerialTransferControl(value uint8) { // 0xFF02
	e.mem.io[258] = value
}

func (e *Emulator) serialTransfer() {
	cyclesElapsed := uint16(e.cycles - e.prevCycles)
	sc := e.getSerialTransferControl()
	sb := e.getSerialTransferData()

	// Transfer Start Flag (0=No transfer is in progress or requested, 1=Transfer in progress, or requested)
	transferStartFlag := GetBit(sc, 7)
	// Clock Speed (0=Normal, 1=Fast) ** CGB Mode Only **
	//clockSpeed := GetBit(sc, 1)
	// Shift Clock (0=External Clock, 1=Internal Clock)
	shiftClock := GetBit(sc, 0)

	// In Non-CGB Mode the Game Boy supplies an internal clock of 8192Hz
	// 4194304Hz / 8192Hz = 512
	clockFrequency := uint16(512)
	increases := e.timer.numDetectFallingEdges(cyclesElapsed, clockFrequency)
	if shiftClock {
		// Master mode
		if transferStartFlag && increases > 0 {
			e.setSerialTransferData(sb << increases)
			e.serialTransferedBits += increases
		}
	} else {
		// Slave mode
		if increases > 0 {
			//e.setSerialTransferData(sb << increases)
			//e.serialTransferedBits += increases
		}
	}

	transferCompleted := false
	if e.serialTransferedBits >= 8 {
		e.serialTransferedBits = 0
		transferCompleted = true
	}

	if transferCompleted {
		// Clear bit 7 of SC
		e.setSerialTransferControl(SetBit8(sc, 7, false))
		e.mem.requestInterruptSerial()
	}
}
