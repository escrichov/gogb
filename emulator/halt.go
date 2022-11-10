package emulator

func (e *Emulator) HaltRun() {
	e.tick()

	// If IME is not set, there are two distinct cases,
	// depending on whether an interrupt is pending as the halt instruction is first executed.
	if e.IME == 0 {
		if e.isInterruptPendingInFirstHaltExecution {
			e.isHaltBugActive = true
		}
	} else if e.IME == 1 {
		if e.isInterruptPendingInFirstHaltExecution && e.isIMEDelayedInFirstHaltExecution {
			e.isHaltBugEIActive = true
		}
	}

	if e.hasPendingInterrupts() {
		e.halt = 0
	}
}
