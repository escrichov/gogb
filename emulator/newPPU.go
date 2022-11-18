package emulator

import "fmt"

// Summary
//
// - Tiles
//   * 8x8 color IDs
//   * color ID = 2 bits per pixel = 0 to 3
//   * color IDs are associated with a palette
//   * In Objects ID 0 means transparent
//
// - Palettes
//   * Array of 4 colors
//
// - Background
//   * is composed of a tilemap
//   * Tilemap contain references to the tiles
//   * color IDs are associated with a palette
//   * In Objects ID 0 means transparent
//
// - Window
//   * no transparency
//   * it’s always a rectangle
//   * only the position of the top-left pixel can be controlled
//
// - Objects
//   * an entry in object attribute memory (OAM)
//   * objects are made of 1 or 2 stacked tiles (8x8 or 8x16 pixels)
//   * it’s always a rectangle
//   * only the position of the top-left pixel can be controlled
//
// - VRAM Tile Data
//   * $8000-$97FF
//   * 16 bytes per tile
//   * 384 tiles
//   * Each tile: 8x8 pixels and has a color depth of 4 colors/gray shades
//   * There are three “blocks” of 128 tiles each
//      * 0: $8000–$87FF, Objs: 0–127, BG/Win if LCDC.4=1: 0–127
//      * 1: $8800–$8FFF, Objs: 128–255, BG/Win if LCDC.4=1: 128–255, BG/Win if LCDC.4=0: 128–255 (or -128–-1)
//      * 2: $9000–$97FF, Objs: Can't use, BG/Win if LCDC.4=0: 0–127
//   * Tiles are always indexed using an 8-bit integer
//   * Each tile occupies 16 bytes, where each line is represented by 2 bytes:
//      * Byte 0-1  Topmost Line (Top 8 pixels)
//      * Byte 2-3  Second Line
//      * etc.
//      * For each line:
//         * the first byte specifies the least significant bit of the color ID
//         * and the second byte specifies the most significant bit
//
// - Tile Maps
//   * two 32x32 tile maps: $9800-$9BFF and $9C00-$9FFF
//   * 1-byte indexes of the tiles
//   * Since one tile has 8x8 pixels, each map holds a 256x256 pixels picture
//   * Only 160x144 of those pixels are displayed on the LCD
//   * In Non-CGB mode, the Background (and the Window) can be disabled using LCDC bit 0
//   * Display the Background or the Window
//
// - Tile Maps Background
//   * The SCY and SCX registers can be used to scroll the Background

