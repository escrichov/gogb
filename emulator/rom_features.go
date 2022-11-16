package emulator

import (
	"encoding/binary"
	"fmt"
)

type RomFeatures struct {
	Logo                  []byte
	Title                 string
	TitleBytes            []byte
	ManufacturerCode      string
	ManufacturerCodeBytes []byte
	ColorGB               uint8
	LicenseCodeBytes      []byte
	LicenseCodeNew        uint16
	GBSGBIndicator        uint8
	CartridgeType         uint8
	RomSizeByte           uint8
	RamSizeByte           uint8
	DestinationCode       uint8
	DestinationCodeName   string
	LicenseCodeOld        uint8
	MaskROMVersionNumber  uint8
	HeaderChecksum        uint8
	GlobalChecksumBytes   []byte
	GlobalChecksum        uint16
	CartridgeTypeName     string

	LicenseCodeName  string
	LicenseCode      uint16
	GlobalChecksumOk bool
	HeaderChecksumOk bool
	LogoOk           bool
	Filename         string

	SupportColor      bool
	SupportMonochrome bool
	SupportSGB        bool
	PGBMode           bool
}

func getCartridgeTypeName(cartridgeType uint8) string {
	switch cartridgeType {
	case 0x00:
		return "ROM ONLY"
	case 0x01:
		return "MBC1"
	case 0x02:
		return "MBC1+RAM"
	case 0x03:
		return "MBC1+RAM+BATTERY"
	case 0x05:
		return "MBC2"
	case 0x06:
		return "MBC2+BATTERY"
	case 0x08:
		return "ROM+RAM 1"
	case 0x09:
		return "ROM+RAM+BATTERY 1"
	case 0x0B:
		return "MMM01"
	case 0x0C:
		return "MMM01+RAM"
	case 0x0D:
		return "MMM01+RAM+BATTERY"
	case 0x0F:
		return "MBC3+TIMER+BATTERY"
	case 0x10:
		return "MBC3+TIMER+RAM+BATTERY 2"
	case 0x11:
		return "MBC3"
	case 0x12:
		return "MBC3+RAM 2"
	case 0x13:
		return "MBC3+RAM+BATTERY 2"
	case 0x19:
		return "MBC5"
	case 0x1A:
		return "MBC5+RAM"
	case 0x1B:
		return "MBC5+RAM+BATTERY"
	case 0x1C:
		return "MBC5+RUMBLE"
	case 0x1D:
		return "MBC5+RUMBLE+RAM"
	case 0x1E:
		return "MBC5+RUMBLE+RAM+BATTERY"
	case 0x20:
		return "MBC6"
	case 0x22:
		return "MBC7+SENSOR+RUMBLE+RAM+BATTERY"
	case 0xFC:
		return "POCKET CAMERA"
	case 0xFD:
		return "BANDAI TAMA5"
	case 0xFE:
		return "HuC3"
	case 0xFF:
		return "HuC1+RAM+BATTERY"
	default:
		return "Unknown"
	}
}

func getNewLicenseeCodeName(code uint16) string {
	switch code {
	case 0x00:
		return "None"
	case 0x01:
		return "Nintendo R&D1"
	case 0x08:
		return "Capcom"
	case 0x13:
		return "Electronic Arts"
	case 0x18:
		return "Hudson Soft"
	case 0x19:
		return "b-ai"
	case 0x20:
		return "kss"
	case 0x22:
		return "pow"
	case 0x24:
		return "PCM Complete"
	case 0x25:
		return "san-x"
	case 0x28:
		return "Kemco Japan"
	case 0x29:
		return "seta"
	case 0x30:
		return "Viacom"
	case 0x31:
		return "Nintendo"
	case 0x32:
		return "Bandai"
	case 0x33:
		return "Ocean/Acclaim"
	case 0x34:
		return "Konami"
	case 0x35:
		return "Hector"
	case 0x37:
		return "Taito"
	case 0x38:
		return "Hudson"
	case 0x39:
		return "Banpresto"
	case 0x41:
		return "Ubi Soft"
	case 0x42:
		return "Atlus"
	case 0x44:
		return "Malibu"
	case 0x46:
		return "angel"
	case 0x47:
		return "Bullet-Proof"
	case 0x49:
		return "irem"
	case 0x50:
		return "Absolute"
	case 0x51:
		return "Acclaim"
	case 0x52:
		return "Activision"
	case 0x53:
		return "American sammy"
	case 0x54:
		return "Konami"
	case 0x55:
		return "Hi tech entertainment"
	case 0x56:
		return "LJN"
	case 0x57:
		return "Matchbox"
	case 0x58:
		return "Mattel"
	case 0x59:
		return "Milton Bradley"
	case 0x60:
		return "Titus"
	case 0x61:
		return "Virgin"
	case 0x64:
		return "LucasArts"
	case 0x67:
		return "Ocean"
	case 0x69:
		return "Electronic Arts"
	case 0x70:
		return "Infogrames"
	case 0x71:
		return "Interplay"
	case 0x72:
		return "Broderbund"
	case 0x73:
		return "sculptured"
	case 0x75:
		return "sci"
	case 0x78:
		return "THQ"
	case 0x79:
		return "Accolade"
	case 0x80:
		return "misawa"
	case 0x83:
		return "lozc"
	case 0x86:
		return "Tokuma Shoten Intermedia"
	case 0x87:
		return "Tsukuda Original"
	case 0x91:
		return "Chunsoft"
	case 0x92:
		return "Video system"
	case 0x93:
		return "Ocean/Acclaim"
	case 0x95:
		return "Varie"
	case 0x96:
		return "Yonezawa/s’pal"
	case 0x97:
		return "Kaneko"
	case 0x99:
		return "Pack in soft"
	case 0xA4:
		return "Konami (Yu-Gi-Oh!)"
	default:
		return fmt.Sprintf("Unknown (%x)", code)
	}
}

