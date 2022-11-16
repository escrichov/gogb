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

func (e *Memory) SetLY(value uint8) {
	e.io[324] = value
}

func (e *Memory) GetLY() uint8 {
	return e.io[324]
}

func (e *Memory) GetLYC() uint8 {
	return e.io[325]
}

func (e *Memory) GetLCDC() *LCDControl {
	return &e.lcdcControl
}

func (e *Memory) SetLCDC(value uint8) {
	e.io[320] = value
	e.lcdcControl.LCDPPUEnable = GetBit(value, 7)
	e.lcdcControl.WindowTileMapArea = GetBit(value, 6)
	e.lcdcControl.WindowEnable = GetBit(value, 5)
	e.lcdcControl.BgWindowTileDataArea = GetBit(value, 4)
	e.lcdcControl.BgTileMapArea = GetBit(value, 3)
	e.lcdcControl.ObjSize = GetBit(value, 2)
	e.lcdcControl.ObjEnable = GetBit(value, 1)
	e.lcdcControl.BgWindowEnablePriority = GetBit(value, 0)
}

func (e *Memory) SaveLCDC() {
	e.io[320] = (BoolToUint8(e.lcdcControl.LCDPPUEnable) << 7) |
		(BoolToUint8(e.lcdcControl.WindowTileMapArea) << 6) |
		(BoolToUint8(e.lcdcControl.WindowEnable) << 5) |
		(BoolToUint8(e.lcdcControl.BgWindowTileDataArea) << 4) |
		(BoolToUint8(e.lcdcControl.BgTileMapArea) << 3) |
		(BoolToUint8(e.lcdcControl.ObjSize) << 2) |
		(BoolToUint8(e.lcdcControl.ObjEnable) << 1) |
		BoolToUint8(e.lcdcControl.BgWindowEnablePriority)
}

func (e *Memory) GetLCDStatus() *LCDStatus {
	return &e.lcdStatus
}

func (e *Memory) SetLCDStatus(value uint8) {
	e.io[321] = value
	e.lcdStatus.LYCLYSTATInterruptSource = GetBit(value, 6)
	e.lcdStatus.Mode2OAMSTATInterruptSource = GetBit(value, 5)
	e.lcdStatus.Mode1VBlankSTATInterruptSource = GetBit(value, 4)
	e.lcdStatus.Mode0HBlankSTATInterruptSource = GetBit(value, 3)
	e.lcdStatus.LYCLYFlag = GetBit(value, 2)
	e.lcdStatus.ModeFlag = value & 0x3
}

func (e *Memory) SaveLCDStatus() {
	e.io[321] = (BoolToUint8(e.lcdStatus.LYCLYSTATInterruptSource) << 6) |
		(BoolToUint8(e.lcdStatus.Mode2OAMSTATInterruptSource) << 5) |
		(BoolToUint8(e.lcdStatus.Mode1VBlankSTATInterruptSource) << 4) |
		(BoolToUint8(e.lcdStatus.Mode0HBlankSTATInterruptSource) << 3) |
		(BoolToUint8(e.lcdStatus.LYCLYFlag) << 2) |
		e.lcdStatus.ModeFlag
}