// - Tile Maps Window
//   * The top left corner of the Window are (WX-7,WY)
//   * it is always displayed starting at the top left tile of its tile map.
//   * is displayed is both: LCDC bit 5 and 0 are set
//   * Enabling the Window makes PPU Mode 3 slightly longer
//   * Window keeps an internal line counter that’s functionally similar to LY
//
// - SCY, SCX: Viewport Y position, X position
//   * Top-left coordinates of the visible 160×144 pixel area within the 256×256 pixels BG map
//   * Values in the range 0–255 may be used.
//   * The visible area of the Background wraps around the Background map
//     (that is, when part of the visible area goes beyond the map edge,
//      it starts displaying the opposite side of the map).
//
// - Window WY, WX
//   * The top left corner of the Window are (WX-7,WY)
//   * The Window is visible (if enabled) when both coordinates are in the ranges WX=0..166, WY=0..143
//   * Values WX=7, WY=0 place the Window at the top left of the screen, completely covering the background.
//   * WX values 0 and 166 are unreliable due to hardware bugs.
//      * If WX is set to 0, the window will “stutter” horizontally when SCX changes (depending on SCX % 8).
//      * If WX is set to 166, the window will span the entirety of the following scanline.
//
// - LY
//   * LY indicates the current horizontal line
//   * from 0 to 153
//   * 144 to 153 indicating the VBlank
//
// - LYC: LY compare
//   * When both values are identical, the “LYC=LY” flag in the STAT register is set,
//     and (if enabled) a STAT interrupt is requested.
//
// - Mid-frame scrolling behavior
//   * Scrolling
//     * Scroll registers are re-read on each tile fetch except for the low 3 bits of SCX
//     * low 3 bits of SCX are only read at the beginning of the scanline
//     * Y coordinate is read once for each bitplane (What is this?)
//   * Window
//     * WY == LY is checked at the start of Mode 2 only
//     * current X + 7 == WX
//     * If WY == LY && window enable bit = 1 at the start of the row
//         => when WX triggered && window enable bit = 0 => glitch pixel
//
// - VRAM Sprite Attribute Table OAM
//   * $FE00-FE9F
//   * 40 entries consists of 4 bytes each
//   * Can display up to 40 sprites either in 8x8 or in 8x16 pixels
//   * Only 10 sprites can be displayed per scan line
//   * 144 to 153 indicating the VBlank
//   * Two ways to write to OAM
//     * During the HBlank and VBlank periods.
//     * write the data to a buffer in normal RAM (typically WRAM) first,
//       then to copy that buffer to OAM using the DMA transfer functionality
//   * Drawing priority
//     * In Non-CGB mode: the smaller the X coordinate, the higher the priority. When X coordinates are identical, the object located first in OAM has higher priority.
//     * In CGB mode: the object located first in OAM has higher priority.
//   * 4 bytes per Entry:
//     * Byte 0 — Y Position => Sprite's vertical position on the screen + 16
//     * Byte 1 — X Position => Sprite's horizontal position on the screen + 8
//     * Byte 2 — Tile Index * In 8x8 mode  (LCDC bit 2 = 0) => tile index ($00-$FF)
//                           * In 8x16 mode (LCDC bit 2 = 1)
//                               => index of the first (top) tile of the sprite
//                               => top 8x8 tile is “NN & $FE” && bottom 8x8 tile is “NN | $01”.
//     * Byte 3 — Attributes/Flags
//       * Bit7   BG and Window over OBJ (0=No, 1=BG and Window colors 1-3 over the OBJ)
//       * Bit6   Y flip          (0=Normal, 1=Vertically mirrored)
//       * Bit5   X flip          (0=Normal, 1=Horizontally mirrored)
//       * Bit4   Palette number  **Non CGB Mode Only** (0=OBP0, 1=OBP1)
//       * Bit3   Tile VRAM-Bank  **CGB Mode Only**     (0=Bank 0, 1=Bank 1)
//       * Bit2-0 Palette number  **CGB Mode Only**     (OBP0-7)
//
// - OAM DMA
//   * Writing to $FF46 launches a DMA transfer from ROM or RAM to OAM
//   * Source:      $XX00-$XX9F   ;XX = $00 to $DF where XX is the value write to $FF46
//   * Destination: $FE00-$FE9F
//   * The transfer takes 160 machine cycles
//   * On DMG, during this time, the CPU can access only HRAM (memory at $FF80-$FFFE)
//
// - FF40 — LCDC: LCD control
//   * 7 LCD and PPU enable - 0=Off, 1=On
//       - Stopping LCD operation (Bit 7 from 1 to 0) may be performed during VBlank ONLY
//       - When the display is disabled the screen is blank, is displayed as a white “whiter” than color #0.
//       - When re-enabling the LCD, the PPU will immediately start drawing again,
//         but the screen will stay blank during the first frame.
//   * 6 Window tile map area - 0=9800-9BFF, 1=9C00-9FFF
//   * 5 Window enable - 0=Off, 1=On
//   * 4 BG and Window tile data area - 0=8800-97FF, 1=8000-8FFF
//   * 3 BG tile map area - 0=9800-9BFF, 1=9C00-9FFF
//   * 2 OBJ size - 0=8x8, 1=8x16
//   * 1 OBJ enable - 0=Off, 1=On
//   * 0 BG and Window enable/priority - 0=Off, 1=On
//       - Non-CGB Mode: 0 = both background and window become blank (white),
//                           and the Window Display Bit is ignored in that case
//       - CGB Mode: 0 = When Bit 0 is cleared,
//                       the background and window lose their priority -
//                       the sprites will always be displayed on top of background and window,
//                       independently of the priority flags in OAM and BG Map attributes.
//
// - FF41 — STAT: LCD status
//   * Bit 6 - LYC=LY STAT Interrupt source         (1=Enable) (Read/Write)
//   * Bit 5 - Mode 2 OAM STAT Interrupt source     (1=Enable) (Read/Write)
//   * Bit 4 - Mode 1 VBlank STAT Interrupt source  (1=Enable) (Read/Write)
//   * Bit 3 - Mode 0 HBlank STAT Interrupt source  (1=Enable) (Read/Write)
//   * Bit 2 - LYC=LY Flag                          (0=Different, 1=Equal) (Read Only)
//   * Bit 1-0 - Mode Flag                          (Mode 0-3, see below) (Read Only)
//          0: HBlank
//          1: VBlank
//          2: Searching OAM
//          3: Transferring Data to LCD Controller
//   * Writing to STAT during OAM scan, HBlank, VBlank or LY=LYC trigger an STAT interrupt
//     It behaves as if $FF were written for one cycle
//     and then the written value were written the next cycle.
//
// - Dots
//   * A dot is the shortest period over which the PPU can output one pixel
//      - is it equivalent to 1 T-state on DMG or 2 T-states on CGB double-speed mode
//      - On each dot during mode 3, either the PPU outputs a pixel or the fetcher is stalling the FIFOs
//   * Frame = 154 scanlines = 70224 dots
//   * On scanlines 0 through 143, the PPU cycles through modes 2, 3, and 0 once every 456 dots.
//   * Scanlines 144 through 153 are mode 1.
//
// - Modes and VRAM Memory Accessible
//    0: HBlank => 85 to 208 dots (Nothing)
//    1: VBlank => 4560 always (Nothing)
//    2: Searching OAM => 80 dots always (Searching OAM for OBJs whose Y coordinate overlap this line)
//    3: Transferring Data to LCD Controller => 168 to 291 dots (Reading OAM and VRAM to generate the picture)
//   * modes 2 and 3, the CPU cannot access OAM ($FE00-FE9F) => return 0xFF
//   * mode 3, the CPU cannot access VRAM ($8000-$9FFF) and CGB palette data registers ($FF69,$FF6B)
//   * mode 3 pauses:
//     - If SCX % 8 is not zero at the start of the scanline
//          => shifter discards that many pixels from the leftmost tile
//          => rendering is paused for that many dots
//     - Window active
//          => pauses for at least 6 dots
//          => as the background fetching mechanism starts over on the left side of the window.
//     - Sprites
//          => each sprite pauses for => 11 - min(5, (x + SCX) % 8) dots
//             Because sprite fetch waits for background fetch to finish,
//             a sprite's cost depends on its position relative to the left side of the background tile under it.
//            It’s greater if a sprite is directly aligned over the background tile, less if the sprite is to the right.
//          => If the sprite’s left side is over the window, => 11 - min(5, (x + (255 - WX)) % 8) dots
//
// - Interrupts
//   * LYC=LY STAT Interrupt source, only if bit 6 of STAT is set (When lyc == ly)
//   * Mode 2 OAM STAT Interrupt source, only if bit 6 of STAT is set (When mode 2 start)
//   * Mode 1 VBlank STAT Interrupt source, only if bit 6 of STAT is set (When mode 1 start)
//   * Mode 0 VBlank STAT Interrupt source, only if bit 6 of STAT is set (When mode 0 start)
//   * VBlank interrupt (When mode 1 start)
//
// - Palettes
//   * FF47 — BGP (Non-CGB Mode only): BG palette data
//     - Bit 7-6 - Color for index 3
//     - Bit 5-4 - Color for index 2
//     - Bit 3-2 - Color for index 1
//     - Bit 1-0 - Color for index 0
//     - 0	White
//     - 1	Light gray
//     - 2	Dark gray
//     - 3	Black
//   * FF48–FF49 — OBP0, OBP1 (Non-CGB Mode only): OBJ palette 0, 1 data
//     - They work exactly like BGP, except that the lower two bits are ignored because color index 0 is transparent for OBJs.

