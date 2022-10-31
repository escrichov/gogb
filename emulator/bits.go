package emulator

// SetBit the value of the register.
func SetBit(input *byte, pos int, value bool) {
	if value {
		*input |= 1 << pos
	} else {
		*input &= ^(1 << pos)
	}
}

// GetBit the value of the register.
func GetBit(input byte, pos int) bool {
	return (input>>pos)&1 == 1
}

func SetBit16(input *uint16, pos int, value bool) {
	if value {
		*input |= 1 << pos
	} else {
		*input &= ^(1 << pos)
	}
}

// GetBit16 the value of the register.
func GetBit16(input uint16, pos int) bool {
	return (input>>pos)&1 == 1
}

func SetBit8(input uint8, pos uint8, value bool) uint8 {
	result := input
	if value {
		result |= 1 << pos
	} else {
		result &= ^(1 << pos)
	}

	return result
}

func BoolToUint8(value bool) uint8 {
	if value {
		return 1
	} else {
		return 0
	}
}
