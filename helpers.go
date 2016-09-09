package main

func mult(x byte, y byte) byte {
	var toReturn byte

	var values []byte
	//find most significant set bit
	bits := 0
	finder := y

	//shift right until the largest bit 'falls off'
	for finder > 0 {
		finder = finder >> 1
		bits = bits + 1
	}

	var curMask byte = 0x01
	//now we know we need to calculate the sums of the xtimes up to bits
	for i := 0; i < bits; i++ {
		//for i = 0 it's the 0/1 bit, so we don't run xtime for this
		if i == 0 {
			values = append(values, x)
		} else {
			//so we xtime with the returned value of the xtime for the bit before.
			values = append(values, xtime(values[i-1]))
		}
		//now we check if the bit in question is set. If so, we XOR (add) it to our current result.
		if y&curMask != 0 {
			toReturn = toReturn ^ values[i]
		}
		curMask = curMask << 1
	}

	return toReturn
}

func xtime(x byte) byte {
	if x == 0 {
		return 0
	}

	toReturn := x << 1 //shift left
	//if most significant bit is set, we need to XOR for overflow.
	if x > 127 {
		toReturn = toReturn ^ 0x1b
	}
	return toReturn
}

func getNRFromKeyLen(key []byte) int {
	switch len(key) {
	case 16:
		return 10
	case 24:
		return 12
	case 32:
		return 14
	default:
		return 0
	}
}

func getRoundsFromKeyLen(key []byte) int {
	switch len(key) {
	case 16:
		return 10
	case 24:
		return 9
	case 32:
		return 8
	default:
		return 0
	}
}

func getSubsByte(x byte) byte {
	first := x & 0xf0 >> 4
	second := x & 0x0f

	return sBox[(first*16)+second]
}

func getInvSubsByte(x byte) byte {
	first := x & 0xf0 >> 4
	second := x & 0x0f

	return inverseSBox[(first*16)+second]
}
