package main

import "fmt"

func shiftRows(message []byte) []byte {

	//three rows
	for i := 1; i < 4; i++ {
		//four values to move
		shiftRow := i % 4
		for j := 0; j < i; j++ {
			temp := message[shiftRow]
			for k := 0; k < 3; k++ {
				message[shiftRow+(4*k)] = message[shiftRow+(4*(k+1))]
			}
			message[shiftRow+(4*3)] = temp
		}
	}
	return message
}

func mixColumnsWrapper(message []byte) []byte {

	stateForMixCol := [4][4]byte{}
	for i := 0; i < len(message); i++ {
		stateForMixCol[i%4][i/4] = message[i]
	}

	cols := mixColumns(stateForMixCol)
	for i := 0; i < len(message); i++ {
		message[i] = cols[i%4][i/4]
	}
	return message
}

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

func addKey(key []byte, message []byte) []byte {
	for i := 0; i < len(message); i++ {
		message[i] = message[i] ^ key[i]
	}
	return message
}

func subBytes(message []byte) []byte {
	for i := range message {
		message[i] = getSubsByte(message[i])
	}
	return message
}

func getKeySchedule(key []byte) []byte {
	rounds := getRoundsFromKeyLen(key)

	finalKey := key

	for i := 0; i < rounds; i++ {
		key = getNextKeySchedule(key, i)
		finalKey = append(finalKey, key...)
	}

	return finalKey
}

func getNextKeySchedule(prevKey []byte, round int) []byte {
	nextKey := make([]byte, len(prevKey))

	//Rotate and substitute last word.
	newLast := prevKey[len(prevKey)-4]
	for i := 1; i < 4; i++ {

		nextKey[i-1] = prevKey[len(prevKey)-(4-i)]
		nextKey[i-1] = getSubsByte(nextKey[i-1])
	}

	nextKey[3] = getSubsByte(newLast)
	//xor with recon
	nextKey[0] = nextKey[0] ^ reconBytes[round]

	//xor with first element
	for i := 0; i < 4; i++ {
		nextKey[i] = prevKey[i] ^ nextKey[i]
	}

	//fmt.Printf("%x\n", nextKey[:4])

	//Finish the rest of the key
	for i := 4; i < len(nextKey); i++ {
		nextKey[i] = prevKey[i] ^ nextKey[i-4]
		//DEBUG
		if i%4 == 0 {
			//fmt.Printf("%v", i)
			//fmt.Printf("%x\n", nextKey[i-4:i])
		}
		///DEBUG
		//extra subs
		if len(nextKey) == 32 && i >= 16 && i < 20 {
			nextKey[i] = nextKey[i-4]
			nextKey[i] = getSubsByte(nextKey[i])
			nextKey[i] = prevKey[i] ^ nextKey[i]
		}
	}

	return nextKey
}

func encryptWPrint(input []byte, key []byte) []byte {

	keySchedule := getKeySchedule(key)
	//init add round key
	fmt.Printf("start[0]\n%x\n", input)
	input = addKey(keySchedule[:16], input)
	fmt.Printf("k_sch[0]\n%x\n", keySchedule[:16])

	nr := getNRFromKeyLen(key)

	//begin the run.
	for i := 1; i < nr; i++ {
		fmt.Printf("start[%v]\n%x\n", i, input)
		input = subBytes(input)
		fmt.Printf("s_box[%v]\n%x\n", i, input)
		input = shiftRows(input)
		fmt.Printf("s_row[%v]\n%x\n", i, input)
		input = mixColumnsWrapper(input)
		fmt.Printf("m_col[%v]\n%x\n", i, input)
		input = addKey(keySchedule[i*16:(i+1)*16], input)
		fmt.Printf("k_sch[%v]\n%x\n", i, keySchedule[i*16:(i+1)*16])
	}
	fmt.Printf("start[%v]\n%x\n", nr, input)
	input = subBytes(input)
	fmt.Printf("s_box[%v]\n%x\n", nr, input)
	input = shiftRows(input)
	fmt.Printf("s_row[%v]\n%x\n", nr, input)
	input = addKey(keySchedule[nr*16:(nr+1)*16], input)
	fmt.Printf("k_sch[%v]\n%x\n", nr, keySchedule[nr*16:(nr+1)*16])

	fmt.Printf("output[%v]\n%x\n", nr, input)

	return input
}

func inverseShiftRows(input []byte) []byte {
	//three rows
	for i := 1; i < 4; i++ {
		//four values to move
		shiftRow := i % 4
		for j := 0; j < i; j++ {
			temp := input[shiftRow+(4*3)]
			for k := 2; k >= 0; k-- {
				input[shiftRow+(4*(k+1))] = input[shiftRow+(4*k)]
			}

			input[shiftRow] = temp
			//Print new matrix
		}

		//fmt.Printf("%x,%x,%x,%x\n", input[shiftRow], input[shiftRow+4], input[shiftRow+8], input[shiftRow+12])
	}
	return input
}

