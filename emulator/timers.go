package emulator

type TimerControl struct {
	TimerEnable    bool
	ClockFrequency int
}

func (e *Emulator) SetInternalDIV(value uint32) {
	e.internalDivider = value
}

// SetDIV Set Internal divider counter
func (e *Emulator) SetDIV(value uint16) {
	e.io[260] = uint8(value)
	e.divTimer = value
}

// GetDIV Get Internal Divider Register
func (e *Emulator) GetDIV() uint8 {
	return e.io[260]
}

// GetTIMA Get Internal Timer Counter
func (e *Emulator) GetTIMA() uint8 {
	return e.io[261]
}

// SetTIMA Set Internal Timer Counter 0xFF05
func (e *Emulator) SetTIMA(value int) {
	timerControl := e.GetTAC()

	tima := value / timerControl.ClockFrequency
	if tima > 0xFF {
		tima = int(e.GetTMA())
		e.SetInterruptTimer()
	}
	e.timaTimer = value

	e.io[261] = uint8(tima)
}

// GetTMA Timer Modulo 0xFF06
func (e *Emulator) GetTMA() uint8 {
	return e.io[262]
}

// GetTAC Timer Control 0xFF07
func (e *Emulator) GetTAC() TimerControl {
	tac := e.io[263]
	timerControl := TimerControl{
		TimerEnable: GetBit(tac, 2),
	}
	switch tac & 0x3 {
	case 0:
		timerControl.ClockFrequency = 1024
	case 1:
		timerControl.ClockFrequency = 16
	case 2:
		timerControl.ClockFrequency = 64
	case 3:
		timerControl.ClockFrequency = 256
	}

	return timerControl
}

func (e *Emulator) incrementTimers() {
	cyclesElapsed := e.cycles - e.prevCycles

	// Increment DIV every 16384Hz
	// 1 Cpu Cycle last 4194304Hz => 4194304Hz / 16384Hz = 256 (Increment DIV every 256 cycles)
	e.divTimer += uint16(cyclesElapsed)
	e.SetDIV(e.divTimer)

	// Increment TIMA
	timerControl := e.GetTAC()
	if timerControl.TimerEnable {
		e.timaTimer += int(cyclesElapsed)
		e.SetTIMA(e.timaTimer)
	}
}