func getOldLicenseeCodeName(code uint8) string {
	switch code {
	case 0x00:
		return "None"
	case 0x01:
		return "Nintendo"
	case 0x08:
		return "Capcom"
	case 0x09:
		return "Hot-B"
	case 0x0A:
		return "Jaleco"
	case 0x0B:
		return "Coconuts Japan"
	case 0x0C:
		return "Elite Systems"
	case 0x13:
		return "EA (Electronic Arts)"
	case 0x18:
		return "Hudsonsoft"
	case 0x19:
		return "ITC Entertainment"
	case 0x1A:
		return "Yanoman"
	case 0x1D:
		return "Japan Clary"
	case 0x1F:
		return "Virgin Interactive"
	case 0x24:
		return "PCM Complete"
	case 0x25:
		return "San-X"
	case 0x28:
		return "Kotobuki Systems"
	case 0x29:
		return "Seta"
	case 0x30:
		return "Infogrames"
	case 0x31:
		return "Nintendo"
	case 0x32:
		return "Bandai"
	case 0x34:
		return "Konami"
	case 0x35:
		return "HectorSoft"
	case 0x38:
		return "Capcom"
	case 0x39:
		return "Banpresto"
	case 0x3C:
		return ".Entertainment i"
	case 0x3E:
		return "Gremlin"
	case 0x41:
		return "Ubisoft"
	case 0x42:
		return "Atlus"
	case 0x44:
		return "Malibu"
	case 0x46:
		return "Angel"
	case 0x47:
		return "Spectrum Holoby"
	case 0x49:
		return "Irem"
	case 0x4A:
		return "Virgin Interactive"
	case 0x4D:
		return "Malibu"
	case 0x4F:
		return "U.S. Gold"
	case 0x50:
		return "Absolute"
	case 0x51:
		return "Acclaim"
	case 0x52:
		return "Activision"
	case 0x53:
		return "American Sammy"
	case 0x54:
		return "GameTek"
	case 0x55:
		return "Park Place"
	case 0x56:
		return "LJN"
	case 0x57:
		return "Matchbox"
	case 0x58:
		return "Mattel"
	case 0x59:
		return "Milton Bradley"
	case 0x5A:
		return "Mindscape"
	case 0x5B:
		return "Romstar"
	case 0x5C:
		return "Naxat Soft"
	case 0x5D:
		return "Tradewest"
	case 0x60:
		return "Titus"
	case 0x61:
		return "Virgin Interactive"
	case 0x67:
		return "Ocean Interactive"
	case 0x69:
		return "EA (Electronic Arts)"
	case 0x6E:
		return "Elite Systems"
	case 0x6F:
		return "Electro Brain"
	case 0x70:
		return "Infogrames"
	case 0x71:
		return "Interplay"
	case 0x72:
		return "Broderbund"
	case 0x73:
		return "Sculptered Soft"
	case 0x75:
		return "The Sales Curve"
	case 0x78:
		return "t.hq"
	case 0x79:
		return "Accolade"
	case 0x7A:
		return "Triffix Entertainment"
	case 0x7C:
		return "Microprose"
	case 0x7F:
		return "Kemco"
	case 0x80:
		return "Misawa Entertainment"
	case 0x83:
		return "Lozc"
	case 0x86:
		return "Tokuma Shoten Intermedia"
	case 0x8B:
		return "Bullet-Proof Software"
	case 0x8C:
		return "Vic Tokai"
	case 0x8E:
		return "Ape"
	case 0x8F:
		return "I’Max"
	case 0x91:
		return "Chunsoft"
	case 0x92:
		return "Video system"
	case 0x93:
		return "Tsubaraya Productions Co."
	case 0x95:
		return "Varie Corporation"
	case 0x96:
		return "Yonezawa/S’Pal"
	case 0x97:
		return "Kaneko"
	case 0x99:
		return "Arc"
	case 0x9A:
		return "Nihon Bussan"
	case 0x9B:
		return "Tecmo"
	case 0x9C:
		return "Tecmo"
	case 0x9D:
		return "Banpresto"
	case 0x9F:
		return "Nova"
	case 0xA1:
		return "Hori Electric"
	case 0xA2:
		return "Bandai"
	case 0xA4:
		return "Konami"
	case 0xA6:
		return "Kawada"
	case 0xA7:
		return "Takara"
	case 0xA9:
		return "Technos Japan"
	case 0xAA:
		return "Broderbund"
	case 0xAC:
		return "Toei Animation"
	case 0xAD:
		return "Toho"
	case 0xB0:
		return "acclaim"
	case 0xB1:
		return "ASCII or Nexsoft"
	case 0xB2:
		return "Bandai"
	case 0xB4:
		return "Square Enix"
	case 0xB6:
		return "HAL Laboratory"
	case 0xB7:
		return "SNK"
	case 0xB9:
		return "Pony Canyon"
	case 0xBA:
		return "Culture Brain"
	case 0xBB:
		return "Sunsoft"
	case 0xBD:
		return "Sony Imagesoft"
	case 0xBF:
		return "Sammy"
	case 0xC0:
		return "Taito"
	case 0xC2:
		return "Kemco"
	case 0xC3:
		return "Squaresoft"
	case 0xC4:
		return "Tokuma Shoten Intermedia"
	case 0xC5:
		return "Data East"
	case 0xC6:
		return "Tonkinhouse"
	case 0xC8:
		return "Koei"
	case 0xC9:
		return "UFL"
	case 0xCA:
		return "Ultra"
	case 0xCB:
		return "Vap"
	case 0xCC:
		return "Use Corporation"
	case 0xCD:
		return "Meldac"
	case 0xCE:
		return ".Pony Canyon or"
	case 0xCF:
		return "Angel"
	case 0xD0:
		return "Taito"
	case 0xD1:
		return "Sofel"
	case 0xD2:
		return "Quest"
	case 0xD3:
		return "Sigma Enterprises"
	case 0xD4:
		return "ASK Kodansha Co."
	case 0xD6:
		return "Naxat Soft"
	case 0xD7:
		return "Copya System"
	case 0xD9:
		return "Banpresto"
	case 0xDA:
		return "Tomy"
	case 0xDB:
		return "LJN"
	case 0xDD:
		return "NCS"
	case 0xDE:
		return "Human"
	case 0xDF:
		return "Altron"
	case 0xE0:
		return "Jaleco"
	case 0xE1:
		return "Towa Chiki"
	case 0xE2:
		return "Yutaka"
	case 0xE3:
		return "Varie"
	case 0xE5:
		return "Epcoh"
	case 0xE7:
		return "Athena"
	case 0xE8:
		return "Asmik ACE Entertainment"
	case 0xE9:
		return "Natsume"
	case 0xEA:
		return "King Records"
	case 0xEB:
		return "Atlus"
	case 0xEC:
		return "Epic/Sony Records"
	case 0xEE:
		return "IGS"
	case 0xF0:
		return "A Wave"
	case 0xF3:
		return "Extreme Entertainment"
	case 0xFF:
		return "LJN"
	default:
		return fmt.Sprintf("Unknown (%x)", code)
	}
}