const (
	modeHblank     int = 0
	modeVblank         = 1
	modeOAMScan        = 2
	modeDrawPixels     = 3
)

type PixelFetcher struct {
	// The FIFO and Pixel Fetcher work together to ensure
	// that the FIFO always contains at least 8 pixels at any given time,
	// as 8 pixels are required for the Pixel Rendering operation to take place.

	// Sprites take priority unless they’re transparent (color 0)
	pixelFIFOBackground PixelFIFO
	pixelFIFOSprites    PixelFIFO

	// Pixel Fetcher
	fetcherX     uint8
	tileDataLow  uint8
	tileDataHigh uint8

	tileX       uint8
	tileY       uint8
	tileIndex   uint8
	tileAddress uint16
}

type NewPPU struct {
	mode   int
	ppuDot int

	// OAM Scan
	spriteIndex        int
	spritesSelected    [10]int
	numSpritesSelected int
	spriteYPosition    int

	fetcher PixelFetcher

	emu *Emulator
}

func (ppu *NewPPU) OAMScanReadYPosition() {
	spriteOffset := uint16(ppu.spriteIndex) * 4
	addr := 0xFE00 + spriteOffset
	ppu.spriteYPosition = int(ppu.emu.mem.io[addr&0x1ff]) - 16
}

func (ppu *NewPPU) OAMScanCheck() {
	// Sprite height
	height := 8
	lcdc := ppu.emu.mem.GetLCDC()
	if lcdc.ObjSize {
		height = 16
	}

	// Check if selected
	ly := int(ppu.emu.mem.GetLY())
	if ppu.numSpritesSelected < 10 && ly >= ppu.spriteYPosition && ly < ppu.spriteYPosition+height {
		ppu.spritesSelected[ppu.numSpritesSelected] = ppu.spriteIndex
		ppu.numSpritesSelected++
	}

	ppu.spriteIndex++
}

