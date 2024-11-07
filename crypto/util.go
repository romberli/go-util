package crypto

import (
	"strings"

	"github.com/romberli/go-util/constant"
)

func padKey(key string) string {
	keyLengths := []int{AESKeySize16, AESKeySize24, AESKeySize32}

	if len(key) > AESKeySize32 {
		return key[:AESKeySize32]
	}

	for _, length := range keyLengths {
		if len(key) < length {
			return key + strings.Repeat(constant.EqualString, length-len(key))
		} else if len(key) == length {
			return key
		}
	}

	return key
}
