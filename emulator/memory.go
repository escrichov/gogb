package emulator

import "github.com/veandco/go-sdl2/sdl"

func (e *Emulator) mem8(addr uint16, val uint8, write bool) uint8 {
    e.tick()

    switch addr >> 13 {
    case 1: // 0x2000 - 0x3FFF
        if write {
            if e.memoryBankController == 3 {
                // Pokemon Blue uses MBC3, which has the ability to swap 64 different 16KiB banks of ROM
                var romBank uint32 = 1
                if val != 0 {
                    romBank = uint32(val & 0x3F)
                }
                e.rom1Pointer = romBank << 14
            } else if e.memoryBankController == 5 {
                if addr <= 0x2FFF {
                    var romBank = uint32(val & 0x3F)
                    e.rom1Pointer = romBank << 14
                } else {
                    // TODO: Implement set bit 9
                }
            }
        }
        return e.rom0[addr]
    case 0: // 0x0000 - 0x1FFF
        if e.bootRomEnabled && addr <= 0xFF {
            return e.bootRom[addr]
        }
        return e.rom0[addr]
    case 2: // 0x4000 - 0x5FFF
        if e.memoryBankController == 3 || e.memoryBankController == 5 {
            // 4 different of 8KiB banks of External Ram (for a total of 32KiB)
            if write && val <= 3 {
                e.extrambankPointer = uint32(val << 13)
            }
            return e.rom0[e.rom1Pointer+uint32(addr&0x3fff)]
        } else {
            return e.rom0[addr]
        }
    case 3: // 0x6000 - 0x7FFF
        if e.memoryBankController == 3 || e.memoryBankController == 5 {
            return e.rom0[e.rom1Pointer+uint32(addr&0x3fff)]
        } else {
            return e.rom0[addr]
        }
    case 4: // 0x8000 - 0x9FFF
        addr &= 8191
        if write {
            e.videoRam[addr] = val
        }
        return e.videoRam[addr]

    case 5: // 0xA000 - 0xBFFF
        if e.memoryBankController == 3 || e.memoryBankController == 5 {
            addr &= 0x1fff
            if write {
                e.extrambank[e.extrambankPointer+uint32(addr)] = val
            }
            return e.extrambank[e.extrambankPointer+uint32(addr)]
        } else {
            return 0
        }
    case 7: // 0xE000 - 0xFFFF
        if addr >= 0xFE00 {
            if write {
                if addr == 0xFF46 {
                    for y := WIDTH - 1; y >= 0; y-- {
                        e.io[y] = e.read8(uint16(val)<<8 | uint16(y))
                    }
                } else if addr == 0xFF40 {
                    e.SetLCDC(val)
                } else if addr == 0xFF50 {
                    e.bootRomEnabled = false
                }
                ioAddr := addr & 0x1ff
                e.io[ioAddr] = val
            }

            if addr == 0xff00 {
                if (^e.io[256] & 16) != 0 {
                    return ^(16 + e.keyboardState[sdl.SCANCODE_DOWN]*8 +
                        e.keyboardState[sdl.SCANCODE_UP]*4 +
                        e.keyboardState[sdl.SCANCODE_LEFT]*2 +
                        e.keyboardState[sdl.SCANCODE_RIGHT])
                }
                if (^e.io[256] & 32) != 0 {
                    return ^(32 + e.keyboardState[sdl.SCANCODE_RETURN]*8 +
                        e.keyboardState[sdl.SCANCODE_TAB]*4 +
                        e.keyboardState[sdl.SCANCODE_Z]*2 +
                        e.keyboardState[sdl.SCANCODE_X])
                }
                return 0xFF
            }
            ioAddr := addr & 0x1ff
            return e.io[ioAddr]
        } else { // Echo internal RAM
            addr &= 0x3fff
            if write {
                e.workRam[addr] = val
            }
            return e.workRam[addr]
        }
    case 6: // 0xC000 - 0xDFFF, Internal RAM
        addr &= 0x3fff
        if write {
            e.workRam[addr] = val
        }
        return e.workRam[addr]
    }

    return 0
}

func (e *Emulator) read16(addr uint16) uint16 {
    tmp8 := e.mem8(addr, 0, false)
    addr++
    result := e.mem8(addr, 0, false)
    addr++
    return uint16(result)<<8 | uint16(tmp8)
}

func (e *Emulator) read8(addr uint16) uint8 {
    return e.mem8(addr, 0, false)
}

func (e *Emulator) write16(addr uint16, val uint16) {
    e.mem8(addr, uint8(val>>8), true)
    addr++
    e.mem8(addr, uint8(val), true)
}

func (e *Emulator) write8(addr uint16, val uint8) {
    e.mem8(addr, val, true)
}