func (ppu *NewPPU) isInsideWindow(x uint8) bool {
	lcdc := ppu.emu.mem.GetLCDC()
	ly := ppu.emu.mem.GetLY()
	wy := ppu.emu.mem.GetWY()
	wx := ppu.emu.mem.GetWX()

	if lcdc.WindowEnable && ly >= wy && x >= (wx-7) {
		return true
	} else {
		return false
	}
}

func (ppu *NewPPU) FetcherGetTile() {
	lcdc := ppu.emu.mem.GetLCDC()
	isInsideWindow := ppu.isInsideWindow(ppu.fetcher.fetcherX)

	tilemap := 0x9800
	if (lcdc.BgTileMapArea && !isInsideWindow) || (lcdc.WindowTileMapArea && isInsideWindow) {
		tilemap = 0x9C00
	}
	fmt.Println(tilemap)

	xWindow := uint8(0)
	yWindow := uint8(0)

	if isInsideWindow {
		ppu.fetcher.tileX = xWindow
		ppu.fetcher.tileY = yWindow
	} else {
		ly := ppu.emu.mem.GetLY()
		scx := ppu.emu.mem.GetSCX()
		scy := ppu.emu.mem.GetSCY()
		ppu.fetcher.tileX = ((scx / 8) + ppu.fetcher.fetcherX) & 0x1F
		ppu.fetcher.tileY = ly + scy
	}

	ppu.fetcher.fetcherX++
}

func (ppu *NewPPU) getTileAddress(tileIndex uint8) uint16 {
	lcdc := ppu.emu.mem.GetLCDC()

	if lcdc.BgWindowTileDataArea { // 8000-8FFF
		return 0x8000 + uint16(tileIndex) // Unsigned addressing
	} else {
		return uint16(int32(0x8800) + int32(ppu.fetcher.tileIndex)) // Signed addressing
	}
}

