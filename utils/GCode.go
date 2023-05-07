package utils

import (
	"math/rand"
)

// RandInt 生成随机数
func RandInt(min int, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min+1) + min
}