func ComputeHeaderChecksum(romData []byte) uint8 {
	var checksum uint8 = 0
	for address := uint16(0x0134); address <= 0x014C; address++ {
		checksum = checksum - romData[address] - 1
	}

	return checksum
}

func ComputeGlobalChecksum(romData []byte) uint16 {
	var checksum uint16 = 0
	for _, num := range romData[:0x14E] {
		checksum += uint16(num)
	}

	for _, num := range romData[0x150:] {
		checksum += uint16(num)
	}

	return checksum
}

func isLogoOk(romData []byte) bool {
	logo := []byte{
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
		0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
		0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	}

	for i, realLogo := range logo {
		if romData[0x104+i] != realLogo {
			return false
		}
	}

	return true
}

func parseRomTitle(bs []byte) string {
	var end = 0
	for n, b := range bs {
		end = n
		if b == 0 {
			break
		}
	}
	return string(bs[0:end])
}

func parseRomHeader(romData []byte) (*RomFeatures, error) {
	var romFeatures RomFeatures

	if len(romData) < 0x150 {
		return nil, fmt.Errorf("incorrect rom size %d", len(romData))
	}

	romFeatures.Logo = romData[0x104:0x134]
	romFeatures.Title = parseRomTitle(romData[0x134:0x144])
	romFeatures.TitleBytes = romData[0x134:0x144]
	romFeatures.ManufacturerCodeBytes = romData[0x013F:0x0143]
	romFeatures.ManufacturerCode = string(romFeatures.ManufacturerCodeBytes)
	romFeatures.ColorGB = romData[0x143]
	romFeatures.LicenseCodeBytes = romData[0x144:0x146]
	romFeatures.LicenseCodeNew = (uint16(romData[0x144]) << 8) | uint16(romData[0x145])
	romFeatures.GBSGBIndicator = romData[0x146]
	romFeatures.CartridgeType = romData[0x147]
	romFeatures.RomSizeByte = romData[0x148]
	romFeatures.RamSizeByte = romData[0x149]
	romFeatures.DestinationCode = romData[0x14A]
	romFeatures.LicenseCodeOld = romData[0x14B]
	romFeatures.MaskROMVersionNumber = romData[0x14C]
	romFeatures.HeaderChecksum = romData[0x14D]
	romFeatures.GlobalChecksumBytes = romData[0x14E:0x150]
	romFeatures.GlobalChecksum = binary.BigEndian.Uint16(romFeatures.GlobalChecksumBytes)
	romFeatures.CartridgeTypeName = getCartridgeTypeName(romFeatures.CartridgeType)

	if romFeatures.DestinationCode == 0x00 {
		romFeatures.DestinationCodeName = "Japan (and possibly overseas)"
	} else if romFeatures.DestinationCode == 0x01 {
		romFeatures.DestinationCodeName = "Overseas only"
	} else {
		romFeatures.DestinationCodeName = "Unknown"
	}

	if romFeatures.LicenseCodeOld == 0x33 {
		romFeatures.LicenseCode = romFeatures.LicenseCodeNew
		romFeatures.LicenseCodeName = getNewLicenseeCodeName(romFeatures.LicenseCodeNew)
	} else {
		romFeatures.LicenseCode = uint16(romFeatures.LicenseCodeOld)
		romFeatures.LicenseCodeName = getOldLicenseeCodeName(romFeatures.LicenseCodeOld)
	}

	romFeatures.HeaderChecksumOk = false
	if ComputeHeaderChecksum(romData) == romFeatures.HeaderChecksum {
		romFeatures.HeaderChecksumOk = true
	}

	romFeatures.GlobalChecksumOk = false
	if ComputeGlobalChecksum(romData) == romFeatures.GlobalChecksum {
		romFeatures.GlobalChecksumOk = true
	}

	romFeatures.LogoOk = isLogoOk(romData)

	romFeatures.SupportMonochrome = true
	if GetBit(romFeatures.ColorGB, 7) {
		// Values with bit 7 and either bit 2 or 3 set will switch the Game Boy into a special non-CGB-mode called “PGB mode”.
		if GetBit(romFeatures.ColorGB, 2) || GetBit(romFeatures.ColorGB, 3) {
			romFeatures.PGBMode = true
			romFeatures.SupportColor = false
			romFeatures.SupportMonochrome = true
		} else if romFeatures.ColorGB == 0x80 {
			// The game supports CGB enhancements, but is backwards compatible with monochrome Game Boys
			romFeatures.SupportColor = true
			romFeatures.SupportMonochrome = true
		} else if romFeatures.ColorGB == 0xC0 {
			// The game works on CGB only (the hardware ignores bit 6, so this really functions the same as $80)
			romFeatures.SupportColor = true
			romFeatures.SupportMonochrome = false
		}
	}

	if romFeatures.GBSGBIndicator == 0x03 {
		romFeatures.SupportSGB = true
	} else {
		romFeatures.SupportSGB = false
	}

	return &romFeatures, nil
}
