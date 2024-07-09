package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"
)

func TestRand(t *testing.T) {
	b := make([]byte, 8)
	rand.Read(b)
	buf := bytes.NewReader(b)

	var i64 uint64
	binary.Read(buf, binary.BigEndian, &i64)
	t.Log(i64)

	t.Log(MakeHexString(6))
}
