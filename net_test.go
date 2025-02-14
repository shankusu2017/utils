package utils

import (
	"testing"
)

func TestMac(t *testing.T) {
	b := make([]byte, 6)
	b[0] = 1
	b[1] = 2
	b[2] = 3
	b[3] = 4
	b[4] = 5
	b[5] = 255
	str, _ := MacByte2String(b)
	if str != "01:02:03:04:05:ff" {
		t.Error(str)
	}

	t.Log(str)
}

func TestIPByte2String(t *testing.T) {
	b := make([]byte, 4)
	b[0] = 192
	b[1] = 168
	b[2] = 20
	b[3] = 136
	ipstr, _ := IPByte2String(b)
	t.Log(ipstr)
}