func inverseSubBytes(input []byte) []byte {
	for i := range input {
		input[i] = getInvSubsByte(input[i])
	}
	return input
}

func inverseMixColsWrapper(input []byte) []byte {
	stateForMixCol := [4][4]byte{}
	for i := 0; i < len(input); i++ {
		stateForMixCol[i%4][i/4] = input[i]
	}

	cols := inverseMixCols(stateForMixCol)
	for i := 0; i < len(input); i++ {
		input[i] = cols[i%4][i/4]
	}
	return input
}

func inverseMixCols(state [4][4]byte) [4][4]byte {
	MultMatrix := [4][4]byte{
		{0x0e, 0x0b, 0x0d, 0x09},
		{0x09, 0x0e, 0x0b, 0x0d},
		{0x0d, 0x09, 0x0e, 0x0b},
		{0x0b, 0x0d, 0x09, 0x0e}}
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

func encrypt(input []byte, key []byte) []byte {

	keySchedule := getKeySchedule(key)
	//init add round key
	input = addKey(keySchedule[:16], input)

	nr := getNRFromKeyLen(key)

	//begin the run.
	for i := 1; i < nr; i++ {
		input = subBytes(input)
		input = shiftRows(input)
		input = mixColumnsWrapper(input)
		input = addKey(keySchedule[i*16:(i+1)*16], input)
	}

	input = subBytes(input)
	input = shiftRows(input)
	input = addKey(keySchedule[nr*16:(nr+1)*16], input)

	return input
}

func decryptEquivWPrint(input []byte, key []byte) []byte {
	keySchedule := getKeySchedule(key)
	nr := getNRFromKeyLen(key)

	//init add round key
	fmt.Printf("start[0]\n%x\n", input)
	key = keySchedule[nr*16 : (nr+1)*16]
	fmt.Printf("rawKey\n%x\n", key)
	//key = inverseMixColsWrapper(key)

	input = addKey(key, input)
	fmt.Printf("k_sch[0]\n%x\n", key)

	//begin the run.
	for i := 1; i < nr; i++ {
		fmt.Printf("start[%v]\n%x\n", i, input)
		input = inverseSubBytes(input)
		fmt.Printf("s_box[%v]\n%x\n", i, input)
		input = inverseShiftRows(input)
		fmt.Printf("s_row[%v]\n%x\n", i, input)
		input = inverseMixColsWrapper(input)
		fmt.Printf("m_col[%v]\n%x\n", i, input)
		key = keySchedule[(nr-i)*16 : ((nr-i)+1)*16]
		key = inverseMixColsWrapper(key)
		input = addKey(key, input)
		fmt.Printf("k_sch[%v]\n%x\n", i, key)
	}
	fmt.Printf("start[%v]\n%x\n", nr, input)
	input = inverseSubBytes(input)
	fmt.Printf("s_box[%v]\n%x\n", nr, input)
	input = inverseShiftRows(input)
	fmt.Printf("s_row[%v]\n%x\n", nr, input)
	key = keySchedule[:16]
	//key = inverseMixColsWrapper(key)
	input = addKey(key, input)
	fmt.Printf("k_sch[%v]\n%x\n", nr, keySchedule[:16])

	fmt.Printf("output[%v]\n%x\n", nr, input)

	return input
}

func decryptWPrint(input []byte, key []byte) []byte {
	keySchedule := getKeySchedule(key)
	nr := getNRFromKeyLen(key)

	//init add round key
	fmt.Printf("start[0]\n%x\n", input)
	input = addKey(keySchedule[nr*16:(nr+1)*16], input)
	fmt.Printf("k_sch[0]\n%x\n", keySchedule[nr*16:(nr+1)*16])

	//begin the run.
	for i := 1; i < nr; i++ {
		fmt.Printf("start[%v]\n%x\n", i, input)
		input = inverseShiftRows(input)
		fmt.Printf("s_row[%v]\n%x\n", i, input)
		input = inverseSubBytes(input)
		fmt.Printf("s_box[%v]\n%x\n", i, input)
		input = addKey(keySchedule[(nr-i)*16:((nr-i)+1)*16], input)
		fmt.Printf("k_sch[%v]\n%x\n", i, keySchedule[(nr-i)*16:((nr-i)+1)*16])
		fmt.Printf("k_add[%v]\n%x\n", i, input)
		input = inverseMixColsWrapper(input)
	}
	fmt.Printf("start[%v]\n%x\n", nr, input)
	input = inverseShiftRows(input)
	fmt.Printf("s_row[%v]\n%x\n", nr, input)
	input = inverseSubBytes(input)
	fmt.Printf("s_box[%v]\n%x\n", nr, input)
	input = addKey(keySchedule[:16], input)
	fmt.Printf("k_sch[%v]\n%x\n", nr, keySchedule[:16])

	fmt.Printf("output[%v]\n%x\n", nr, input)

	return input
}

func main() {

}
