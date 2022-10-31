package rom

type RomHeader struct {
	Title                string
	TitleBytes           []byte
	ColorGB              uint8
	LicenseCodeNew       string
	GBSGBIndicator       uint8
	CartridgeType        uint8
	RomSize              uint8
	RamSize              uint8
	DestinationCode      uint8
	LicenseCodeOld       uint8
	MaskROMVersionNumber uint8
	ComplementCheck      uint8
	CheckSum             uint16
}
