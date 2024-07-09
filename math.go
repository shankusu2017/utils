package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// RandWeigh 带权重的随机
func RandWeigh(r []int) (int, error) {
	ttl := 0
	for _, weight := range r {
		if weight < 0 {
			return -1, errors.New(fmt.Sprintf("963c1597 invalid weight:%d", weight))
		}
		ttl += weight
	}

	if ttl < 1 {
		return -1, errors.New(fmt.Sprintf("e1e2b1de invalid slice ttl:%d", ttl))
	}

	b := make([]byte, 8)
	rand.Read(b)
	buf := bytes.NewReader(b)

	var i64 uint64
	binary.Read(buf, binary.BigEndian, &i64)

	i64 %= uint64(ttl)
	i64 += 1 // weight ==0 的项则被忽略

	for i, weight := range r {
		if i64 <= uint64(weight) {
			return i, nil
		} else {
			i64 -= uint64(weight)
		}
	}

	return -1, errors.New("7457d252 library error")
}

// MakeMagic2Buf 用随机值填充buf
func MakeMagic2Buf(buf []byte) {
	rand.Read(buf)
}

// MakeHexString buf.len = 4,那么产生随机的4个Byte，转换成8个CHAR "6eb88a73"
func MakeHexString(len int) string {
	buf := make([]byte, len)
	rand.Read(buf)
	str := hex.EncodeToString(buf)
	return str
}
