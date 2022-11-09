package emulator

type TimerControl struct {
	TimerEnable    bool
	ClockFrequency uint16
}

// GetDIV Get Internal Divider Register
func (e *Emulator) GetDIV() uint8 {
	return e.io[260]
}

// SetDIV Set divider counter
func (e *Emulator) SetDIV(value uint8) {
	e.io[260] = value
}

// GetTIMA Get Internal Timer Counter
func (e *Emulator) GetTIMA() uint8 {
	return e.io[261]
}

// SetTIMA Set Internal Timer Counter
func (e *Emulator) SetTIMA(value uint8) {
	e.io[261] = value
}

// GetTMA Get Timer Modulo 0xFF06
func (e *Emulator) GetTMA() uint8 {
	return e.io[262]
}

// SetTMA Set Timer Modulo 0xFF06
func (e *Emulator) SetTMA(value uint8) {
	e.io[262] = value
}

// GetTAC Timer Control 0xFF07
func (e *Emulator) GetTAC() TimerControl {
	tac := e.io[263]
	return parseTimerControl(tac)
}

func parseTimerControl(tac uint8) TimerControl {
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

// SetInternalCounter Set Internal Timer
func (e *Emulator) SetInternalTimer(value uint16) {
	e.internalTimer = value
}

// updateDiv update divider counter with internal timer
func (e *Emulator) updateDiv(cyclesElapsed uint16) {
	// Increment DIV every 16384Hz
	// 1 Cpu Cycle takes 4194304Hz => 4194304Hz / 16384Hz = 256 (Increment DIV every 256 cycles)
	newTimer := e.internalTimer + cyclesElapsed
	div := uint8(newTimer >> 8)
	e.SetDIV(div)
}

func numDetectFallingEdges(oldTimer uint16, newTimer uint16, clockFrequency uint16) uint8 {
	value := oldTimer / clockFrequency
	newValue := newTimer / clockFrequency
	return uint8(newValue - value)
}

func (e *Emulator) isFallingEdgeWritingDIV() bool {
	// When writing to DIV, if the current output is 1 and timer is enabled,
	// as the new value after reseting DIV will be 0,
	// the falling edge detector will detect a falling edge and TIMA will increase.
	timerControl := e.GetTAC()
	if timerControl.TimerEnable {
		if timerControl.ClockFrequency == 1024 && GetBit16(e.internalTimer, 9) {
			return true
		} else if timerControl.ClockFrequency == 16 && GetBit(uint8(e.internalTimer), 3) {
			return true
		} else if timerControl.ClockFrequency == 64 && GetBit(uint8(e.internalTimer), 5) {
			return true
		} else if timerControl.ClockFrequency == 256 && GetBit(uint8(e.internalTimer), 7) {
			return true
		}
	}

	return false
}

func (e *Emulator) isFallingEdgeWritingTAC(newTac uint8) bool {
	// When writing to DIV, if the current output is 1 and timer is enabled,
	// as the new value after reseting DIV will be 0,
	// the falling edge detector will detect a falling edge and TIMA will increase.
	oldTimerControl := e.GetTAC()
	newTimerControl := parseTimerControl(newTac)
	if oldTimerControl.TimerEnable {
		if newTimerControl.TimerEnable {
			return ((e.internalTimer & (oldTimerControl.ClockFrequency / 2)) != 0) && ((e.internalTimer & (newTimerControl.ClockFrequency / 2)) == 0)
		} else {
			return (e.internalTimer & (oldTimerControl.ClockFrequency / 2)) != 0
		}
	}

	return false
}

// UpdateTIMA Update TIMA with internal counter
func (e *Emulator) updateTIMA(cyclesElapsed uint16) {
	timerControl := e.GetTAC()

	if timerControl.TimerEnable {
		increases := numDetectFallingEdges(e.internalTimer, e.internalTimer+cyclesElapsed, timerControl.ClockFrequency)

		e.increaseTIMA(increases)
	}
}

func (e *Emulator) increaseTIMA(increases uint8) {
	tima := e.GetTIMA()
	newTima := tima + increases
	if uint16(tima)+uint16(increases) > 0xFF {
		newTima = e.GetTMA()
		e.SetInterruptTimer()
	}
	e.SetTIMA(newTima)
}

func (e *Emulator) incrementTimers() {
	cyclesElapsed := uint16(e.cycles - e.prevCycles)

	// Increment DIV
	e.updateDiv(cyclesElapsed)

	// Increment TIMA
	e.updateTIMA(cyclesElapsed)

	// Update internal timer
	e.internalTimer += cyclesElapsed
}
