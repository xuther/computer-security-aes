package main

import "fmt"

func mixColumns(state [4][4]byte) [4][4]byte {
	MultMatrix := [4][4]byte{{2, 3, 1, 1}, {1, 2, 3, 1}, {1, 1, 2, 3}, {3, 1, 1, 2}}
	toReturn := [4][4]byte{}

	//Run through each column in state
	for i := 0; i < 4; i++ {
		curCol := i
		//Matrix Mult
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				toReturn[j][curCol] = toReturn[j][curCol] ^ mult(state[k][curCol], MultMatrix[j][k])
			}
		}
	}
	return toReturn
}

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

func main() {
	state := [4][4]byte{{0xd4, 0, 0, 0}, {0xbf, 0, 0, 0}, {0x5d, 0, 0, 0}, {0x30, 0, 0, 0}}
	fmt.Printf("%v", mixColumns(state))
	fmt.Printf("%v", mult(0x57, 0x13))
}
