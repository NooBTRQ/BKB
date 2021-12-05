package infrastructure

import (
	"math/rand"
	"time"
)

func GetRandNumber(start, end, seed int) int {
	if end <= start {
		panic("get random error")
	}
	rand.Seed(int64(seed) + time.Now().UnixNano())
	return rand.Intn(end-start) + start
}

func Max(num1, num2 int) int {
	if num1 > num2 {
		return num1
	}
	return num2
}

func Min(num1, num2 int) int {
	if num1 > num2 {
		return num2
	}
	return num1
}
