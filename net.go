package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
)

func MacByte2String(b []byte) (string, error) {
	if len(b) != 6 {
		return "", errors.New(fmt.Sprintf("invalid length: %d", len(b)))
	}
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s",
		hex.EncodeToString(b[0:1]),
		hex.EncodeToString(b[1:2]),
		hex.EncodeToString(b[2:3]),
		hex.EncodeToString(b[3:4]),
		hex.EncodeToString(b[4:5]),
		hex.EncodeToString(b[5:6])), nil
}
