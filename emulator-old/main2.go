package emulator_old

//import (
//	"emulator-go/emulator/gb/utils"
//	"errors"
//	"fmt"
//	"github.com/veandco/go-sdl2/sdl"
//	"log"
//	"os"
//	"syscall"
//	"unsafe"
//)
//
//const WIDTH = 160
//const HEIGHT = 144
//
//var cycles, prevCycles uint16
//var totalCycles uint64
//var workRam [16384]uint8
//var videoRam [8192]uint8
//var io [512]uint8
//var extrambank *[32768]uint8
//var ppuDot = 32
//var rom0 []byte
//var extrambankPointer uint32
//var rom1Pointer uint32 = 32768
//var keyboardState []uint8
//var frameBuffer [WIDTH * HEIGHT]int32
//
//var palette = []int32{-1, -23197, -65536, -1 << 24, -1, -8092417, -12961132, -1 << 24}
//
////var palette = []int32{-23197, -23197, -23197, -23197, -23197, -23197, -23197, -23197}
////var palette = []int32{-12961132, -12961132, -12961132, -12961132, -12961132, -12961132, -12961132, -12961132}
////var palette = []int32{-8092417, -8092417, -8092417, -8092417, -8092417, -8092417, -8092417, -8092417}
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
//func (cpu *CPU) popPC() uint8 {
//	result := read8(cpu.PC)
//	cpu.PC++
//	return result
//}
//
//func (cpu *CPU) popPC16() uint16 {
//	result := read16(cpu.PC)
//	cpu.PC += 2
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
//		return cpu.BC.Get()
//	case 1:
//		return cpu.DE.Get()
//	case 2:
//		return cpu.HL.Get()
//	case 3:
//		return cpu.SP.Get()
//	default:
//		return 0
//	}
//}
//
//func (cpu *CPU) r16Group1Set(number uint8, val uint16) {
//	switch number {
//	case 0:
//		cpu.BC.Set(val)
//	case 1:
//		cpu.DE.Set(val)
//	case 2:
//		cpu.HL.Set(val)
//	case 3:
//		cpu.SP.Set(val)
//	default:
//	}
//}
//
//func (cpu *CPU) r16Group2Get(number uint8) uint16 {
//	switch number {
//	case 0:
//		return cpu.BC.Get()
//	case 1:
//		return cpu.DE.Get()
//	case 2:
//		value := cpu.HL.Get()
//		cpu.HL.Set(value + 1)
//		return value
//	case 3:
//		value := cpu.HL.Get()
//		cpu.HL.Set(value - 1)
//		return value
//	default:
//		return 0
//	}
//}
//
//func (cpu *CPU) r16Group3Get(number uint8) uint16 {
//	switch number {
//	case 0:
//		return cpu.BC.Get()
//	case 1:
//		return cpu.DE.Get()
//	case 2:
//		return cpu.HL.Get()
//	case 3:
//		return cpu.AF.Get()
//	default:
//		return 0
//	}
//}
//
//func (cpu *CPU) r16Group3Set(number uint8, val uint16) {
//	switch number {
//	case 0:
//		cpu.BC.Set(val)
//	case 1:
//		cpu.DE.Set(val)
//	case 2:
//		cpu.HL.Set(val)
//	case 3:
//		cpu.AF.Set(val)
//	default:
//	}
//}
//
//func (cpu *CPU) r8Get(number uint8) uint8 {
//	switch number {
//	case 0:
//		return cpu.BC.GetHi()
//	case 1:
//		return cpu.BC.GetLo()
//	case 2:
//		return cpu.DE.GetHi()
//	case 3:
//		return cpu.DE.GetLo()
//	case 4:
//		return cpu.HL.GetHi()
//	case 5:
//		return cpu.HL.GetLo()
//	case 6:
//		return read8(cpu.HL.Get())
//	case 7:
//		return cpu.AF.GetHi()
//	default:
//		return 0
//	}
//}
//
//func (cpu *CPU) r8Set(number uint8, val uint8) {
//	switch number {
//	case 0:
//		cpu.BC.SetHi(val)
//	case 1:
//		cpu.BC.SetLo(val)
//	case 2:
//		cpu.DE.SetHi(val)
//	case 3:
//		cpu.DE.SetLo(val)
//	case 4:
//		cpu.HL.SetHi(val)
//	case 5:
//		cpu.HL.SetLo(val)
//	case 6:
//		write8(cpu.HL.Get(), val)
//	case 7:
//		cpu.AF.SetHi(val)
//	default:
//	}
//}
//
//func tick() {
//	cycles += 4
//	totalCycles += 4
//}
//
//func mem8(addr uint16, val uint8, write bool) uint8 {
//	tick()
//	switch addr >> 13 {
//	case 1:
//		if write {
//			// Pokemon Blue uses MBC3, which has the ability to swap 64 different 16KiB banks of ROM
//			var romBank uint32 = 1
//			if val != 0 {
//				romBank = uint32(val & 63)
//			}
//			rom1Pointer = romBank << 14
//		}
//		return rom0[addr]
//	case 0:
//		return rom0[addr]
//	case 2:
//		// 4 different of 8KiB banks of External Ram (for a total of 32KiB)
//		if write && val <= 3 {
//			extrambankPointer = uint32(val << 13)
//		}
//		return rom0[rom1Pointer+uint32(addr&16383)]
//	case 3:
//		return rom0[rom1Pointer+uint32(addr&16383)]
//	case 4:
//		addr &= 8191
//		if write {
//			videoRam[addr] = val
//		}
//		return videoRam[addr]
//
//	case 5:
//		addr &= 8191
//		if write {
//			extrambank[extrambankPointer+uint32(addr)] = val
//		}
//		return extrambank[extrambankPointer+uint32(addr)]
//
//	case 7:
//		if addr >= 65024 {
//			if write {
//				if addr == 65350 {
//					for y := WIDTH - 1; y >= 0; y-- {
//						io[y] = read8(uint16(val)<<8 | uint16(y))
//					}
//				}
//				ioAddr := addr & 511
//				io[ioAddr] = val
//			}
//
//			if addr == 65280 {
//				if (^io[256] & 16) != 0 {
//					return ^(16 + keyboardState[sdl.SCANCODE_DOWN]*8 +
//						keyboardState[sdl.SCANCODE_UP]*4 +
//						keyboardState[sdl.SCANCODE_LEFT]*2 +
//						keyboardState[sdl.SCANCODE_RIGHT])
//				}
//				if (^io[256] & 32) != 0 {
//					return ^(32 + keyboardState[sdl.SCANCODE_RETURN]*8 +
//						keyboardState[sdl.SCANCODE_TAB]*4 +
//						keyboardState[sdl.SCANCODE_Z]*2 +
//						keyboardState[sdl.SCANCODE_X])
//				}
//				return 255
//			}
//			ioAddr := addr & 511
//			return io[ioAddr]
//		}
//	case 6:
//		addr &= 16383
//		if write {
//			workRam[addr] = val
//		}
//		return workRam[addr]
//	}
//
//	return 0
//}
//
//func getColor(tile, yOffset, xOffset int) uint8 {
//	videoRamIndex := tile*16 + yOffset*2
//	tileData := videoRam[videoRamIndex]
//	tileData1 := videoRam[videoRamIndex+1]
//	return ((tileData1>>xOffset)%2)*2 + (tileData>>xOffset)%2
//}
//
//func read16(addr uint16) uint16 {
//	tmp8 := mem8(addr, 0, false)
//	addr++
//	result := mem8(addr, 0, false)
//	addr++
//	return uint16(result)<<8 | uint16(tmp8)
//}
//
//func read8(addr uint16) uint8 {
//	return mem8(addr, 0, false)
//}
//
//func write16(addr uint16, val uint16) {
//	mem8(addr, uint8(val>>8), true)
//	addr++
//	mem8(addr, uint8(val), true)
//}
//
//func write8(addr uint16, val uint8) {
//	mem8(addr, val, true)
//}
//
//func (cpu *CPU) push(val uint16) {
//	sp := cpu.SP.Get()
//	sp--
//	write8(sp, uint8(val>>8))
//	sp--
//	write8(sp, uint8(val))
//	cpu.SP.Set(sp)
//
//	tick()
//}
//
//func (cpu *CPU) pop() uint16 {
//	sp := cpu.SP.Get()
//	result := read16(sp)
//	cpu.SP.Set(sp + 2)
//
//	return result
//}
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
//func SetIF(value uint8) {
//	io[271] = value
//}
//
//func GetIF() uint8 {
//	return io[271]
//}
//
//func SetLCDC(value uint8) {
//	io[320] = value
//}
//
//func GetLCDC() uint8 {
//	return io[320]
//}
//
//func SetLY(value uint8) {
//	io[324] = value
//}
//
//func GetLY() uint8 {
//	return io[324]
//}
//
//func SetDIV(value uint16) {
//	io[260] = uint8(value >> 8)
//	io[259] = uint8(value & 0xFF)
//}
//
//func GetDIV() uint16 {
//	return uint16(io[260])<<8 | uint16(io[259])
//}
//
//func main() {
//
//	var err error
//	rom0, err = os.ReadFile("./assets/roms/pokered.bin")
//	if err != nil {
//		fmt.Println("Boot Rom file not found:", err)
//		return
//	}
//
//	t := int(unsafe.Sizeof(uint8(8))) * 32768
//	var mapFile *os.File
//	if _, err := os.Stat("save.bin"); errors.Is(err, os.ErrNotExist) {
//		mapFile, err = os.Create("save.bin")
//		if err != nil {
//			log.Fatal("Error opening file: ", err)
//		}
//		_, err = mapFile.Seek(int64(t-1), 0)
//		if err != nil {
//			log.Fatal("Error opening file: ", err)
//		}
//		_, err = mapFile.Write([]byte(" "))
//		if err != nil {
//			log.Fatal("Error writing file: ", err)
//		}
//	} else {
//		mapFile, err = os.OpenFile("save.bin", os.O_RDWR|os.O_CREATE, 0666)
//		if err != nil {
//			log.Fatal("Error opening file: ", err)
//		}
//	}
//	mmap, err := syscall.Mmap(int(mapFile.Fd()), 0, int(t), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	extrambank = (*[32768]uint8)(unsafe.Pointer(&mmap[0]))
//
//	var IME uint8
//	var halt uint8
//	var cpu CPU
//
//	// Initialization
//	cpu.PC = 256
//	SetLCDC(145)
//	SetLY(0)
//	SetDIV(44032)
//	cpu.SP.Set(65534)
//	cpu.AF.SetHi(1)
//	cpu.AF.SetLo(176)
//	cpu.BC.SetHi(1)
//	cpu.BC.SetLo(0)
//	cpu.DE.SetHi(0)
//	cpu.DE.SetLo(216)
//	cpu.HL.SetHi(1)
//	cpu.HL.SetLo(77)
//
//	var numInstructions uint64 = 0
//	var vsyncEnabled = true
//
//	// Debug
//	// create a file and check for errors
//	debugFile, err := os.Create("debug_go.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// close the file
//	defer debugFile.Close()
//
//	// Framebuffer set to black
//	for i, _ := range frameBuffer {
//		frameBuffer[i] = 0
//	}
//
//	// SDL Initialization
//	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
//		log.Fatal("Error initializing SDL:", err)
//		return
//	}
//
//	window, err := sdl.CreateWindow(
//		"GB",
//		sdl.WINDOWPOS_CENTERED,
//		sdl.WINDOWPOS_CENTERED,
//		WIDTH*4,
//		HEIGHT*4,
//		sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
//	if err != nil {
//		log.Fatal("Error creating window:", err)
//	}
//
//	// sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_ACCELERATED
//	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_ACCELERATED)
//	if err != nil {
//		log.Fatal("Error creating renderer:", err)
//	}
//
//	// Creating a SDL texture that is used to display the color buffer
//	texture, err := renderer.CreateTexture(
//		uint32(sdl.PIXELFORMAT_RGBA32),
//		sdl.TEXTUREACCESS_STREAMING,
//		WIDTH,
//		HEIGHT,
//	)
//	if err != nil {
//		log.Fatal("Error creating texture:", err)
//	}
//
//	keyboardState = sdl.GetKeyboardState()
//	for {
//		prevCycles = cycles
//		if (IME & GetIF() & io[511]) != 0 {
//			SetIF(0)
//			halt = 0
//			IME = 0
//			tick()
//			tick()
//			cpu.push(cpu.PC)
//			cpu.PC = 64
//		} else if halt != 0 {
//			tick()
//		} else {
//			// CPU Switch
//			opcode := cpu.popPC()
//			numInstructions++
//			//_, err := debugFile.WriteString(fmt.Sprintf("OP: %d PC: %d Cycles: %d A: %d F: %d DIV: %d\n", opcode, cpu.PC, totalCycles, cpu.AF.GetHi(), cpu.AF.GetLo(), GetDIV()))
//			//if err != nil {
//			//	log.Fatal(err)
//			//}
//
//			switch opcode {
//			case 0: // NOP
//			case 8: // LD (u16), SP
//				write16(cpu.popPC16(), cpu.SP.Get())
//			case 16: // STOP (TODO: Not implemented)
//				halt = 1
//				fmt.Println("HALT = 1")
//			case 24: // JR (unconditional)
//				i8 := int8(cpu.popPC())
//				addr := int32(cpu.PC) + int32(i8)
//				cpu.PC = uint16(addr)
//				tick()
//			case 32, 40, 48, 56: // JR (conditional)
//				i8 := int8(cpu.popPC())
//				if cpu.checkCondition((opcode >> 3) & 0x3) {
//					addr := int32(cpu.PC) + int32(i8)
//					cpu.PC = uint16(addr)
//					tick()
//				}
//			case 1, 17, 33, 49: // LD r16, u16
//				u16 := cpu.popPC16()
//				number := (opcode >> 4) & 0x3
//				cpu.r16Group1Set(number, u16)
//			case 9, 25, 41, 57: // ADD HL, r16
//				number := (opcode >> 4) & 0x3
//				r16 := cpu.r16Group1Get(number)
//				hl := cpu.HL.Get()
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag((hl%4096 + r16%4096) > 4095)
//				cpu.SetCarryFlag((uint32(hl) + uint32(r16)) > 65535)
//				cpu.HL.Set(hl + r16)
//				tick()
//			case 2, 18, 34, 50: // LD (r16), A
//				number := (opcode >> 4) & 0x3
//				write8(cpu.r16Group2Get(number), cpu.AF.GetHi())
//			case 10, 26, 42, 58: // LD A, (r16)
//				number := (opcode >> 4) & 0x3
//				cpu.AF.SetHi(read8(cpu.r16Group2Get(number)))
//			case 3, 19, 35, 51: // INC r16
//				number := (opcode >> 4) & 0x3
//				cpu.r16Group1Set(number, cpu.r16Group1Get(number)+1)
//				tick()
//			case 11, 27, 43, 59: // DEC r16
//				number := (opcode >> 4) & 0x3
//				r16 := cpu.r16Group1Get(number)
//				cpu.r16Group1Set(number, r16-1)
//				tick()
//			case 4, 12, 20, 28, 36, 44, 52, 60: // INC r8
//				number := (opcode >> 3) & 0x7
//				r8 := cpu.r8Get(number)
//				result := r8 + 1
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(result&15 == 0)
//				cpu.r8Set(number, r8+1)
//			case 5, 13, 21, 29, 37, 45, 53, 61: // DEC r8
//				number := (opcode >> 3) & 0x7
//				r8 := cpu.r8Get(number)
//				result := r8 - 1
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(true)
//				cpu.SetHalfCarryFlag((result+1)&15 == 0)
//				cpu.r8Set(number, r8-1)
//				// TODO: Remove this
//				if opcode == 53 {
//					tick()
//				}
//			case 6, 14, 22, 30, 38, 46, 54, 62: // LD r8, u8
//				u8 := cpu.popPC()
//				number := (opcode >> 3) & 0x7
//				cpu.r8Set(number, u8)
//				// TODO: Remove this
//				if opcode == 54 {
//					tick()
//				}
//			case 7: // RLCA
//				a := cpu.AF.GetHi()
//				bit7 := utils.GetBit(a, 7)
//				result := (a << 1) | BoolToUint8(bit7)
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(false)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(bit7)
//			case 15: // RRCA
//				a := cpu.AF.GetHi()
//				bit0 := utils.GetBit(a, 0)
//				result := (a >> 1) | (BoolToUint8(bit0) << 7)
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(false)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(bit0)
//			case 23: // RLA
//				a := cpu.AF.GetHi()
//				bit7 := utils.GetBit(a, 7)
//				carry := BoolToUint8(cpu.GetCarryFlag())
//				result := (a << 1) | carry
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(false)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(bit7)
//			case 31: // RRA
//				a := cpu.AF.GetHi()
//				bit0 := utils.GetBit(a, 0)
//				carry := BoolToUint8(cpu.GetCarryFlag())
//				result := (a >> 1) | (carry << 7)
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(false)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(bit0)
//			case 39: // DAA
//				tmp8 := uint8(0)
//				carry := false
//				if cpu.GetHalfCarryFlag() || !cpu.GetSubtractFlag() && cpu.AF.GetHi()%16 > 9 {
//					tmp8 = 6
//				}
//				if cpu.GetCarryFlag() || !cpu.GetSubtractFlag() && cpu.AF.GetHi() > 153 {
//					tmp8 |= 96
//					carry = true
//				}
//
//				a := cpu.AF.GetHi()
//				result := a
//				if cpu.GetSubtractFlag() {
//					result -= tmp8
//				} else {
//					result += tmp8
//				}
//
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(carry)
//			case 47: // CPL
//				a := cpu.AF.GetHi()
//				result := a ^ 0xFF
//				cpu.AF.SetHi(result)
//				cpu.SetSubtractFlag(true)
//				cpu.SetHalfCarryFlag(true)
//			case 55: // SCF
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(true)
//			case 63: // CCF
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(!cpu.GetCarryFlag())
//			case 118: // HALT
//				halt = 1
//			case 64, 65, 66, 67, 68, 69, 70, 71, // LD r8, r8
//				72, 73, 74, 75, 76, 77, 78, 79,
//				80, 81, 82, 83, 84, 85, 86, 87,
//				88, 89, 90, 91, 92, 93, 94, 95,
//				96, 97, 98, 99, 100, 101, 102, 103,
//				104, 105, 106, 107, 108, 109, 110, 111,
//				112, 113, 114, 115, 116, 117, 119, // 118 is missing because is defined HALT,
//				120, 121, 122, 123, 124, 125, 126, 127:
//				numberSource := opcode & 0x7
//				numberDestination := (opcode >> 3) & 0x7
//				cpu.r8Set(numberDestination, cpu.r8Get(numberSource))
//			case 128, 129, 130, 131, 132, 133, 134, 135: // ADD A, r8
//				number := opcode & 0x7
//				r8 := cpu.r8Get(number)
//				a := cpu.AF.GetHi()
//				result := a + r8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(a%16+r8%16 > 15)
//				cpu.SetCarryFlag(uint16(a)+uint16(r8) > 255)
//			case 136, 137, 138, 139, 140, 141, 142, 143: // ADC A, r8
//				number := opcode & 0x7
//				r8 := cpu.r8Get(number)
//				carry := BoolToUint8(cpu.GetCarryFlag())
//				a := cpu.AF.GetHi()
//				result := a + r8 + carry
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(a%16+r8%16+carry > 15)
//				cpu.SetCarryFlag(uint16(a)+uint16(r8)+uint16(carry) > 255)
//			case 144, 145, 146, 147, 148, 149, 150, 151: // SUB A, r8
//				number := opcode & 0x7
//				r8 := cpu.r8Get(number)
//				a := cpu.AF.GetHi()
//				result := a - r8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(true)
//				cpu.SetHalfCarryFlag(a%16-r8%16 > 15)
//				cpu.SetCarryFlag(uint16(a)-uint16(r8) > 255)
//			case 152, 153, 154, 155, 156, 157, 158, 159: // SBC A, r8
//				number := opcode & 0x7
//				r8 := cpu.r8Get(number)
//				carry := BoolToUint8(cpu.GetCarryFlag())
//				a := cpu.AF.GetHi()
//				result := a - r8 - carry
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(true)
//				cpu.SetHalfCarryFlag(a%16-r8%16-carry > 15)
//				cpu.SetCarryFlag(uint16(a)-uint16(r8)-uint16(carry) > 255)
//			case 160, 161, 162, 163, 164, 165, 166, 167: // AND A, r8
//				number := opcode & 0x7
//				r8 := cpu.r8Get(number)
//				a := cpu.AF.GetHi()
//				result := a & r8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(true)
//				cpu.SetCarryFlag(false)
//			case 168, 169, 170, 171, 172, 173, 174, 175: // XOR A, r8
//				number := opcode & 0x7
//				r8 := cpu.r8Get(number)
//				a := cpu.AF.GetHi()
//				result := a ^ r8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(false)
//			case 176, 177, 178, 179, 180, 181, 182, 183: // OR A, r8
//				number := opcode & 0x7
//				r8 := cpu.r8Get(number)
//				a := cpu.AF.GetHi()
//				result := a | r8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(false)
//			case 184, 185, 186, 187, 188, 189, 190, 191: // CP A, r8
//				number := opcode & 0x7
//				r8 := cpu.r8Get(number)
//				a := cpu.AF.GetHi()
//				result := a - r8
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(true)
//				cpu.SetHalfCarryFlag(a%16-r8%16 > 15)
//				cpu.SetCarryFlag(uint16(a)-uint16(r8) > 255)
//			case 192, 200, 208, 216: // RET condition
//				tick()
//				if cpu.checkCondition((opcode >> 3) & 0x3) {
//					addr := cpu.pop()
//					cpu.PC = addr
//				}
//			case 224: // LD (FF00 + u8), A
//				u8 := cpu.popPC()
//				addr := 0xFF00 + uint16(u8)
//				write8(addr, cpu.AF.GetHi())
//			case 232: // ADD SP, i8
//				i8 := int8(cpu.popPC())
//				sp := cpu.SP.Get()
//				result := int32(sp) + int32(i8)
//				cpu.SP.Set(uint16(result))
//				cpu.SetZeroFlag(false)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(int16(i8%16)+int16(uint8(i8)%16) > 15)
//				cpu.SetCarryFlag(int16(sp)+int16(i8) > 255)
//				tick()
//				tick()
//			case 240: // LD A, (FF00 + u8)
//				u8 := cpu.popPC()
//				addr := 0xFF00 + uint16(u8)
//				cpu.AF.SetHi(read8(addr))
//			case 248: // LD HL, SP + i8
//				i8 := int8(cpu.popPC())
//				sp := cpu.SP.Get()
//				result := int32(sp) + int32(i8)
//				cpu.HL.Set(uint16(result))
//				cpu.SetZeroFlag(false)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(int16(sp%16)+int16(i8%16) > 15)
//				cpu.SetCarryFlag(uint16(uint8(sp))+uint16(i8) > 255)
//				tick()
//			case 193, 209, 225, 241: // POP r16
//				number := (opcode >> 4) & 0x3
//				cpu.r16Group3Set(number, cpu.pop())
//			case 201: // RET
//				pc := cpu.pop()
//				cpu.PC = pc
//				tick()
//			case 217: // RETI
//				IME = 1
//				pc := cpu.pop()
//				cpu.PC = pc
//				tick()
//			case 233: // JP HL
//				cpu.PC = cpu.HL.Get()
//			case 249: // LD SP, HL
//				cpu.SP.Set(cpu.HL.Get())
//				tick()
//			case 194, 202, 210, 218: // JP condition
//				addr := cpu.popPC16()
//				if cpu.checkCondition((opcode >> 3) & 0x3) {
//					cpu.PC = addr
//					tick()
//				}
//			case 226: // LD (FF00+C), A
//				addr := 0xFF00 + uint16(cpu.BC.GetLo())
//				write8(addr, cpu.AF.GetHi())
//			case 234: // LD (u16), A
//				u16 := cpu.popPC16()
//				write8(u16, cpu.AF.GetHi())
//			case 242: // LD A, (0xFF00+C)
//				addr := 0xFF00 + uint16(cpu.BC.GetLo())
//				cpu.AF.SetHi(read8(addr))
//			case 250: // LD A, (u16)
//				u16 := cpu.popPC16()
//				cpu.AF.SetHi(read8(u16))
//			case 195: // JP u16
//				cpu.PC = cpu.popPC16()
//				tick()
//			case 243: // DI
//				IME = 0
//			case 251: // EI
//				IME = 1
//			case 196, 204, 212, 220: // CALL condition
//				u16 := cpu.popPC16()
//				if cpu.checkCondition((opcode >> 3) & 0x3) {
//					cpu.push(cpu.PC)
//					cpu.PC = u16
//				}
//			case 197, 213, 229, 245: // PUSH r16
//				number := (opcode >> 4) & 0x3
//				cpu.push(cpu.r16Group3Get(number))
//			case 205: // CALL u16
//				u16 := cpu.popPC16()
//				cpu.push(cpu.PC)
//				cpu.PC = u16
//			case 198: // ADD A, u8
//				u8 := cpu.popPC()
//				a := cpu.AF.GetHi()
//				result := a + u8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(a%16+u8%16 > 15)
//				cpu.SetCarryFlag(uint16(a)+uint16(u8) > 255)
//			case 206: // ADC A, u8
//				u8 := cpu.popPC()
//				carry := BoolToUint8(cpu.GetCarryFlag())
//				a := cpu.AF.GetHi()
//				result := a + u8 + carry
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(a%16+u8%16+carry > 15)
//				cpu.SetCarryFlag(uint16(a)+uint16(u8)+uint16(carry) > 255)
//			case 214: // SUB A, u8
//				u8 := cpu.popPC()
//				a := cpu.AF.GetHi()
//				result := a - u8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(true)
//				cpu.SetHalfCarryFlag(a%16-u8%16 > 15)
//				cpu.SetCarryFlag(uint16(a)-uint16(u8) > 255)
//			case 222: // SBC A, u8
//				u8 := cpu.popPC()
//				carry := BoolToUint8(cpu.GetCarryFlag())
//				a := cpu.AF.GetHi()
//				result := a - u8 - carry
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(true)
//				cpu.SetHalfCarryFlag(a%16-u8%16-carry > 15)
//				cpu.SetCarryFlag(uint16(a)-uint16(u8)-uint16(carry) > 255)
//			case 230: // AND A, u8
//				u8 := cpu.popPC()
//				a := cpu.AF.GetHi()
//				result := a & u8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(true)
//				cpu.SetCarryFlag(false)
//			case 238: // XOR A, u8
//				u8 := cpu.popPC()
//				a := cpu.AF.GetHi()
//				result := a ^ u8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(false)
//			case 246: // OR A, u8
//				u8 := cpu.popPC()
//				a := cpu.AF.GetHi()
//				result := a | u8
//				cpu.AF.SetHi(result)
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(false)
//				cpu.SetHalfCarryFlag(false)
//				cpu.SetCarryFlag(false)
//			case 254: // CP A, u8
//				u8 := cpu.popPC()
//				a := cpu.AF.GetHi()
//				result := a - u8
//				cpu.SetZeroFlag(result == 0)
//				cpu.SetSubtractFlag(true)
//				cpu.SetHalfCarryFlag(a%16-u8%16 > 15)
//				cpu.SetCarryFlag(uint16(a)-uint16(u8) > 255)
//			case 199, 207, 215, 223, 231, 239, 247, 255: // RST (Call to 00EXP000)
//				// TODO: RST (Call to 00EXP000) Not implemented
//			case 0xCB:
//				opcode := cpu.popPC()
//				switch opcode {
//				case 0, 1, 2, 3, 4, 5, 6, 7: // RLC
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					bit7 := utils.GetBit(r8, 7)
//					result := (r8 << 1) | BoolToUint8(bit7)
//					cpu.r8Set(number, result)
//					cpu.SetZeroFlag(result == 0)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(false)
//					cpu.SetCarryFlag(bit7)
//				case 8, 9, 10, 11, 12, 13, 14, 15: // RRC
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					bit0 := utils.GetBit(r8, 0)
//					result := (r8 >> 1) | (BoolToUint8(bit0) << 7)
//					cpu.r8Set(number, result)
//					cpu.SetZeroFlag(result == 0)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(false)
//					cpu.SetCarryFlag(bit0)
//					// TODO: Remove this
//					if opcode == 14 {
//						tick()
//						tick()
//						tick()
//					}
//				case 16, 17, 18, 19, 20, 21, 22, 23: // RL
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					bit7 := utils.GetBit(r8, 7)
//					carry := BoolToUint8(cpu.GetCarryFlag())
//					result := (r8 << 1) | carry
//					cpu.r8Set(number, result)
//					cpu.SetZeroFlag(result == 0)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(false)
//					cpu.SetCarryFlag(bit7)
//				case 24, 25, 26, 27, 28, 29, 30, 31: // RR
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					bit0 := utils.GetBit(r8, 0)
//					carry := BoolToUint8(cpu.GetCarryFlag())
//					result := (r8 >> 1) | (carry << 7)
//					cpu.r8Set(number, result)
//					cpu.SetZeroFlag(result == 0)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(false)
//					cpu.SetCarryFlag(bit0)
//				case 32, 33, 34, 35, 36, 37, 38, 39: // SLA
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					bit7 := utils.GetBit(r8, 7)
//					result := r8 << 1
//					cpu.r8Set(number, result)
//					cpu.SetZeroFlag(result == 0)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(false)
//					cpu.SetCarryFlag(bit7)
//				case 40, 41, 42, 43, 44, 45, 46, 47: // SRA
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					bit0 := utils.GetBit(r8, 0)
//					bit7 := utils.GetBit(r8, 7)
//					result := r8>>1 | (BoolToUint8(bit7) << 7)
//					cpu.r8Set(number, result)
//					cpu.SetZeroFlag(result == 0)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(false)
//					cpu.SetCarryFlag(bit0)
//				case 48, 49, 50, 51, 52, 53, 54, 55: // SWAP
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					lower := r8 & 0xF
//					upper := r8 >> 4
//					result := (lower << 4) | upper
//					cpu.r8Set(number, result)
//					cpu.SetZeroFlag(result == 0)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(false)
//					cpu.SetCarryFlag(false)
//				case 56, 57, 58, 59, 60, 61, 62, 63: // SRL
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					bit0 := utils.GetBit(r8, 0)
//					result := r8 >> 1
//					cpu.r8Set(number, result)
//					cpu.SetZeroFlag(result == 0)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(false)
//					cpu.SetCarryFlag(bit0)
//				case 64, 65, 66, 67, 68, 69, 70, 71, // BIT bit, r8
//					72, 73, 74, 75, 76, 77, 78, 79,
//					80, 81, 82, 83, 84, 85, 86, 87,
//					88, 89, 90, 91, 92, 93, 94, 95,
//					96, 97, 98, 99, 100, 101, 102, 103,
//					104, 105, 106, 107, 108, 109, 110, 111,
//					112, 113, 114, 115, 116, 117, 118, 119,
//					120, 121, 122, 123, 124, 125, 126, 127:
//					bit := (opcode >> 3) & 0x7
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					cpu.SetSubtractFlag(false)
//					cpu.SetHalfCarryFlag(true)
//					cpu.SetZeroFlag(!utils.GetBit(r8, int(bit)))
//				case 128, 129, 130, 131, 132, 133, 134, 135, // RES bit, r8
//					136, 137, 138, 139, 140, 141, 142, 143,
//					144, 145, 146, 147, 148, 149, 150, 151,
//					152, 153, 154, 155, 156, 157, 158, 159,
//					160, 161, 162, 163, 164, 165, 166, 167,
//					168, 169, 170, 171, 172, 173, 174, 175,
//					176, 177, 178, 179, 180, 181, 182, 183,
//					184, 185, 186, 187, 188, 189, 190, 191:
//					bit := (opcode >> 3) & 0x7
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					cpu.r8Set(number, SetBit8(r8, bit, false))
//				case 192, 193, 194, 195, 196, 197, 198, 199, // SET bit, r8
//					200, 201, 202, 203, 204, 205, 206, 207,
//					208, 209, 210, 211, 212, 213, 214, 215,
//					216, 217, 218, 219, 220, 221, 222, 223,
//					224, 225, 226, 227, 228, 229, 230, 231,
//					232, 233, 234, 235, 236, 237, 238, 239,
//					240, 241, 242, 243, 244, 245, 246, 247,
//					248, 249, 250, 251, 252, 253, 254, 255:
//					bit := (opcode >> 3) & 0x7
//					number := opcode & 0x7
//					r8 := cpu.r8Get(number)
//					cpu.r8Set(number, SetBit8(r8, bit, true))
//				default:
//					fmt.Println("CB Opcode: ", opcode, " not found")
//					return
//				}
//			default:
//				fmt.Println("Opcode: ", opcode, " not found")
//				return
//			}
//		}
//
//		// PPU
//		div := GetDIV()
//		SetDIV(div + cycles - prevCycles)
//		for ; prevCycles != cycles; prevCycles++ {
//			lcdc := GetLCDC()
//			if (lcdc & 128) != 0 {
//				ppuDot++
//
//				// Render Scanline (Every 256 PPU Dots)
//				if ppuDot == 456 {
//					ly := GetLY()
//
//					// Only render visible lines (up to line 144)
//					if ly < HEIGHT {
//						for tmp := WIDTH - 1; tmp >= 0; tmp-- {
//
//							// IsWindow
//							isWindow := false
//							if (lcdc&32) != 0 && ly >= io[330] && uint8(tmp) >= (io[331]-7) {
//								isWindow = true
//							}
//
//							// xOffset
//							var xOffset uint8
//							if isWindow {
//								xOffset = uint8(tmp) - io[331] + 7
//							} else {
//								xOffset = uint8(tmp) + io[323]
//							}
//
//							// yOffset
//							var yOffset uint8
//							if isWindow {
//								yOffset = ly - io[330]
//							} else {
//								yOffset = ly + io[322]
//							}
//
//							// PaletteIndex
//							var paletteIndex uint16 = 0
//
//							// Tile
//							mask := uint8(8)
//							if isWindow {
//								mask = 64
//							}
//
//							videoRamIndex := uint16(6)
//							if (lcdc & mask) != 0 {
//								videoRamIndex = 7
//							}
//							videoRamIndex = videoRamIndex<<10 | uint16(yOffset)/8*32 + uint16(xOffset)/8
//							var tile = videoRam[videoRamIndex]
//
//							// Color
//							var tileValue int
//							if (lcdc & 16) != 0 {
//								tileValue = int(tile)
//							} else {
//								tileValue = 256 + int(int8(tile))
//							}
//							color := getColor(tileValue, int(yOffset&7), int(7-xOffset&7))
//
//							// Sprites
//							if (lcdc & 2) != 0 {
//								for spriteIndex := uint8(0); spriteIndex < WIDTH; spriteIndex += 4 {
//									spriteX := uint8(tmp) - io[spriteIndex+1] + 8
//									spriteY := ly - io[spriteIndex] + 16
//
//									spriteYOffset := uint8(0)
//									if (io[spriteIndex+3] & 64) != 0 {
//										spriteYOffset = 7
//									}
//									spriteYOffset = spriteY ^ spriteYOffset
//
//									spriteXOffset := uint8(7)
//									if (io[spriteIndex+3] & 32) != 0 {
//										spriteXOffset = 0
//									}
//									spriteXOffset = spriteX ^ spriteXOffset
//
//									spriteColor := getColor(int(io[spriteIndex+2]), int(spriteYOffset), int(spriteXOffset))
//
//									if spriteX < 8 && spriteY < 8 && !((io[spriteIndex+3]&128) != 0 && color != 0) && spriteColor != 0 {
//										color = spriteColor
//										if io[spriteIndex+3]&16 == 0 {
//											paletteIndex = uint16(1)
//										} else {
//											paletteIndex = uint16(2)
//										}
//										break
//									}
//								}
//							}
//
//							paletteIndexValue := uint16((io[327+paletteIndex]>>(2*color))%4) + paletteIndex*4&7
//							frameBufferIndex := uint16(ly)*WIDTH + uint16(tmp)
//							frameBuffer[frameBufferIndex] = palette[paletteIndexValue]
//						}
//					}
//
//					if ly == (HEIGHT - 1) {
//						SetIF(GetIF() | 1)
//
//						// Render Framebuffer
//						if numInstructions != 0 {
//							renderer.Clear()
//							buf := unsafe.Pointer(&frameBuffer[0])
//							framebufferBytes := unsafe.Slice((*byte)(buf), WIDTH*HEIGHT)
//							texture.Update(nil, framebufferBytes, WIDTH*4)
//							renderer.Copy(texture, nil, nil)
//							renderer.Present()
//						}
//
//						// Manage Keyboard Events
//						for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
//							switch event.(type) {
//							case *sdl.KeyboardEvent:
//								keyboardEvent, ok := event.(*sdl.KeyboardEvent)
//								if !ok {
//									continue
//								}
//								if keyboardEvent.Type == sdl.KEYDOWN {
//									switch keyboardEvent.Keysym.Sym {
//									case sdl.K_ESCAPE:
//										texture.Destroy()
//										renderer.Destroy()
//										window.Destroy()
//										sdl.Quit()
//										return
//									case sdl.K_s:
//										vsyncEnabled = !vsyncEnabled
//										renderer.RenderSetVSync(vsyncEnabled)
//									}
//								}
//							case *sdl.QuitEvent:
//								texture.Destroy()
//								renderer.Destroy()
//								window.Destroy()
//								sdl.Quit()
//								return
//							}
//						}
//					}
//
//					// Increment Line
//					SetLY((ly + 1) % 154)
//					ppuDot = 0
//				}
//			} else {
//				SetLY(0)
//				ppuDot = 0
//			}
//		}
//	}
//}
