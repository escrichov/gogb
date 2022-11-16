package emulator

type TimerControl struct {
	TimerEnable    bool
	ClockFrequency uint16
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

// GetDIV Get Internal Divider Register
func (m *Memory) GetDIV() uint8 {
	return m.io[260]
}

// SetDIV Set divider counter
func (m *Memory) SetDIV(value uint8) {
	m.io[260] = value
}

// GetTIMA Get Internal Timer Counter
func (m *Memory) GetTIMA() uint8 {
	return m.io[261]
}

// SetTIMA Set Internal Timer Counter
func (m *Memory) SetTIMA(value uint8) {
	m.io[261] = value
}

// GetTMA Get Timer Modulo 0xFF06
func (m *Memory) GetTMA() uint8 {
	return m.io[262]
}

// SetTMA Set Timer Modulo 0xFF06
func (m *Memory) SetTMA(value uint8) {
	m.io[262] = value
}

// GetTAC Timer Control 0xFF07
func (m *Memory) GetTAC() TimerControl {
	tac := m.io[263]
	return parseTimerControl(tac)
}