func (ppu *NewPPU) FetcherGetTileDataLow() {
	// 2 Dots
	address := ppu.getTileAddress(ppu.fetcher.tileIndex)
	ppu.fetcher.tileAddress = address
	ppu.fetcher.tileDataLow = ppu.emu.mem.videoRam[address&0x1fff]
}

func (ppu *NewPPU) FetcherGetTileDataHigh() {
	// 2 Dots
	ppu.fetcher.tileAddress++
	ppu.fetcher.tileDataLow = ppu.emu.mem.videoRam[ppu.fetcher.tileAddress&0x1fff]
}

func (ppu *NewPPU) FetcherSleep() {
	// 2 Dots
	// Do nothing
}

func (ppu *NewPPU) FetcherPush() {
	// ?? dots

}

func (ppu *NewPPU) FetcherRun() {
	// ?? dots

}

func GetColorFromPalette(colorIndex uint8, paletteValue uint8) uint8 {
	// Value Color
	// 0	 White
	// 1	 Light gray
	// 2	 Dark gray
	// 3	 Black
	return (paletteValue >> (colorIndex * 2)) & 0x03
}

func (ppu *NewPPU) GetBackgroundColor(colorIndex uint8) uint8 {
	return GetColorFromPalette(colorIndex, ppu.emu.mem.GetBGP())
}

func (ppu *NewPPU) GetObjectColor(colorIndex uint8, palette uint8) uint8 {
	if palette == 0 {
		return GetColorFromPalette(colorIndex, ppu.emu.mem.GetOBP0())
	} else {
		return GetColorFromPalette(colorIndex, ppu.emu.mem.GetOBP1())
	}
}

func (ppu *NewPPU) PixelRendering() int {
	if !ppu.fetcher.pixelFIFOSprites.IsEmpty() && !ppu.fetcher.pixelFIFOBackground.IsEmpty() {
		lcdc := ppu.emu.mem.GetLCDC()

		pixelSprite, _ := ppu.fetcher.pixelFIFOSprites.Pop()
		pixelBackground, _ := ppu.fetcher.pixelFIFOBackground.Pop()

		// If LCDC.0 is disabled then the background is disabled on DMG and the background pixel won't have priority on CGB.
		// When the background pixel is disabled the pixel color value will be 0
		if lcdc.BgWindowEnablePriority {
			pixelBackground.color = 0x00
		}

		pixel := pixelBackground
		isObj := false
		if lcdc.ObjEnable {
			if pixelSprite.color != 0 && pixelSprite.backgroundPriority {
				pixel = pixelSprite
				isObj = true
			}
		}
		fmt.Println(pixel, isObj)
	}

	return 0
}

func (ppu *NewPPU) PPURun() bool {
	switch ppu.mode {
	case modeHblank: // H-Blank
	case modeVblank: // V-Blank
	case modeOAMScan: // OAM Scan
		if ppu.ppuDot%2 == 0 {
			ppu.OAMScanReadYPosition()
		} else {
			ppu.OAMScanCheck()
		}
	case modeDrawPixels: // Drawing pixels
		ppu.FetcherRun()
		ppu.PixelRendering()
	}

	ppu.ppuDot++

	// Move to next mode
	ly := ppu.emu.mem.GetLY()
	if ly < 143 {
		if ppu.ppuDot == 80 {
			ppu.mode = modeDrawPixels
			ppu.fetcher.pixelFIFOBackground.Clear()
			ppu.fetcher.pixelFIFOSprites.Clear()
		}
	}
	if ppu.ppuDot == 456 {
		if ly < 143 {
			ppu.mode = modeOAMScan
			ppu.numSpritesSelected = 0
			ppu.spriteIndex = 0
		} else if ly == 143 {
			ppu.mode = modeVblank
		}
	}

	if ppu.ppuDot == 456 {
		ppu.ppuDot = 0
		ly++
		ppu.emu.mem.SetLY(ly)
	}

	return false
}
