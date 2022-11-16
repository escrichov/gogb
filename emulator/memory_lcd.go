package emulator

type LCDControl struct {
	LCDPPUEnable           bool // Bit 7, 0=Off, 1=On
	WindowTileMapArea      bool // Bit 6, 0=9800-9BFF, 1=9C00-9FFF
	WindowEnable           bool // Bit 5, 0=Off, 1=On
	BgWindowTileDataArea   bool // Bit 4, 0=8800-97FF, 1=8000-8FFF
	BgTileMapArea          bool // Bit 3, 0=9800-9BFF, 1=9C00-9FFF
	ObjSize                bool // Bit 2, 0=8x8, 1=8x16
	ObjEnable              bool // Bit 1, 0=Off, 1=On
	BgWindowEnablePriority bool // Bit 0, 0=Off, 1=On
}

type LCDStatus struct {
	LYCLYSTATInterruptSource       bool  // Bit 6, 0=Off, 1=Enable
	Mode2OAMSTATInterruptSource    bool  // Bit 5, 0=Off, 1=Enable
	Mode1VBlankSTATInterruptSource bool  // Bit 4, 0=Off, 1=Enable
	Mode0HBlankSTATInterruptSource bool  // Bit 3, 0=Off, 1=Enable
	LYCLYFlag                      bool  // Bit 2, 0=Different, 1=Equal
	ModeFlag                       uint8 // Bit 1-0, 0: HBlank, 1: VBlank, 2: Searching OAM, 3: Transferring Data to LCD Controller
}

func (m *Memory) SetLY(value uint8) {
	m.io[324] = value
}

func (m *Memory) GetLY() uint8 {
	return m.io[324]
}

func (m *Memory) GetLYC() uint8 {
	return m.io[325]
}

func (m *Memory) GetLCDC() *LCDControl {
	return &m.lcdcControl
}

func (m *Memory) SetLCDC(value uint8) {
	m.io[320] = value
	m.lcdcControl.LCDPPUEnable = GetBit(value, 7)
	m.lcdcControl.WindowTileMapArea = GetBit(value, 6)
	m.lcdcControl.WindowEnable = GetBit(value, 5)
	m.lcdcControl.BgWindowTileDataArea = GetBit(value, 4)
	m.lcdcControl.BgTileMapArea = GetBit(value, 3)
	m.lcdcControl.ObjSize = GetBit(value, 2)
	m.lcdcControl.ObjEnable = GetBit(value, 1)
	m.lcdcControl.BgWindowEnablePriority = GetBit(value, 0)
}

func (m *Memory) SaveLCDC() {
	m.io[320] = (BoolToUint8(m.lcdcControl.LCDPPUEnable) << 7) |
		(BoolToUint8(m.lcdcControl.WindowTileMapArea) << 6) |
		(BoolToUint8(m.lcdcControl.WindowEnable) << 5) |
		(BoolToUint8(m.lcdcControl.BgWindowTileDataArea) << 4) |
		(BoolToUint8(m.lcdcControl.BgTileMapArea) << 3) |
		(BoolToUint8(m.lcdcControl.ObjSize) << 2) |
		(BoolToUint8(m.lcdcControl.ObjEnable) << 1) |
		BoolToUint8(m.lcdcControl.BgWindowEnablePriority)
}

func (m *Memory) GetLCDStatus() *LCDStatus {
	return &m.lcdStatus
}

func (m *Memory) SetLCDStatus(value uint8) {
	m.io[321] = value
	m.lcdStatus.LYCLYSTATInterruptSource = GetBit(value, 6)
	m.lcdStatus.Mode2OAMSTATInterruptSource = GetBit(value, 5)
	m.lcdStatus.Mode1VBlankSTATInterruptSource = GetBit(value, 4)
	m.lcdStatus.Mode0HBlankSTATInterruptSource = GetBit(value, 3)
	m.lcdStatus.LYCLYFlag = GetBit(value, 2)
	m.lcdStatus.ModeFlag = value & 0x3
}

func (m *Memory) SaveLCDStatus() {
	m.io[321] = (BoolToUint8(m.lcdStatus.LYCLYSTATInterruptSource) << 6) |
		(BoolToUint8(m.lcdStatus.Mode2OAMSTATInterruptSource) << 5) |
		(BoolToUint8(m.lcdStatus.Mode1VBlankSTATInterruptSource) << 4) |
		(BoolToUint8(m.lcdStatus.Mode0HBlankSTATInterruptSource) << 3) |
		(BoolToUint8(m.lcdStatus.LYCLYFlag) << 2) |
		m.lcdStatus.ModeFlag
}
