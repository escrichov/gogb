package emulator_old

//
//import (
//	"bytes"
//	"emulator-go/emulator/gb/utils"
//	"encoding/binary"
//	"errors"
//	"fmt"
//	"github.com/veandco/go-sdl2/sdl"
//	"log"
//	"os"
//	"path/filepath"
//	"syscall"
//	"unsafe"
//)
//
//func SetBit8(input uint8, pos uint8, value bool) uint8 {
//	result := input
//	if value {
//		result |= 1 << pos
//	} else {
//		result &= ^(1 << pos)
//	}
//
//	return result
//}
//
//func BoolToUint8(value bool) uint8 {
//	if value {
//		return 1
//	} else {
//		return 0
//	}
//}
//
//const WIDTH = 160
//const HEIGHT = 144
//
//type LCDControl struct {
//	LCDPPUEnable           bool // Bit 7, 0=Off, 1=On
//	WindowTileMapArea      bool // Bit 6, 0=9800-9BFF, 1=9C00-9FFF
//	WindowEnable           bool // Bit 5, 0=Off, 1=On
//	BgWindowTileDataArea   bool // Bit 4, 0=8800-97FF, 1=8000-8FFF
//	BgTileMapArea          bool // Bit 3, 0=9800-9BFF, 1=9C00-9FFF
//	ObjSize                bool // Bit 2, 0=8x8, 1=8x16
//	ObjEnable              bool // Bit 1, 0=Off, 1=On
//	BgWindowEnablePriority bool // Bit 0, 0=Off, 1=On
//}
//
//type RomHeader struct {
//	Title                string
//	TitleBytes           []byte
//	ColorGB              uint8
//	LicenseCodeNew       string
//	GBSGBIndicator       uint8
//	CartridgeType        uint8
//	RomSize              uint8
//	RamSize              uint8
//	DestinationCode      uint8
//	LicenseCodeOld       uint8
//	MaskROMVersionNumber uint8
//	ComplementCheck      uint8
//	CheckSum             uint16
//}
//
//type Emulator struct {
//	cycles, prevCycles uint16
//	totalCycles        uint64
//
//	IME  uint8
//	halt uint8
//
//	workRam  [0x4000]uint8
//	videoRam [0x2000]uint8
//
//	io                [0x200]uint8
//	extrambank        *[0x8000]uint8
//	ppuDot            int
//	rom0              []byte
//	bootRom           []byte
//	extrambankPointer uint32
//	rom1Pointer       uint32
//	keyboardState     []uint8
//	frameBuffer       [WIDTH * HEIGHT]int32
//	lcdcControl       LCDControl
//
//	palette []int32
//
//	cpu CPU
//
//	window   *sdl.Window
//	renderer *sdl.Renderer
//	texture  *sdl.Texture
//
//	numInstructions      uint64
//	vsyncEnabled         bool
//	stop                 bool
//	bootRomEnabled       bool
//	romHeader            RomHeader
//	memoryBankController int
//}
//
//func newEmulator(romFilename, saveFilename, bootRomFilename string) (*Emulator, error) {
//	emulator := Emulator{
//		ppuDot:       32,
//		rom1Pointer:  32768,
//		palette:      []int32{-1, -23197, -65536, -1 << 24, -1, -8092417, -12961132, -1 << 24},
//		vsyncEnabled: false,
//	}
//
//	if bootRomFilename == "" {
//		emulator.initializeBootRomValues()
//		emulator.bootRomEnabled = false
//	} else {
//		err := emulator.loadBootRom(bootRomFilename)
//		if err != nil {
//			return nil, err
//		}
//		emulator.bootRomEnabled = true
//	}
//
//	// Framebuffer set to black
//	for i, _ := range emulator.frameBuffer {
//		emulator.frameBuffer[i] = 0
//	}
//
//	err := emulator.loadRom(romFilename)
//	if err != nil {
//		return nil, err
//	}
//
//	err = emulator.initializeSDL(utils.ToCamel(emulator.romHeader.Title), 4)
//	if err != nil {
//		return nil, err
//	}
//
//	err = emulator.initializeSaveFile(saveFilename)
//	if err != nil {
//		return nil, err
//	}
//
//	return &emulator, nil
//}
//
//func (e *Emulator) initializeBootRomValues() {
//	e.cpu.PC = 256
//	e.SetLCDC(145)
//	e.SetLY(0)
//	e.SetDIV(44032)
//	e.cpu.SetSP(65534)
//	e.cpu.SetA(1)
//	e.cpu.SetF(176)
//	e.cpu.SetB(19)
//	e.cpu.SetC(0)
//	e.cpu.SetD(0)
//	e.cpu.SetE(216)
//	e.cpu.SetH(1)
//	e.cpu.SetL(77)
//}
//
//func (e *Emulator) initializeSDL(windowName string, windowScale float64) error {
//	var err error
//
//	e.window, err = sdl.CreateWindow(
//		windowName,
//		sdl.WINDOWPOS_CENTERED,
//		sdl.WINDOWPOS_CENTERED,
//		int32(WIDTH*windowScale),
//		int32(HEIGHT*windowScale),
//		sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
//	if err != nil {
//		log.Println("Error creating window:", err)
//		return err
//	}
//
//	e.renderer, err = sdl.CreateRenderer(e.window, -1, sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_ACCELERATED)
//	if err != nil {
//		log.Println("Error creating renderer:", err)
//		return err
//	}
//
//	// Creating a SDL texture that is used to display the color buffer
//	e.texture, err = e.renderer.CreateTexture(
//		uint32(sdl.PIXELFORMAT_RGBA32),
//		sdl.TEXTUREACCESS_STREAMING,
//		WIDTH,
//		HEIGHT,
//	)
//	if err != nil {
//		log.Println("Error creating texture:", err)
//		return err
//	}
//
//	// Point to Keyboard State
//	e.keyboardState = sdl.GetKeyboardState()
//
//	// Vsync
//	e.renderer.RenderSetVSync(e.vsyncEnabled)
//
//	return nil
//}
//
//func (e *Emulator) initializeSaveFile(fileName string) error {
//
//	t := int(unsafe.Sizeof(uint8(8))) * 32768
//	var mapFile *os.File
//
//	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
//		mapFile, err = os.Create(fileName)
//		if err != nil {
//			log.Println("Error opening file: ", err)
//			return err
//		}
//		_, err = mapFile.Seek(int64(t-1), 0)
//		if err != nil {
//			log.Println("Error opening file: ", err)
//			return err
//		}
//		_, err = mapFile.Write([]byte(" "))
//		if err != nil {
//			log.Println("Error writing file: ", err)
//			return err
//		}
//	} else {
//		mapFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
//		if err != nil {
//			log.Println("Error opening file: ", err)
//			return err
//		}
//	}
//
//	mmap, err := syscall.Mmap(int(mapFile.Fd()), 0, int(t), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	e.extrambank = (*[32768]uint8)(unsafe.Pointer(&mmap[0]))
//
//	return nil
//}
//
//func (e *Emulator) loadRom(fileName string) error {
//	var err error
//
//	bs, err := os.ReadFile(fileName)
//	if err != nil {
//		log.Println("Rom file not found:", err)
//		return err
//	}
//
//	fileExtension := filepath.Ext(fileName)
//	if fileExtension == ".zip" {
//		e.rom0, err = utils.UnzipBytes(bs)
//		if err != nil {
//			return err
//		}
//	} else {
//		e.rom0 = bs
//	}
//
//	err = e.parseRomHeader()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (e *Emulator) loadBootRom(fileName string) error {
//	var err error
//
//	e.bootRom, err = os.ReadFile(fileName)
//	if err != nil {
//		log.Println("BootRom file not found:", err)
//		return err
//	}
//
//	return nil
//}
//
//func parseRomTitle(bs []byte) string {
//	var end = 0
//	for n, b := range bs {
//		end = n
//		if b == 0 {
//			break
//		}
//	}
//	return string(bs[0:end])
//}
//
//func (e *Emulator) parseRomHeader() error {
//	if len(e.rom0) < 0x150 {
//		return fmt.Errorf("incorrect rom size %d", len(e.rom0))
//	}
//
//	e.romHeader.Title = parseRomTitle(e.rom0[0x134:0x144])
//	e.romHeader.TitleBytes = e.rom0[0x134:0x144]
//	e.romHeader.ColorGB = e.rom0[0x143]
//	e.romHeader.LicenseCodeNew = string(e.rom0[0x144:0x146])
//	e.romHeader.GBSGBIndicator = e.rom0[0x146]
//	e.romHeader.CartridgeType = e.rom0[0x147]
//	e.romHeader.RomSize = e.rom0[0x148]
//	e.romHeader.RamSize = e.rom0[0x149]
//	e.romHeader.DestinationCode = e.rom0[0x14A]
//	e.romHeader.LicenseCodeOld = e.rom0[0x14B]
//	e.romHeader.MaskROMVersionNumber = e.rom0[0x14C]
//	e.romHeader.ComplementCheck = e.rom0[0x14D]
//	e.romHeader.CheckSum = binary.BigEndian.Uint16(e.rom0[0x14E:0x150])
//
//	switch e.romHeader.CartridgeType {
//	case 0:
//		e.memoryBankController = 0
//	case 1, 2, 3:
//		e.memoryBankController = 1
//	case 5, 6:
//		e.memoryBankController = 2
//	case 8, 9:
//		e.memoryBankController = 0
//	case 0xB, 0xC, 0xD:
//		e.memoryBankController = 1
//	case 0xF, 0x10, 0x11, 0x12, 0x13:
//		e.memoryBankController = 3
//	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
//		e.memoryBankController = 5
//	}
//
//	return nil
//}
//
//type Register struct {
//	// The value of the register.
//	value uint16
//}
//
//// GetHi gets the higher byte of the register.
//func (reg *Register) GetHi() byte {
//	return byte(reg.value >> 8)
//}
//
//// GetLo gets the lower byte of the register.
//func (reg *Register) GetLo() byte {
//	return byte(reg.value & 0xFF)
//}
//
//// Get gets the 2 byte value of the register.
//func (reg *Register) Get() uint16 {
//	return reg.value
//}
//
//// SetHi sets the higher byte of the register.
//func (reg *Register) SetHi(val byte) {
//	reg.value = uint16(val)<<8 | (uint16(reg.value) & 0xFF)
//}
//
//// SetLo sets the lower byte of the register.
//func (reg *Register) SetLo(val byte) {
//	reg.value = uint16(val) | (uint16(reg.value) & 0xFF00)
//}
//
//// Set the value of the register.
//func (reg *Register) Set(val uint16) {
//	reg.value = val
//}
//
//// SetBit the value of the register.
//func (reg *Register) SetBit(pos int, value bool) {
//	utils.SetBit16(&reg.value, pos, value)
//}
//
//// GetBit the value of the register.
//func (reg *Register) GetBit(pos int) bool {
//	return utils.GetBit16(reg.value, pos)
//}
//
//// SetHighBit the value of the high part of the register.
//func (reg *Register) SetHighBit(highPos int, value bool) {
//	pos := highPos + 8
//	reg.SetBit(pos, value)
//}
//
//// SetLowBit the value of the low part of the register.
//func (reg *Register) SetLowBit(lowPos int, value bool) {
//	reg.SetBit(lowPos, value)
//}
//
//// GetHighBit the value of the high part of the register.
//func (reg *Register) GetHighBit(highPos int) bool {
//	pos := highPos + 8
//	return reg.GetBit(pos)
//}
//
//// GetLowBit the value of the low part of the register.
//func (reg *Register) GetLowBit(lowPos int) bool {
//	return reg.GetBit(lowPos)
//}
//
//type CPU struct {
//	AF Register // Accumulator & Flags Register (ZNHC---) -> N & H flags are not used -> (Z--C---)
//	BC Register
//	DE Register
//	HL Register
//
//	SP Register
//	PC uint16
//}
//
//// SetZeroFlag sets the value of the Zero flag.
//func (cpu *CPU) SetZeroFlag(value bool) {
//	cpu.AF.SetLowBit(7, value)
//}
//
//// SetSubtractFlag sets the value of Subtract flag.
//func (cpu *CPU) SetSubtractFlag(value bool) {
//	cpu.AF.SetLowBit(6, value)
//}
//
//// SetHalfCarryFlag sets the value of the Half Carry flag.
//func (cpu *CPU) SetHalfCarryFlag(value bool) {
//	cpu.AF.SetLowBit(5, value)
//}
//
//// SetCarryFlag sets the value of the Carry flag.
//func (cpu *CPU) SetCarryFlag(value bool) {
//	cpu.AF.SetLowBit(4, value)
//}
//
//// GetZeroFlag gets the value of the Zero flag.
//func (cpu *CPU) GetZeroFlag() bool {
//	return cpu.AF.GetLowBit(7)
//}
//
//// GetSubtractFlag gets the value of Subtract flag.
//func (cpu *CPU) GetSubtractFlag() bool {
//	return cpu.AF.GetLowBit(6)
//}
//
//// GetHalfCarryFlag gets the value of the Half Carry flag.
//func (cpu *CPU) GetHalfCarryFlag() bool {
//	return cpu.AF.GetLowBit(5)
//}
//
//// GetCarryFlag gets the value of the Carry flag.
//func (cpu *CPU) GetCarryFlag() bool {
//	return cpu.AF.GetLowBit(4)
//}
//
//func (cpu *CPU) GetA() uint8 {
//	return cpu.AF.GetHi()
//}
//
//func (cpu *CPU) SetA(value uint8) {
//	cpu.AF.SetHi(value)
//}
//
//func (cpu *CPU) GetF() uint8 {
//	return cpu.AF.GetLo()
//}
//
//func (cpu *CPU) SetF(value uint8) {
//	cpu.AF.SetLo(value)
//}
//
//func (cpu *CPU) GetB() uint8 {
//	return cpu.BC.GetHi()
//}
//
//func (cpu *CPU) SetB(value uint8) {
//	cpu.BC.SetHi(value)
//}
//
//func (cpu *CPU) GetC() uint8 {
//	return cpu.BC.GetLo()
//}
//
//func (cpu *CPU) SetC(value uint8) {
//	cpu.BC.SetLo(value)
//}
//
//func (cpu *CPU) GetD() uint8 {
//	return cpu.DE.GetHi()
//}
//
//func (cpu *CPU) SetD(value uint8) {
//	cpu.DE.SetHi(value)
//}
//
//func (cpu *CPU) GetE() uint8 {
//	return cpu.DE.GetLo()
//}
//
//func (cpu *CPU) SetE(value uint8) {
//	cpu.DE.SetLo(value)
//}
//
//func (cpu *CPU) GetH() uint8 {
//	return cpu.HL.GetHi()
//}
//
//func (cpu *CPU) SetH(value uint8) {
//	cpu.HL.SetHi(value)
//}
//
//func (cpu *CPU) GetL() uint8 {
//	return cpu.HL.GetLo()
//}
//
//func (cpu *CPU) SetL(value uint8) {
//	cpu.HL.SetLo(value)
//}
//
//func (cpu *CPU) GetAF() uint16 {
//	return cpu.AF.Get()
//}
//
//func (cpu *CPU) SetAF(value uint16) {
//	cpu.AF.Set(value)
//}
//
//func (cpu *CPU) GetBC() uint16 {
//	return cpu.BC.Get()
//}
//
//func (cpu *CPU) SetBC(value uint16) {
//	cpu.BC.Set(value)
//}
//
//func (cpu *CPU) GetDE() uint16 {
//	return cpu.DE.Get()
//}
//
//func (cpu *CPU) SetDE(value uint16) {
//	cpu.DE.Set(value)
//}
//
//func (cpu *CPU) GetHL() uint16 {
//	return cpu.HL.Get()
//}
//
//func (cpu *CPU) SetHL(value uint16) {
//	cpu.HL.Set(value)
//}
//
//func (cpu *CPU) GetSP() uint16 {
//	return cpu.SP.Get()
//}
//
//func (cpu *CPU) SetSP(value uint16) {
//	cpu.SP.Set(value)
//}
//
//func (e *Emulator) popPC() uint8 {
//	result := e.read8(e.cpu.PC)
//	e.cpu.PC++
//	return result
//}
//
//func (e *Emulator) popPC16() uint16 {
//	result := e.read16(e.cpu.PC)
//	e.cpu.PC += 2
//	return result
//}
//
//func (cpu *CPU) checkCondition(conditionNumber uint8) bool {
//	switch conditionNumber {
//	case 0:
//		return !cpu.GetZeroFlag()
//	case 1:
//		return cpu.GetZeroFlag()
//	case 2:
//		return !cpu.GetCarryFlag()
//	case 3:
//		return cpu.GetCarryFlag()
//	default:
//		return false
//	}
//}
//
//func (cpu *CPU) r16Group1Get(number uint8) uint16 {
//	switch number {
//	case 0:
//		return cpu.GetBC()
//	case 1:
//		return cpu.GetDE()
//	case 2:
//		return cpu.GetHL()
//	case 3:
//		return cpu.GetSP()
//	default:
//		return 0
//	}
//}
//
//func (cpu *CPU) r16Group1Set(number uint8, val uint16) {
//	switch number {
//	case 0:
//		cpu.SetBC(val)
//	case 1:
//		cpu.SetDE(val)
//	case 2:
//		cpu.SetHL(val)
//	case 3:
//		cpu.SetSP(val)
//	default:
//	}
//}
//
//func (cpu *CPU) r16Group2Get(number uint8) uint16 {
//	switch number {
//	case 0:
//		return cpu.GetBC()
//	case 1:
//		return cpu.GetDE()
//	case 2:
//		value := cpu.GetHL()
//		cpu.SetHL(value + 1)
//		return value
//	case 3:
//		value := cpu.GetHL()
//		cpu.SetHL(value - 1)
//		return value
//	default:
//		return 0
//	}
//}
//
//func (cpu *CPU) r16Group3Get(number uint8) uint16 {
//	switch number {
//	case 0:
//		return cpu.GetBC()
//	case 1:
//		return cpu.GetDE()
//	case 2:
//		return cpu.GetHL()
//	case 3:
//		return cpu.GetAF()
//	default:
//		return 0
//	}
//}
//
//func (cpu *CPU) r16Group3Set(number uint8, val uint16) {
//	switch number {
//	case 0:
//		cpu.SetBC(val)
//	case 1:
//		cpu.SetDE(val)
//	case 2:
//		cpu.SetHL(val)
//	case 3:
//		cpu.SetAF(val)
//	default:
//	}
//}
//
//func (e *Emulator) r8Get(number uint8) uint8 {
//	switch number {
//	case 0:
//		return e.cpu.GetB()
//	case 1:
//		return e.cpu.GetC()
//	case 2:
//		return e.cpu.GetD()
//	case 3:
//		return e.cpu.GetE()
//	case 4:
//		return e.cpu.GetH()
//	case 5:
//		return e.cpu.GetL()
//	case 6:
//		return e.read8(e.cpu.GetHL())
//	case 7:
//		return e.cpu.GetA()
//	default:
//		return 0
//	}
//}
//
//func (e *Emulator) r8Set(number uint8, val uint8) {
//	switch number {
//	case 0:
//		e.cpu.SetB(val)
//	case 1:
//		e.cpu.SetC(val)
//	case 2:
//		e.cpu.SetD(val)
//	case 3:
//		e.cpu.SetE(val)
//	case 4:
//		e.cpu.SetH(val)
//	case 5:
//		e.cpu.SetL(val)
//	case 6:
//		e.write8(e.cpu.GetHL(), val)
//	case 7:
//		e.cpu.SetA(val)
//	default:
//	}
//}
//
//func (e *Emulator) tick() {
//	e.cycles += 4
//	e.totalCycles += 4
//}
//
//func (e *Emulator) mem8(addr uint16, val uint8, write bool) uint8 {
//	e.tick()
//
//	switch addr >> 13 {
//	case 1: // 0x2000 - 0x3FFF
//		if write {
//			if e.memoryBankController == 3 {
//				// Pokemon Blue uses MBC3, which has the ability to swap 64 different 16KiB banks of ROM
//				var romBank uint32 = 1
//				if val != 0 {
//					romBank = uint32(val & 0x3F)
//				}
//				e.rom1Pointer = romBank << 14
//			} else if e.memoryBankController == 5 {
//				if addr <= 0x2FFF {
//					var romBank = uint32(val & 0x3F)
//					e.rom1Pointer = romBank << 14
//				} else {
//					// TODO: Implement set bit 9
//				}
//			}
//		}
//		return e.rom0[addr]
//	case 0: // 0x0000 - 0x1FFF
//		if e.bootRomEnabled && addr <= 0xFF {
//			return e.bootRom[addr]
//		}
//		return e.rom0[addr]
//	case 2: // 0x4000 - 0x5FFF
//		if e.memoryBankController == 3 || e.memoryBankController == 5 {
//			// 4 different of 8KiB banks of External Ram (for a total of 32KiB)
//			if write && val <= 3 {
//				e.extrambankPointer = uint32(val << 13)
//			}
//			return e.rom0[e.rom1Pointer+uint32(addr&0x3fff)]
//		} else {
//			return e.rom0[addr]
//		}
//	case 3: // 0x6000 - 0x7FFF
//		if e.memoryBankController == 3 || e.memoryBankController == 5 {
//			return e.rom0[e.rom1Pointer+uint32(addr&0x3fff)]
//		} else {
//			return e.rom0[addr]
//		}
//	case 4: // 0x8000 - 0x9FFF
//		addr &= 8191
//		if write {
//			e.videoRam[addr] = val
//		}
//		return e.videoRam[addr]
//
//	case 5: // 0xA000 - 0xBFFF
//		if e.memoryBankController == 3 || e.memoryBankController == 5 {
//			addr &= 0x1fff
//			if write {
//				e.extrambank[e.extrambankPointer+uint32(addr)] = val
//			}
//			return e.extrambank[e.extrambankPointer+uint32(addr)]
//		} else {
//			return 0
//		}
//	case 7: // 0xE000 - 0xFFFF
//		if addr >= 0xFE00 {
//			if write {
//				if addr == 0xFF46 {
//					for y := WIDTH - 1; y >= 0; y-- {
//						e.io[y] = e.read8(uint16(val)<<8 | uint16(y))
//					}
//				} else if addr == 0xFF40 {
//					e.SetLCDC(val)
//				} else if addr == 0xFF50 {
//					e.bootRomEnabled = false
//				}
//				ioAddr := addr & 0x1ff
//				e.io[ioAddr] = val
//			}
//
//			if addr == 0xff00 {
//				if (^e.io[256] & 16) != 0 {
//					return ^(16 + e.keyboardState[sdl.SCANCODE_DOWN]*8 +
//						e.keyboardState[sdl.SCANCODE_UP]*4 +
//						e.keyboardState[sdl.SCANCODE_LEFT]*2 +
//						e.keyboardState[sdl.SCANCODE_RIGHT])
//				}
//				if (^e.io[256] & 32) != 0 {
//					return ^(32 + e.keyboardState[sdl.SCANCODE_RETURN]*8 +
//						e.keyboardState[sdl.SCANCODE_TAB]*4 +
//						e.keyboardState[sdl.SCANCODE_Z]*2 +
//						e.keyboardState[sdl.SCANCODE_X])
//				}
//				return 0xFF
//			}
//			ioAddr := addr & 0x1ff
//			return e.io[ioAddr]
//		} else { // Echo internal RAM
//			addr &= 0x3fff
//			if write {
//				e.workRam[addr] = val
//			}
//			return e.workRam[addr]
//		}
//	case 6: // 0xC000 - 0xDFFF, Internal RAM
//		addr &= 0x3fff
//		if write {
//			e.workRam[addr] = val
//		}
//		return e.workRam[addr]
//	}
//
//	return 0
//}
//
//func (e *Emulator) getColor(tile, yOffset, xOffset int) uint8 {
//	videoRamIndex := tile*16 + yOffset*2
//	tileData := e.videoRam[videoRamIndex]
//	tileData1 := e.videoRam[videoRamIndex+1]
//	return ((tileData1>>xOffset)%2)*2 + (tileData>>xOffset)%2
//}
//
//func (e *Emulator) read16(addr uint16) uint16 {
//	tmp8 := e.mem8(addr, 0, false)
//	addr++
//	result := e.mem8(addr, 0, false)
//	addr++
//	return uint16(result)<<8 | uint16(tmp8)
//}
//
//func (e *Emulator) read8(addr uint16) uint8 {
//	return e.mem8(addr, 0, false)
//}
//
//func (e *Emulator) write16(addr uint16, val uint16) {
//	e.mem8(addr, uint8(val>>8), true)
//	addr++
//	e.mem8(addr, uint8(val), true)
//}
//
//func (e *Emulator) write8(addr uint16, val uint8) {
//	e.mem8(addr, val, true)
//}
//
//func (e *Emulator) push(val uint16) {
//	sp := e.cpu.GetSP()
//	sp--
//	e.write8(sp, uint8(val>>8))
//	sp--
//	e.write8(sp, uint8(val))
//	e.cpu.SetSP(sp)
//
//	e.tick()
//}
//
//func (e *Emulator) pop() uint16 {
//	sp := e.cpu.GetSP()
//	result := e.read16(sp)
//	e.cpu.SetSP(sp + 2)
//
//	return result
//}
//
//func (e *Emulator) SetIF(value uint8) {
//	e.io[271] = value
//}
//
//func (e *Emulator) GetIF() uint8 {
//	return e.io[271]
//}
//
//func (e *Emulator) SetLCDC(value uint8) {
//	e.io[320] = value
//	e.lcdcControl.LCDPPUEnable = utils.GetBit(value, 7)
//	e.lcdcControl.WindowTileMapArea = utils.GetBit(value, 6)
//	e.lcdcControl.WindowEnable = utils.GetBit(value, 5)
//	e.lcdcControl.BgWindowTileDataArea = utils.GetBit(value, 4)
//	e.lcdcControl.BgTileMapArea = utils.GetBit(value, 3)
//	e.lcdcControl.ObjSize = utils.GetBit(value, 2)
//	e.lcdcControl.ObjEnable = utils.GetBit(value, 1)
//	e.lcdcControl.BgWindowEnablePriority = utils.GetBit(value, 0)
//}
//
//func (e *Emulator) GetLCDC() *LCDControl {
//	return &e.lcdcControl
//}
//
//func (e *Emulator) SetLY(value uint8) {
//	e.io[324] = value
//}
//
//func (e *Emulator) GetLY() uint8 {
//	return e.io[324]
//}
//
//func (e *Emulator) SetDIV(value uint16) {
//	e.io[260] = uint8(value >> 8)
//	e.io[259] = uint8(value & 0xFF)
//}
//
//func (e *Emulator) GetDIV() uint16 {
//	return uint16(e.io[260])<<8 | uint16(e.io[259])
//}
//func (e *Emulator) CPURun() {
//	opcode := e.popPC()
//	e.numInstructions++
//
//	switch opcode {
//	case 0: // NOP
//	case 8: // LD (u16), SP
//		e.write16(e.popPC16(), e.cpu.GetSP())
//	case 16: // STOP (TODO: STOP Not implemented)
//		e.halt = 1
//		e.popPC()
//	case 24: // JR (unconditional)
//		i8 := int8(e.popPC())
//		addr := int32(e.cpu.PC) + int32(i8)
//		e.cpu.PC = uint16(addr)
//		e.tick()
//	case 32, 40, 48, 56: // JR (conditional)
//		i8 := int8(e.popPC())
//		if e.cpu.checkCondition((opcode >> 3) & 0x3) {
//			addr := int32(e.cpu.PC) + int32(i8)
//			e.cpu.PC = uint16(addr)
//			e.tick()
//		}
//	case 1, 17, 33, 49: // LD r16, u16
//		u16 := e.popPC16()
//		number := (opcode >> 4) & 0x3
//		e.cpu.r16Group1Set(number, u16)
//	case 9, 25, 41, 57: // ADD HL, r16
//		number := (opcode >> 4) & 0x3
//		r16 := e.cpu.r16Group1Get(number)
//		hl := e.cpu.GetHL()
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag((hl%4096 + r16%4096) > 4095)
//		e.cpu.SetCarryFlag((uint32(hl) + uint32(r16)) > 65535)
//		e.cpu.SetHL(hl + r16)
//		e.tick()
//	case 2, 18, 34, 50: // LD (r16), A
//		number := (opcode >> 4) & 0x3
//		e.write8(e.cpu.r16Group2Get(number), e.cpu.GetA())
//	case 10, 26, 42, 58: // LD A, (r16)
//		number := (opcode >> 4) & 0x3
//		e.cpu.SetA(e.read8(e.cpu.r16Group2Get(number)))
//	case 3, 19, 35, 51: // INC r16
//		number := (opcode >> 4) & 0x3
//		e.cpu.r16Group1Set(number, e.cpu.r16Group1Get(number)+1)
//		e.tick()
//	case 11, 27, 43, 59: // DEC r16
//		number := (opcode >> 4) & 0x3
//		r16 := e.cpu.r16Group1Get(number)
//		e.cpu.r16Group1Set(number, r16-1)
//		e.tick()
//	case 4, 12, 20, 28, 36, 44, 52, 60: // INC r8
//		number := (opcode >> 3) & 0x7
//		r8 := e.r8Get(number)
//		result := r8 + 1
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(result&0xF == 0)
//		e.r8Set(number, result)
//	case 5, 13, 21, 29, 37, 45, 53, 61: // DEC r8
//		number := (opcode >> 3) & 0x7
//		r8 := e.r8Get(number)
//		result := r8 - 1
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(true)
//		e.cpu.SetHalfCarryFlag((result+1)&15 == 0)
//		e.r8Set(number, r8-1)
//	case 6, 14, 22, 30, 38, 46, 54, 62: // LD r8, u8
//		u8 := e.popPC()
//		number := (opcode >> 3) & 0x7
//		e.r8Set(number, u8)
//	case 7: // RLCA
//		a := e.cpu.GetA()
//		bit7 := utils.GetBit(a, 7)
//		result := (a << 1) | BoolToUint8(bit7)
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(false)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(bit7)
//	case 15: // RRCA
//		a := e.cpu.GetA()
//		bit0 := utils.GetBit(a, 0)
//		result := (a >> 1) | (BoolToUint8(bit0) << 7)
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(false)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(bit0)
//	case 23: // RLA
//		a := e.cpu.GetA()
//		bit7 := utils.GetBit(a, 7)
//		carry := BoolToUint8(e.cpu.GetCarryFlag())
//		result := (a << 1) | carry
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(false)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(bit7)
//	case 31: // RRA
//		a := e.cpu.GetA()
//		bit0 := utils.GetBit(a, 0)
//		carry := BoolToUint8(e.cpu.GetCarryFlag())
//		result := (a >> 1) | (carry << 7)
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(false)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(bit0)
//	case 39: // DAA
//		tmp8 := uint8(0)
//		carry := false
//		if e.cpu.GetHalfCarryFlag() || !e.cpu.GetSubtractFlag() && e.cpu.GetA()%16 > 9 {
//			tmp8 = 6
//		}
//		if e.cpu.GetCarryFlag() || !e.cpu.GetSubtractFlag() && e.cpu.GetA() > 153 {
//			tmp8 |= 96
//			carry = true
//		}
//
//		a := e.cpu.GetA()
//		result := a
//		if e.cpu.GetSubtractFlag() {
//			result -= tmp8
//		} else {
//			result += tmp8
//		}
//
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(carry)
//	case 47: // CPL
//		a := e.cpu.GetA()
//		result := a ^ 0xFF
//		e.cpu.SetA(result)
//		e.cpu.SetSubtractFlag(true)
//		e.cpu.SetHalfCarryFlag(true)
//	case 55: // SCF
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(true)
//	case 63: // CCF
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(!e.cpu.GetCarryFlag())
//	case 118: // HALT
//		e.halt = 1
//	case 64, 65, 66, 67, 68, 69, 70, 71, // LD r8, r8
//		72, 73, 74, 75, 76, 77, 78, 79,
//		80, 81, 82, 83, 84, 85, 86, 87,
//		88, 89, 90, 91, 92, 93, 94, 95,
//		96, 97, 98, 99, 100, 101, 102, 103,
//		104, 105, 106, 107, 108, 109, 110, 111,
//		112, 113, 114, 115, 116, 117, 119, // 118 is missing because is defined HALT,
//		120, 121, 122, 123, 124, 125, 126, 127:
//		numberSource := opcode & 0x7
//		numberDestination := (opcode >> 3) & 0x7
//		e.r8Set(numberDestination, e.r8Get(numberSource))
//	case 128, 129, 130, 131, 132, 133, 134, 135: // ADD A, r8
//		number := opcode & 0x7
//		r8 := e.r8Get(number)
//		a := e.cpu.GetA()
//		result := a + r8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(a%16+r8%16 > 15)
//		e.cpu.SetCarryFlag(uint16(a)+uint16(r8) > 255)
//	case 136, 137, 138, 139, 140, 141, 142, 143: // ADC A, r8
//		number := opcode & 0x7
//		r8 := e.r8Get(number)
//		carry := BoolToUint8(e.cpu.GetCarryFlag())
//		a := e.cpu.GetA()
//		result := a + r8 + carry
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(a%16+r8%16+carry > 15)
//		e.cpu.SetCarryFlag(uint16(a)+uint16(r8)+uint16(carry) > 255)
//	case 144, 145, 146, 147, 148, 149, 150, 151: // SUB A, r8
//		number := opcode & 0x7
//		r8 := e.r8Get(number)
//		a := e.cpu.GetA()
//		result := a - r8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(true)
//		e.cpu.SetHalfCarryFlag(a%16-r8%16 > 15)
//		e.cpu.SetCarryFlag(uint16(a)-uint16(r8) > 255)
//	case 152, 153, 154, 155, 156, 157, 158, 159: // SBC A, r8
//		number := opcode & 0x7
//		r8 := e.r8Get(number)
//		carry := BoolToUint8(e.cpu.GetCarryFlag())
//		a := e.cpu.GetA()
//		result := a - r8 - carry
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(true)
//		e.cpu.SetHalfCarryFlag(a%16-r8%16-carry > 15)
//		e.cpu.SetCarryFlag(uint16(a)-uint16(r8)-uint16(carry) > 255)
//	case 160, 161, 162, 163, 164, 165, 166, 167: // AND A, r8
//		number := opcode & 0x7
//		r8 := e.r8Get(number)
//		a := e.cpu.GetA()
//		result := a & r8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(true)
//		e.cpu.SetCarryFlag(false)
//	case 168, 169, 170, 171, 172, 173, 174, 175: // XOR A, r8
//		number := opcode & 0x7
//		r8 := e.r8Get(number)
//		a := e.cpu.GetA()
//		result := a ^ r8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(false)
//	case 176, 177, 178, 179, 180, 181, 182, 183: // OR A, r8
//		number := opcode & 0x7
//		r8 := e.r8Get(number)
//		a := e.cpu.GetA()
//		result := a | r8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(false)
//	case 184, 185, 186, 187, 188, 189, 190, 191: // CP A, r8
//		number := opcode & 0x7
//		r8 := e.r8Get(number)
//		a := e.cpu.GetA()
//		result := a - r8
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(true)
//		e.cpu.SetHalfCarryFlag(a%16-r8%16 > 15)
//		e.cpu.SetCarryFlag(uint16(a)-uint16(r8) > 255)
//	case 192, 200, 208, 216: // RET condition
//		e.tick()
//		if e.cpu.checkCondition((opcode >> 3) & 0x3) {
//			addr := e.pop()
//			e.cpu.PC = addr
//		}
//	case 224: // LD (FF00 + u8), A
//		u8 := e.popPC()
//		addr := 0xFF00 + uint16(u8)
//		e.write8(addr, e.cpu.GetA())
//	case 232: // ADD SP, i8
//		i8 := int8(e.popPC())
//		sp := e.cpu.GetSP()
//		result := int32(sp) + int32(i8)
//		e.cpu.SetSP(uint16(result))
//		e.cpu.SetZeroFlag(false)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(int16(i8%16)+int16(uint8(i8)%16) > 15)
//		e.cpu.SetCarryFlag(int16(sp)+int16(i8) > 255)
//		e.tick()
//		e.tick()
//	case 240: // LD A, (FF00 + u8)
//		u8 := e.popPC()
//		addr := 0xFF00 + uint16(u8)
//		e.cpu.SetA(e.read8(addr))
//	case 248: // LD HL, SP + i8
//		i8 := int8(e.popPC())
//		sp := e.cpu.GetSP()
//		result := int32(sp) + int32(i8)
//		e.cpu.SetHL(uint16(result))
//		e.cpu.SetZeroFlag(false)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(int16(sp%16)+int16(i8%16) > 15)
//		e.cpu.SetCarryFlag(uint16(uint8(sp))+uint16(i8) > 255)
//		e.tick()
//	case 193, 209, 225, 241: // POP r16
//		number := (opcode >> 4) & 0x3
//		value := e.pop()
//		e.cpu.r16Group3Set(number, value)
//	case 201: // RET
//		pc := e.pop()
//		e.cpu.PC = pc
//		e.tick()
//	case 217: // RETI
//		e.IME = 1
//		pc := e.pop()
//		e.cpu.PC = pc
//		e.tick()
//	case 233: // JP HL
//		e.cpu.PC = e.cpu.GetHL()
//	case 249: // LD SP, HL
//		e.cpu.SetSP(e.cpu.GetHL())
//		e.tick()
//	case 194, 202, 210, 218: // JP condition
//		addr := e.popPC16()
//		if e.cpu.checkCondition((opcode >> 3) & 0x3) {
//			e.cpu.PC = addr
//			e.tick()
//		}
//	case 226: // LD (FF00+C), A
//		addr := 0xFF00 + uint16(e.cpu.GetC())
//		e.write8(addr, e.cpu.GetA())
//	case 234: // LD (u16), A
//		u16 := e.popPC16()
//		e.write8(u16, e.cpu.GetA())
//	case 242: // LD A, (0xFF00+C)
//		addr := 0xFF00 + uint16(e.cpu.GetC())
//		e.cpu.SetA(e.read8(addr))
//	case 250: // LD A, (u16)
//		u16 := e.popPC16()
//		e.cpu.SetA(e.read8(u16))
//	case 195: // JP u16
//		e.cpu.PC = e.popPC16()
//		e.tick()
//	case 243: // DI
//		e.IME = 0
//	case 251: // EI
//		e.IME = 1
//	case 196, 204, 212, 220: // CALL condition
//		u16 := e.popPC16()
//		if e.cpu.checkCondition((opcode >> 3) & 0x3) {
//			e.push(e.cpu.PC)
//			e.cpu.PC = u16
//		}
//	case 197, 213, 229, 245: // PUSH r16
//		number := (opcode >> 4) & 0x3
//		value := e.cpu.r16Group3Get(number)
//		e.push(value)
//	case 205: // CALL u16
//		u16 := e.popPC16()
//		e.push(e.cpu.PC)
//		e.cpu.PC = u16
//	case 198: // ADD A, u8
//		u8 := e.popPC()
//		a := e.cpu.GetA()
//		result := a + u8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(a%16+u8%16 > 15)
//		e.cpu.SetCarryFlag(uint16(a)+uint16(u8) > 255)
//	case 206: // ADC A, u8
//		u8 := e.popPC()
//		carry := BoolToUint8(e.cpu.GetCarryFlag())
//		a := e.cpu.GetA()
//		result := a + u8 + carry
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(a%16+u8%16+carry > 15)
//		e.cpu.SetCarryFlag(uint16(a)+uint16(u8)+uint16(carry) > 255)
//	case 214: // SUB A, u8
//		u8 := e.popPC()
//		a := e.cpu.GetA()
//		result := a - u8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(true)
//		e.cpu.SetHalfCarryFlag(a%16-u8%16 > 15)
//		e.cpu.SetCarryFlag(uint16(a)-uint16(u8) > 255)
//	case 222: // SBC A, u8
//		u8 := e.popPC()
//		carry := BoolToUint8(e.cpu.GetCarryFlag())
//		a := e.cpu.GetA()
//		result := a - u8 - carry
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(true)
//		e.cpu.SetHalfCarryFlag(a%16-u8%16-carry > 15)
//		e.cpu.SetCarryFlag(uint16(a)-uint16(u8)-uint16(carry) > 255)
//	case 230: // AND A, u8
//		u8 := e.popPC()
//		a := e.cpu.GetA()
//		result := a & u8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(true)
//		e.cpu.SetCarryFlag(false)
//	case 238: // XOR A, u8
//		u8 := e.popPC()
//		a := e.cpu.GetA()
//		result := a ^ u8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(false)
//	case 246: // OR A, u8
//		u8 := e.popPC()
//		a := e.cpu.GetA()
//		result := a | u8
//		e.cpu.SetA(result)
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(false)
//		e.cpu.SetHalfCarryFlag(false)
//		e.cpu.SetCarryFlag(false)
//	case 254: // CP A, u8
//		u8 := e.popPC()
//		a := e.cpu.GetA()
//		result := a - u8
//		e.cpu.SetZeroFlag(result == 0)
//		e.cpu.SetSubtractFlag(true)
//		e.cpu.SetHalfCarryFlag(a%16-u8%16 > 15)
//		e.cpu.SetCarryFlag(uint16(a)-uint16(u8) > 255)
//	case 199, 207, 215, 223, 231, 239, 247, 255: // RST (Call to 00EXP000)
//		addr := uint16(opcode & 0x38)
//		e.push(e.cpu.PC)
//		e.cpu.PC = addr
//	case 0xCB:
//		opcode := e.popPC()
//		switch opcode {
//		case 0, 1, 2, 3, 4, 5, 6, 7: // RLC
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			bit7 := utils.GetBit(r8, 7)
//			result := (r8 << 1) | BoolToUint8(bit7)
//			e.r8Set(number, result)
//			e.cpu.SetZeroFlag(result == 0)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(false)
//			e.cpu.SetCarryFlag(bit7)
//		case 8, 9, 10, 11, 12, 13, 14, 15: // RRC
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			bit0 := utils.GetBit(r8, 0)
//			result := (r8 >> 1) | (BoolToUint8(bit0) << 7)
//			e.r8Set(number, result)
//			e.cpu.SetZeroFlag(result == 0)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(false)
//			e.cpu.SetCarryFlag(bit0)
//		case 16, 17, 18, 19, 20, 21, 22, 23: // RL
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			bit7 := utils.GetBit(r8, 7)
//			carry := BoolToUint8(e.cpu.GetCarryFlag())
//			result := (r8 << 1) | carry
//			e.r8Set(number, result)
//			e.cpu.SetZeroFlag(result == 0)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(false)
//			e.cpu.SetCarryFlag(bit7)
//		case 24, 25, 26, 27, 28, 29, 30, 31: // RR
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			bit0 := utils.GetBit(r8, 0)
//			carry := BoolToUint8(e.cpu.GetCarryFlag())
//			result := (r8 >> 1) | (carry << 7)
//			e.r8Set(number, result)
//			e.cpu.SetZeroFlag(result == 0)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(false)
//			e.cpu.SetCarryFlag(bit0)
//		case 32, 33, 34, 35, 36, 37, 38, 39: // SLA
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			bit7 := utils.GetBit(r8, 7)
//			result := r8 << 1
//			e.r8Set(number, result)
//			e.cpu.SetZeroFlag(result == 0)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(false)
//			e.cpu.SetCarryFlag(bit7)
//		case 40, 41, 42, 43, 44, 45, 46, 47: // SRA
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			bit0 := utils.GetBit(r8, 0)
//			bit7 := utils.GetBit(r8, 7)
//			result := r8>>1 | (BoolToUint8(bit7) << 7)
//			e.r8Set(number, result)
//			e.cpu.SetZeroFlag(result == 0)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(false)
//			e.cpu.SetCarryFlag(bit0)
//		case 48, 49, 50, 51, 52, 53, 54, 55: // SWAP
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			lower := r8 & 0xF
//			upper := r8 >> 4
//			result := (lower << 4) | upper
//			e.r8Set(number, result)
//			e.cpu.SetZeroFlag(result == 0)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(false)
//			e.cpu.SetCarryFlag(false)
//		case 56, 57, 58, 59, 60, 61, 62, 63: // SRL
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			bit0 := utils.GetBit(r8, 0)
//			result := r8 >> 1
//			e.r8Set(number, result)
//			e.cpu.SetZeroFlag(result == 0)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(false)
//			e.cpu.SetCarryFlag(bit0)
//		case 64, 65, 66, 67, 68, 69, 70, 71, // BIT bit, r8
//			72, 73, 74, 75, 76, 77, 78, 79,
//			80, 81, 82, 83, 84, 85, 86, 87,
//			88, 89, 90, 91, 92, 93, 94, 95,
//			96, 97, 98, 99, 100, 101, 102, 103,
//			104, 105, 106, 107, 108, 109, 110, 111,
//			112, 113, 114, 115, 116, 117, 118, 119,
//			120, 121, 122, 123, 124, 125, 126, 127:
//			bit := (opcode >> 3) & 0x7
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			e.cpu.SetSubtractFlag(false)
//			e.cpu.SetHalfCarryFlag(true)
//			e.cpu.SetZeroFlag(!utils.GetBit(r8, int(bit)))
//		case 128, 129, 130, 131, 132, 133, 134, 135, // RES bit, r8
//			136, 137, 138, 139, 140, 141, 142, 143,
//			144, 145, 146, 147, 148, 149, 150, 151,
//			152, 153, 154, 155, 156, 157, 158, 159,
//			160, 161, 162, 163, 164, 165, 166, 167,
//			168, 169, 170, 171, 172, 173, 174, 175,
//			176, 177, 178, 179, 180, 181, 182, 183,
//			184, 185, 186, 187, 188, 189, 190, 191:
//			bit := (opcode >> 3) & 0x7
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			e.r8Set(number, SetBit8(r8, bit, false))
//		case 192, 193, 194, 195, 196, 197, 198, 199, // SET bit, r8
//			200, 201, 202, 203, 204, 205, 206, 207,
//			208, 209, 210, 211, 212, 213, 214, 215,
//			216, 217, 218, 219, 220, 221, 222, 223,
//			224, 225, 226, 227, 228, 229, 230, 231,
//			232, 233, 234, 235, 236, 237, 238, 239,
//			240, 241, 242, 243, 244, 245, 246, 247,
//			248, 249, 250, 251, 252, 253, 254, 255:
//			bit := (opcode >> 3) & 0x7
//			number := opcode & 0x7
//			r8 := e.r8Get(number)
//			e.r8Set(number, SetBit8(r8, bit, true))
//		default:
//			log.Println("CB Opcode: ", opcode, " not found")
//			return
//		}
//	default:
//		log.Println("Opcode: ", opcode, " not found")
//		return
//	}
//}
//
//func (e *Emulator) PPURun() {
//	// PPU
//	div := e.GetDIV()
//	e.SetDIV(div + e.cycles - e.prevCycles)
//	for ; e.prevCycles != e.cycles; e.prevCycles++ {
//		lcdc := e.GetLCDC()
//		if lcdc.LCDPPUEnable {
//			e.ppuDot++
//
//			// Render Scanline (Every 256 PPU Dots)
//			if e.ppuDot == 456 {
//				ly := e.GetLY()
//
//				// Only render visible lines (up to line 144)
//				if ly < HEIGHT {
//					for tmp := WIDTH - 1; tmp >= 0; tmp-- {
//
//						// IsWindow
//						isWindow := false
//						if lcdc.WindowEnable && ly >= e.io[330] && uint8(tmp) >= (e.io[331]-7) {
//							isWindow = true
//						}
//
//						// xOffset
//						var xOffset uint8
//						if isWindow {
//							xOffset = uint8(tmp) - e.io[331] + 7
//						} else {
//							xOffset = uint8(tmp) + e.io[323]
//						}
//
//						// yOffset
//						var yOffset uint8
//						if isWindow {
//							yOffset = ly - e.io[330]
//						} else {
//							yOffset = ly + e.io[322]
//						}
//
//						// PaletteIndex
//						var paletteIndex uint16 = 0
//
//						// Tile
//						tileMapArea := lcdc.BgTileMapArea
//						if isWindow {
//							tileMapArea = lcdc.WindowTileMapArea
//						}
//
//						videoRamIndex := uint16(6)
//						if tileMapArea {
//							videoRamIndex = 7
//						}
//						videoRamIndex = videoRamIndex<<10 | uint16(yOffset)/8*32 + uint16(xOffset)/8
//						var tile = e.videoRam[videoRamIndex]
//
//						// Color
//						var tileValue int
//						if lcdc.BgWindowTileDataArea {
//							tileValue = int(tile)
//						} else {
//							tileValue = 256 + int(int8(tile))
//						}
//						color := e.getColor(tileValue, int(yOffset&7), int(7-xOffset&7))
//
//						// Sprites
//						if lcdc.ObjEnable {
//							for spriteIndex := uint8(0); spriteIndex < WIDTH; spriteIndex += 4 {
//								spriteX := uint8(tmp) - e.io[spriteIndex+1] + 8
//								spriteY := ly - e.io[spriteIndex] + 16
//
//								spriteYOffset := uint8(0)
//								if (e.io[spriteIndex+3] & 64) != 0 {
//									spriteYOffset = 7
//								}
//								spriteYOffset = spriteY ^ spriteYOffset
//
//								spriteXOffset := uint8(7)
//								if (e.io[spriteIndex+3] & 32) != 0 {
//									spriteXOffset = 0
//								}
//								spriteXOffset = spriteX ^ spriteXOffset
//
//								spriteColor := e.getColor(int(e.io[spriteIndex+2]), int(spriteYOffset), int(spriteXOffset))
//
//								if spriteX < 8 && spriteY < 8 && !((e.io[spriteIndex+3]&128) != 0 && color != 0) && spriteColor != 0 {
//									color = spriteColor
//									if e.io[spriteIndex+3]&16 == 0 {
//										paletteIndex = uint16(1)
//									} else {
//										paletteIndex = uint16(2)
//									}
//									break
//								}
//							}
//						}
//
//						paletteIndexValue := uint16((e.io[327+paletteIndex]>>(2*color))%4) + paletteIndex*4&7
//						frameBufferIndex := uint16(ly)*WIDTH + uint16(tmp)
//						e.frameBuffer[frameBufferIndex] = e.palette[paletteIndexValue]
//					}
//				}
//
//				if ly == (HEIGHT - 1) {
//					e.SetIF(e.GetIF() | 1)
//
//					e.renderFrame()
//					e.manageKeyboardEvents()
//				}
//
//				// Increment Line
//				e.SetLY((ly + 1) % 154)
//				e.ppuDot = 0
//			}
//		} else {
//			e.SetLY(0)
//			e.ppuDot = 0
//		}
//	}
//}
//
//func (e *Emulator) renderFrame() {
//	e.renderer.Clear()
//	buf := unsafe.Pointer(&e.frameBuffer[0])
//	framebufferBytes := unsafe.Slice((*byte)(buf), WIDTH*HEIGHT)
//	e.texture.Update(nil, framebufferBytes, WIDTH*4)
//	e.renderer.Copy(e.texture, nil, nil)
//	e.renderer.Present()
//}
//
//func (e *Emulator) manageKeyboardEvents() {
//	// Manage Keyboard Events
//	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
//		switch event.(type) {
//		case *sdl.KeyboardEvent:
//			keyboardEvent, ok := event.(*sdl.KeyboardEvent)
//			if !ok {
//				continue
//			}
//			if keyboardEvent.Type == sdl.KEYDOWN {
//				switch keyboardEvent.Keysym.Sym {
//				case sdl.K_ESCAPE:
//					e.stop = true
//				case sdl.K_v:
//					e.vsyncEnabled = !e.vsyncEnabled
//					e.renderer.RenderSetVSync(e.vsyncEnabled)
//				case sdl.K_s:
//					e.BessStore("save.bess")
//				}
//			}
//		case *sdl.QuitEvent:
//			e.stop = true
//		}
//	}
//}
//
//func (e *Emulator) Destroy() {
//	e.texture.Destroy()
//	e.renderer.Destroy()
//	e.window.Destroy()
//	sdl.Quit()
//}
//
//func (e *Emulator) Run() {
//	e.stop = false
//	for {
//		e.prevCycles = e.cycles
//		if (e.IME & e.GetIF() & e.io[511]) != 0 {
//			e.SetIF(0)
//			e.halt = 0
//			e.IME = 0
//			e.tick()
//			e.tick()
//			e.push(e.cpu.PC)
//			e.cpu.PC = 64
//		} else if e.halt != 0 {
//			e.tick()
//		} else {
//			e.CPURun()
//		}
//
//		e.PPURun()
//
//		if e.stop {
//			break
//		}
//	}
//}
//
//func (e *Emulator) BessStore(filename string) error {
//	data := new(bytes.Buffer)
//
//	// Name Block
//	data.WriteString("NAME")
//	binary.Write(data, binary.LittleEndian, int32(15))
//	data.WriteString("EMULATOR-GO 0.1")
//
//	// Info Block
//	data.WriteString("INFO")
//	binary.Write(data, binary.LittleEndian, int32(0x12))
//	data.Write(e.rom0[0x134:0x144]) // ROM (Title)
//	data.Write(e.rom0[0x14E:0x150]) // ROM (Global checksum)
//
//	// Core Block
//	data.WriteString("CORE")
//	binary.Write(data, binary.LittleEndian, int32(0xD0))
//	binary.Write(data, binary.LittleEndian, int16(1)) // Major BESS version as a 16-bit integer
//	binary.Write(data, binary.LittleEndian, int16(1)) // Major Minor version as a 16-bit integer
//	data.WriteString("GDA ")                          // A four-character ASCII model identifier
//	binary.Write(data, binary.LittleEndian, e.cpu.PC)
//	binary.Write(data, binary.LittleEndian, e.cpu.AF.value)
//	binary.Write(data, binary.LittleEndian, e.cpu.BC.value)
//	binary.Write(data, binary.LittleEndian, e.cpu.DE.value)
//	binary.Write(data, binary.LittleEndian, e.cpu.HL.value)
//	binary.Write(data, binary.LittleEndian, e.cpu.SP.value)
//	binary.Write(data, binary.LittleEndian, e.IME)
//	binary.Write(data, binary.LittleEndian, e.GetIF()) // The value of the IE register
//	binary.Write(data, binary.LittleEndian, e.halt)    // Execution state (0 = running; 1 = halted; 2 = stopped)
//	binary.Write(data, binary.LittleEndian, uint8(0))  // Reserved, must be 0
//	data.Write(e.io[0x100:0x180])                      // Memory-mapped Registers
//
//	binary.Write(data, binary.LittleEndian, int32(0x4000)) // The size of RAM (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of RAM from file start (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0x2000)) // The size of VRAM (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of VRAM from file start (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0))      // The size of MBC RAM (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of MBC RAM from file start (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0xA0))   // The size of OAM (=0xA0, 32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of OAM from file start (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0x7F))   // The size of HRAM (=0x7F, 32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of HRAM from file start (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0x40))   // The size of background palettes (=0x40 or 0, 32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of background palettes from file start (32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0x40))   // The size of object palettes (=0x40 or 0, 32-bit integer)
//	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of object palettes from file start (32-bit integer)
//
//	// XOAM block - Not implemented
//	// MBC block - Not implemented
//	// RTC block - Not implemented
//	// HUC3 block - Not implemented
//	// TPP1 block - Not implemented
//	// MBC7 block - Not implemented
//	// SGB block - Not implemented
//
//	// End Block
//	data.WriteString("END ")
//	binary.Write(data, binary.LittleEndian, int32(0))
//
//	// Footer
//	binary.Write(data, binary.LittleEndian, int32(0))
//	data.WriteString("BESS")
//
//	err := os.WriteFile(filename, data.Bytes(), 0644)
//	if err != nil {
//		log.Println("Error creating save state:", err)
//		return err
//	}
//
//	return nil
//}
//
//func (e *Emulator) BessLoad(filename string) error {
//	bs, err := os.ReadFile(filename)
//	if err != nil {
//		log.Println("Error creating save state:", err)
//		return err
//	}
//
//	// Footer
//	bess := string(bs[len(bs)-4:])
//	startOfFile := binary.LittleEndian.Uint32(bs[len(bs)-8 : len(bs)-4])
//
//	if bess != "BESS" {
//		return fmt.Errorf("bess not in footer")
//	}
//
//	buffer := bytes.NewReader(bs[startOfFile : len(bs)-8])
//
//	var blockNameTmp []byte = make([]byte, 4)
//
//	// Name Block
//	_, err = buffer.Read(blockNameTmp)
//	if err != nil {
//		return err
//	}
//	blockName := string(blockNameTmp)
//
//	var blockSize int32
//	err = binary.Read(buffer, binary.LittleEndian, &blockSize)
//	if err != nil {
//		return err
//	}
//
//	var blockContent []byte = make([]byte, 15)
//	_, err = buffer.Read(blockContent)
//	if err != nil {
//		return err
//	}
//	fmt.Println(blockName, blockSize, string(blockContent))
//
//	// Info  Block
//	_, err = buffer.Read(blockNameTmp)
//	if err != nil {
//		return err
//	}
//	blockName = string(blockNameTmp)
//
//	err = binary.Read(buffer, binary.LittleEndian, &blockSize)
//	if err != nil {
//		return err
//	}
//
//	blockContent = make([]byte, blockSize)
//	_, err = buffer.Read(blockContent)
//	if err != nil {
//		return err
//	}
//	romTitle := blockContent[:16]
//	romCheckSum := binary.BigEndian.Uint16(blockContent[16:])
//	fmt.Println(blockName, blockSize, string(romTitle), romCheckSum)
//
//	// Core Block
//	_, err = buffer.Read(blockNameTmp)
//	if err != nil {
//		return err
//	}
//	blockName = string(blockNameTmp)
//
//	err = binary.Read(buffer, binary.LittleEndian, &blockSize)
//	if err != nil {
//		return err
//	}
//	fmt.Println(blockName, blockSize)
//
//	var majorVersion int16
//	err = binary.Read(buffer, binary.LittleEndian, &majorVersion)
//	if err != nil {
//		return err
//	}
//
//	var minorVersion int16
//	err = binary.Read(buffer, binary.LittleEndian, &minorVersion)
//	if err != nil {
//		return err
//	}
//	fmt.Printf("BESS Version: %d.%d\n", majorVersion, minorVersion)
//
//	var modelIdentifier []byte = make([]byte, 4)
//	_, err = buffer.Read(modelIdentifier)
//	if err != nil {
//		return err
//	}
//	fmt.Println("Model Identifier:", string(modelIdentifier))
//
//	var pc, af, bc, de, hl, sp uint16
//	var ime, ie, halt uint8
//
//	err = binary.Read(buffer, binary.LittleEndian, &pc)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &af)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &bc)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &de)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &hl)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &sp)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &ime)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &ie)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &halt)
//	if err != nil {
//		return err
//	}
//	reserved, err := buffer.ReadByte()
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("PC: 0x%x\n", pc)
//	fmt.Printf("AF: 0x%x\n", af)
//	fmt.Printf("BC: 0x%x\n", bc)
//	fmt.Printf("DE: 0x%x\n", de)
//	fmt.Printf("HL: 0x%x\n", hl)
//	fmt.Printf("SP: 0x%x\n", sp)
//	fmt.Printf("IME: %d\n", ime)
//	fmt.Printf("IE: %d\n", ie)
//	fmt.Printf("HALT: %d\n", halt)
//
//	var io [128]byte
//	_, err = buffer.Read(io[:])
//	if err != nil {
//		return err
//	}
//
//	var ramSize, vramSize, mbcRamSize, OAMSize, HRAMSize, BackgroundPalettesSize, ObjectPalettesSize int32
//	var ramOffset, vramOffset, mbcRamOffset, OAMOffset, HRAMOffset, BackgroundPalettesOffset, ObjectPalettesOffset int32
//	err = binary.Read(buffer, binary.LittleEndian, &ramSize)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &ramOffset)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &vramSize)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &vramOffset)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &mbcRamSize)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &mbcRamOffset)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &OAMSize)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &OAMOffset)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &HRAMSize)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &HRAMOffset)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &BackgroundPalettesSize)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &BackgroundPalettesOffset)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &ObjectPalettesSize)
//	if err != nil {
//		return err
//	}
//	err = binary.Read(buffer, binary.LittleEndian, &ObjectPalettesOffset)
//	if err != nil {
//		return err
//	}
//
//	for {
//		_, err = buffer.Read(blockNameTmp)
//		if err != nil {
//			return err
//		}
//		blockName = string(blockNameTmp)
//
//		err = binary.Read(buffer, binary.LittleEndian, &blockSize)
//		if err != nil {
//			return err
//		}
//
//		fmt.Println([]byte(blockName), blockName, blockSize)
//
//		if blockName == "END " {
//			if blockSize != 0 {
//				return fmt.Errorf("bess end block size (%d) != 0", blockSize)
//			}
//			fmt.Println(blockName, blockSize)
//			break
//		} else {
//			blockContent = make([]byte, blockSize)
//			_, err = buffer.Read(blockContent)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	// Extra checks
//	if res := bytes.Compare(romTitle, e.romHeader.TitleBytes); res != 0 {
//		return fmt.Errorf("incorrect rom title: %v != %v", []byte(romTitle), []byte(e.romHeader.Title))
//	}
//
//	if romCheckSum != e.romHeader.CheckSum {
//		return fmt.Errorf("incorrect rom checksum: %d != %d", romCheckSum, e.romHeader.CheckSum)
//	}
//
//	if majorVersion != 1 {
//		return fmt.Errorf("major version not supported: %d", majorVersion)
//	}
//
//	if modelIdentifier[0] != 'G' {
//		return fmt.Errorf("model identifier not supported: %c", modelIdentifier)
//	}
//
//	if reserved != 0 {
//		return fmt.Errorf("reserved byte with offset 0x17 0x%x != 0", reserved)
//	}
//
//	if ime > 1 {
//		return fmt.Errorf("incorrect ime: %d", ime)
//	}
//
//	if halt > 2 {
//		return fmt.Errorf("incorrect execution state: %d", halt)
//	}
//
//	if ramSize != 0x2000 {
//		return fmt.Errorf("incorrect ram size: %d", ramSize)
//	}
//
//	if vramSize != 0x2000 {
//		return fmt.Errorf("incorrect vram size: %d", vramSize)
//	}
//
//	if OAMSize != 0xA0 {
//		return fmt.Errorf("incorrect oam size: %d", OAMSize)
//	}
//
//	if HRAMSize != 0x7F {
//		return fmt.Errorf("incorrect hram size: %d", HRAMSize)
//	}
//
//	if BackgroundPalettesSize != 0x40 {
//		return fmt.Errorf("incorrect background palettes size: %d", BackgroundPalettesSize)
//	}
//
//	if ObjectPalettesSize != 0x40 {
//		return fmt.Errorf("incorrect object palettes size: %d", ObjectPalettesSize)
//	}
//
//	// Set values
//	e.cpu.PC = pc
//	e.cpu.SetAF(af)
//	e.cpu.SetBC(bc)
//	e.cpu.SetDE(de)
//	e.cpu.SetHL(hl)
//	e.cpu.SetSP(sp)
//	e.SetIF(ie)
//	e.IME = ime
//	e.halt = halt
//
//	return nil
//}
//
//func main() {
//
//	// SDL Initialization
//	var subsystemMask uint32 = sdl.INIT_VIDEO | sdl.INIT_AUDIO
//	if sdl.WasInit(subsystemMask) != subsystemMask {
//		if err := sdl.Init(subsystemMask); err != nil {
//			log.Fatal("Error initializing SDL:", err)
//		}
//	}
//
//	// Failed
//	// 01-special.gb
//	// 02-interrupts.gb
//	// 03-op sp,hl.gb
//	// 07-jr,jp,call,ret,rst.gb
//	// 08-misc instrs.gb
//	// 11-op a,(hl).gb
//
//	// Passed
//	// 04-op r,imm.gb
//	// 05-op rp.gb
//	// 06-ld r,r.gb
//	// 09-op r,r.gb
//	// 10-bit ops.gb
//	//emulator, err := newEmulator("./assets/roms/gb-test-roms/cpu_instrs/individual/01-special.gb", "save10.bin")
//	//emulator, err := newEmulator("./assets/roms/pokeblue.bin", "save11.bin", "./assets/roms/DMG_ROM.bin")
//	//emulator, err := newEmulator("./assets/roms/gb-test-roms/cpu_instrs/individual/01-special.gb", "save10.bin", "")
//	emulator, err := newEmulator("./assets/roms/pokeblue.bin", "s1.bin", "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer emulator.Destroy()
//
//	err = emulator.BessLoad("rom.s3")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	emulator.Run()
//}
