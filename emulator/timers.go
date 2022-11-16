package emulator

type Timer struct {
	internalTimer                  uint16
	timaUpdateWithTMADelayedCycles uint64

	emulator *Emulator
}

// SetInternalTimer Set Internal Timer
func (timer *Timer) SetInternalTimer(value uint16) {
	timer.internalTimer = value
}

// updateDiv update divider counter with internal timer
func (timer *Timer) updateDiv(cyclesElapsed uint16) {
	// Increment DIV every 16384Hz
	// 1 Cpu Cycle takes 4194304Hz => 4194304Hz / 16384Hz = 256 (Increment DIV every 256 cycles)
	newTimer := timer.internalTimer + cyclesElapsed
	div := uint8(newTimer >> 8)
	timer.emulator.mem.SetDIV(div)
}

func (timer *Timer) numDetectFallingEdges(cyclesElapsed, clockFrequency uint16) uint8 {
	oldTimer := timer.internalTimer
	newTimer := oldTimer + cyclesElapsed

	value := oldTimer / clockFrequency
	newValue := newTimer / clockFrequency
	return uint8(newValue - value)
}

func (timer *Timer) isFallingEdgeWritingDIV() bool {
	// When writing to DIV, if the current output is 1 and timer is enabled,
	// as the new value after reseting DIV will be 0,
	// the falling edge detector will detect a falling edge and TIMA will increase.
	timerControl := timer.emulator.mem.GetTAC()
	if timerControl.TimerEnable {
		if timerControl.ClockFrequency == 1024 && GetBit16(timer.internalTimer, 9) {
			return true
		} else if timerControl.ClockFrequency == 16 && GetBit(uint8(timer.internalTimer), 3) {
			return true
		} else if timerControl.ClockFrequency == 64 && GetBit(uint8(timer.internalTimer), 5) {
			return true
		} else if timerControl.ClockFrequency == 256 && GetBit(uint8(timer.internalTimer), 7) {
			return true
		}
	}

	return false
}

func (timer *Timer) isFallingEdgeWritingTAC(newTac uint8) bool {
	// When writing to DIV, if the current output is 1 and timer is enabled,
	// as the new value after reseting DIV will be 0,
	// the falling edge detector will detect a falling edge and TIMA will increase.
	oldTimerControl := timer.emulator.mem.GetTAC()
	newTimerControl := parseTimerControl(newTac)
	if oldTimerControl.TimerEnable {
		if newTimerControl.TimerEnable {
			return ((timer.internalTimer & (oldTimerControl.ClockFrequency / 2)) != 0) && ((timer.internalTimer & (newTimerControl.ClockFrequency / 2)) == 0)
		} else {
			return (timer.internalTimer & (oldTimerControl.ClockFrequency / 2)) != 0
		}
	}

	return false
}

// UpdateTIMA Update TIMA with internal counter
func (timer *Timer) updateTIMA(cyclesElapsed uint16) {
	timerControl := timer.emulator.mem.GetTAC()

	if timerControl.TimerEnable {
		increases := timer.numDetectFallingEdges(cyclesElapsed, timerControl.ClockFrequency)

		timer.increaseTIMA(increases, true)
	}
}

func (timer *Timer) increaseTIMA(increases uint8, delayTimaReloadAndInterrupt bool) {
	tima := timer.emulator.mem.GetTIMA()
	newTima := tima + increases
	if uint16(tima)+uint16(increases) > 0xFF {
		// When TIMA overflows, the value from TMA is loaded and IF timer flag is set to 1,
		// but this doesn't happen immediately.
		// Timer interrupt is delayed 1 cycle (4 clocks) from the TIMA overflow.
		// The TMA reload to TIMA is also delayed. For one cycle, after overflowing TIMA,
		// the value in TIMA is 00h, not TMA.
		if delayTimaReloadAndInterrupt {
			newTima = 0
			timer.timaUpdateWithTMADelayedCycles = timer.emulator.cycles + 4*4
		} else {
			newTima = timer.emulator.mem.GetTMA()
			timer.emulator.mem.requestInterruptTimer()
		}
	}
	timer.emulator.mem.SetTIMA(newTima)
}

func (timer *Timer) reloadTIMAwithTMA() {
	newTima := timer.emulator.mem.GetTMA()
	timer.emulator.mem.SetTIMA(newTima)
	timer.emulator.mem.requestInterruptTimer()
}

func (timer *Timer) incrementTimers(cyclesElapsed uint16) {
	// Increment DIV
	timer.updateDiv(cyclesElapsed)

	// Increment TIMA
	timer.updateTIMA(cyclesElapsed)

	// Update internal timer
	timer.internalTimer += cyclesElapsed
}
