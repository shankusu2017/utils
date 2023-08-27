package utils

import "math/rand"

func SliceRandOne(sli []string) string {
	l := len(sli)
	if l == 0 {
		return ""
	}

	ret := rand.Uint32()
	return sli[(int)(ret)%l]
}
