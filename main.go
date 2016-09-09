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

	//Finish the rest of the key
	for i := 4; i < len(nextKey); i++ {
		nextKey[i] = prevKey[i] ^ nextKey[i-4]
		//extra subs
		if len(nextKey) == 32 && i >= 16 && i < 20 {
			nextKey[i] = getSubsByte(nextKey[i])
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

func decrypt(input []byte, key []byte) []byte {
	return []byte{}
}

func main() {
	input := []byte{
		0x00, 0x11, 0x22, 0x33,
		0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb,
		0xcc, 0xdd, 0xee, 0xff}
	//key := []byte{
	//	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	//	0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	//	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
	//	0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	key := []byte{
		0x60, 0x3d, 0xeb, 0x10, 0x15, 0xca, 0x71, 0xbe,
		0x2b, 0x73, 0xae, 0xf0, 0x85, 0x7d, 0x77, 0x81,
		0x1f, 0x35, 0x2c, 0x07, 0x3b, 0x61, 0x08, 0xd7,
		0x2d, 0x98, 0x10, 0xa3, 0x09, 0x14, 0xdf, 0xf4}
	output := []byte{
		0x8e, 0xa2, 0xb7, 0xca,
		0x51, 0x67, 0x45, 0xbf,
		0xea, 0xfc, 0x49, 0x90,
		0x4b, 0x49, 0x60, 0x89}

	fmt.Printf("%x\n%x", encryptWPrint(input, key), output)
}
