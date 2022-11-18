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

func (mem *Memory) SetLY(value uint8) {
	mem.io[324] = value
}

func (mem *Memory) GetLY() uint8 {
	return mem.io[324]
}

func (mem *Memory) GetLYC() uint8 {
	return mem.io[325]
}

func (mem *Memory) GetLCDC() *LCDControl {
	return &mem.lcdcControl
}

func (mem *Memory) GetWY() uint8 {
	return mem.io[330]
}

func (mem *Memory) GetWX() uint8 {
	return mem.io[331]
}

func (mem *Memory) GetSCY() uint8 {
	return mem.io[322]
}

func (mem *Memory) GetSCX() uint8 {
	return mem.io[323]
}

func (mem *Memory) GetBGP() uint8 {
	return mem.io[328]
}

func (mem *Memory) GetOBP0() uint8 {
	return mem.io[328]
}

func (mem *Memory) GetOBP1() uint8 {
	return mem.io[329]
}

func (mem *Memory) SetLCDC(value uint8) {
	mem.io[320] = value
	mem.lcdcControl.LCDPPUEnable = GetBit(value, 7)
	mem.lcdcControl.WindowTileMapArea = GetBit(value, 6)
	mem.lcdcControl.WindowEnable = GetBit(value, 5)
	mem.lcdcControl.BgWindowTileDataArea = GetBit(value, 4)
	mem.lcdcControl.BgTileMapArea = GetBit(value, 3)
	mem.lcdcControl.ObjSize = GetBit(value, 2)
	mem.lcdcControl.ObjEnable = GetBit(value, 1)
	mem.lcdcControl.BgWindowEnablePriority = GetBit(value, 0)
}

func (mem *Memory) SaveLCDC() {
	mem.io[320] = (BoolToUint8(mem.lcdcControl.LCDPPUEnable) << 7) |
		(BoolToUint8(mem.lcdcControl.WindowTileMapArea) << 6) |
		(BoolToUint8(mem.lcdcControl.WindowEnable) << 5) |
		(BoolToUint8(mem.lcdcControl.BgWindowTileDataArea) << 4) |
		(BoolToUint8(mem.lcdcControl.BgTileMapArea) << 3) |
		(BoolToUint8(mem.lcdcControl.ObjSize) << 2) |
		(BoolToUint8(mem.lcdcControl.ObjEnable) << 1) |
		BoolToUint8(mem.lcdcControl.BgWindowEnablePriority)
}

func (mem *Memory) GetLCDStatus() *LCDStatus {
	return &mem.lcdStatus
}

func (mem *Memory) SetLCDStatus(value uint8) {
	mem.io[321] = value
	mem.lcdStatus.LYCLYSTATInterruptSource = GetBit(value, 6)
	mem.lcdStatus.Mode2OAMSTATInterruptSource = GetBit(value, 5)
	mem.lcdStatus.Mode1VBlankSTATInterruptSource = GetBit(value, 4)
	mem.lcdStatus.Mode0HBlankSTATInterruptSource = GetBit(value, 3)
	mem.lcdStatus.LYCLYFlag = GetBit(value, 2)
	mem.lcdStatus.ModeFlag = value & 0x3
}

func (mem *Memory) SaveLCDStatus() {
	mem.io[321] = (BoolToUint8(mem.lcdStatus.LYCLYSTATInterruptSource) << 6) |
		(BoolToUint8(mem.lcdStatus.Mode2OAMSTATInterruptSource) << 5) |
		(BoolToUint8(mem.lcdStatus.Mode1VBlankSTATInterruptSource) << 4) |
		(BoolToUint8(mem.lcdStatus.Mode0HBlankSTATInterruptSource) << 3) |
		(BoolToUint8(mem.lcdStatus.LYCLYFlag) << 2) |
		mem.lcdStatus.ModeFlag
}
