package utils

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

// GetBitInt the value of the register.
func GetBitInt(input byte, pos int) int {
    return int((input >> pos) & 1)
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
