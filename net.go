package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
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

func IPByte2String(b []byte) (string, error) {
	if len(b) == 4 {
		ip := net.IPv4(b[0], b[1], b[2], b[3])
		return ip.String(), nil
	} else if len(b) == 16 {
		ip := net.IP(b)
		return ip.String(), nil
	} else {
		return "", errors.New(fmt.Sprintf("invalid length: %d", len(b)))
	}
}
