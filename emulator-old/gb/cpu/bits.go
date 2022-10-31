package cpu

// HalfCarryAdd half carries two values.
func HalfCarryAdd(val1 byte, val2 byte) bool {
	return (val1&0xF)+(val2&0xF) > 0xF
}
