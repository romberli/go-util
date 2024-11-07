package crypto

func padKey(key []byte) []byte {
	keyLengths := []int{AESKeySize16, AESKeySize24, AESKeySize32}

	if len(key) > AESKeySize32 {
		return key[:AESKeySize32]
	}

	for _, length := range keyLengths {
		if len(key) < length {
			return append(key, make([]byte, length-len(key))...)
		} else if len(key) == length {
			return key
		}
	}

	return key
}
