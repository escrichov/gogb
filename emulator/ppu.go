package emulator

const (
	lcdMode2Bounds = 80
	lcdMode3Bounds = lcdMode2Bounds + 172
)

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

func (e *Emulator) SetLY(value uint8) {
	e.io[324] = value
}

func (e *Emulator) GetLY() uint8 {
	return e.io[324]
}

func (e *Emulator) GetLYC() uint8 {
	return e.io[325]
}

func (e *Emulator) GetLCDC() *LCDControl {
	return &e.lcdcControl
}

func (e *Emulator) SetLCDC(value uint8) {
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

func (e *Emulator) SaveLCDC() {
	e.io[320] = (BoolToUint8(e.lcdcControl.LCDPPUEnable) << 7) |
		(BoolToUint8(e.lcdcControl.WindowTileMapArea) << 6) |
		(BoolToUint8(e.lcdcControl.WindowEnable) << 5) |
		(BoolToUint8(e.lcdcControl.BgWindowTileDataArea) << 4) |
		(BoolToUint8(e.lcdcControl.BgTileMapArea) << 3) |
		(BoolToUint8(e.lcdcControl.ObjSize) << 2) |
		(BoolToUint8(e.lcdcControl.ObjEnable) << 1) |
		BoolToUint8(e.lcdcControl.BgWindowEnablePriority)
}

func (e *Emulator) GetLCDStatus() *LCDStatus {
	return &e.lcdStatus
}

func (e *Emulator) SetLCDStatus(value uint8) {
	e.io[321] = value
	e.lcdStatus.LYCLYSTATInterruptSource = GetBit(value, 6)
	e.lcdStatus.Mode2OAMSTATInterruptSource = GetBit(value, 5)
	e.lcdStatus.Mode1VBlankSTATInterruptSource = GetBit(value, 4)
	e.lcdStatus.Mode0HBlankSTATInterruptSource = GetBit(value, 3)
	e.lcdStatus.LYCLYFlag = GetBit(value, 2)
	e.lcdStatus.ModeFlag = value & 0x3
}

func (e *Emulator) SaveLCDStatus() {
	e.io[321] = (BoolToUint8(e.lcdStatus.LYCLYSTATInterruptSource) << 6) |
		(BoolToUint8(e.lcdStatus.Mode2OAMSTATInterruptSource) << 5) |
		(BoolToUint8(e.lcdStatus.Mode1VBlankSTATInterruptSource) << 4) |
		(BoolToUint8(e.lcdStatus.Mode0HBlankSTATInterruptSource) << 3) |
		(BoolToUint8(e.lcdStatus.LYCLYFlag) << 2) |
		e.lcdStatus.ModeFlag
}

func (e *Emulator) getColor(tile, yOffset, xOffset int) uint8 {
	videoRamIndex := tile*16 + yOffset*2
	tileData := e.videoRam[videoRamIndex]
	tileData1 := e.videoRam[videoRamIndex+1]
	return ((tileData1>>xOffset)%2)*2 + (tileData>>xOffset)%2
}

// Set the status of the LCD based on the current state of memory.
func (e *Emulator) setLCDStatus() {
	status := e.GetLCDStatus()

	if !e.lcdcControl.LCDPPUEnable {
		status.LYCLYSTATInterruptSource = false
		status.Mode2OAMSTATInterruptSource = false
		status.Mode1VBlankSTATInterruptSource = false
		status.Mode0HBlankSTATInterruptSource = false
		status.LYCLYFlag = false
		status.ModeFlag = 0
	}

	ly := e.GetLY()
	currentMode := status.ModeFlag

	var mode uint8
	requestInterrupt := false

	switch {
	case ly >= 144:
		mode = 1
		requestInterrupt = status.Mode1VBlankSTATInterruptSource
	case e.ppuDot < lcdMode2Bounds:
		mode = 2
		requestInterrupt = status.Mode2OAMSTATInterruptSource
	case e.ppuDot < lcdMode3Bounds:
		mode = 3
	default:
		mode = 0
		requestInterrupt = status.Mode0HBlankSTATInterruptSource
		if mode != currentMode {
			//gb.Memory.doHDMATransfer()
		}
	}

	if requestInterrupt && mode != currentMode {
		e.requestInterruptLCDStat()
	}

	// Check if LYC == LY (coincidence flag)
	lyc := e.GetLYC()
	if ly == lyc {
		e.lcdStatus.LYCLYFlag = true
		if e.lcdStatus.LYCLYSTATInterruptSource {
			e.requestInterruptLCDStat()
		}
	} else {
		e.lcdStatus.LYCLYFlag = false
	}

	e.SaveLCDStatus()
}

func (e *Emulator) PPURun() bool {
	renderFrame := false
	e.setLCDStatus()

	// PPU
	cyclesElapsed := e.cycles - e.prevCycles
	for i := uint64(0); i < cyclesElapsed; i++ {
		lcdc := e.GetLCDC()
		if lcdc.LCDPPUEnable {
			e.ppuDot++

			// Render Scanline (Every 256 PPU Dots)
			if e.ppuDot == 456 {
				ly := e.GetLY()

				// Only render visible lines (up to line 144)
				if ly < HEIGHT {
					for tmp := WIDTH - 1; tmp >= 0; tmp-- {

						// IsWindow
						isWindow := false
						if lcdc.WindowEnable && ly >= e.io[330] && uint8(tmp) >= (e.io[331]-7) {
							isWindow = true
						}

						// xOffset
						var xOffset uint8
						if isWindow {
							xOffset = uint8(tmp) - e.io[331] + 7
						} else {
							xOffset = uint8(tmp) + e.io[323]
						}

						// yOffset
						var yOffset uint8
						if isWindow {
							yOffset = ly - e.io[330]
						} else {
							yOffset = ly + e.io[322]
						}

						// PaletteIndex
						var paletteIndex uint16 = 0

						// Tile
						tileMapArea := lcdc.BgTileMapArea
						if isWindow {
							tileMapArea = lcdc.WindowTileMapArea
						}

						videoRamIndex := uint16(6)
						if tileMapArea {
							videoRamIndex = 7
						}
						videoRamIndex = videoRamIndex<<10 | uint16(yOffset)/8*32 + uint16(xOffset)/8
						var tile = e.videoRam[videoRamIndex]

						// Color
						var tileValue int
						if lcdc.BgWindowTileDataArea {
							tileValue = int(tile)
						} else {
							tileValue = 256 + int(int8(tile))
						}
						color := e.getColor(tileValue, int(yOffset&7), int(7-xOffset&7))

						// Sprites
						if lcdc.ObjEnable {
							for spriteIndex := uint8(0); spriteIndex < WIDTH; spriteIndex += 4 {
								spriteX := uint8(tmp) - e.io[spriteIndex+1] + 8
								spriteY := ly - e.io[spriteIndex] + 16

								spriteYOffset := uint8(0)
								if (e.io[spriteIndex+3] & 64) != 0 {
									spriteYOffset = 7
								}
								spriteYOffset = spriteY ^ spriteYOffset

								spriteXOffset := uint8(7)
								if (e.io[spriteIndex+3] & 32) != 0 {
									spriteXOffset = 0
								}
								spriteXOffset = spriteX ^ spriteXOffset

								spriteColor := e.getColor(int(e.io[spriteIndex+2]), int(spriteYOffset), int(spriteXOffset))

								if spriteX < 8 && spriteY < 8 && !((e.io[spriteIndex+3]&128) != 0 && color != 0) && spriteColor != 0 {
									color = spriteColor
									if e.io[spriteIndex+3]&16 == 0 {
										paletteIndex = uint16(1)
									} else {
										paletteIndex = uint16(2)
									}
									break
								}
							}
						}

						paletteIndexValue := uint16((e.io[327+paletteIndex]>>(2*color))%4) + paletteIndex*4&7
						frameBufferIndex := uint16(ly)*WIDTH + uint16(tmp)
						e.frameBuffer[frameBufferIndex] = e.palette[paletteIndexValue]
					}
				}

				if ly == (HEIGHT - 1) {
					e.requestInterruptVBlank()
					renderFrame = true
				}

				// Increment Line
				e.SetLY((ly + 1) % 154)
				e.ppuDot = 0
			}
		} else {
			e.SetLY(0)
			e.ppuDot = 0
		}
	}

	return renderFrame
}
