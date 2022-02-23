package infrastructure

import (
	"errors"
	"math/rand"
	"time"
)

func GetRandNumber(start, end, seed int) (int, error) {
	if end <= start {
		return 0, errors.New("excepted start and ends")
	}
	rand.Seed(int64(seed) + time.Now().UnixNano())
	return rand.Intn(end-start) + start, nil
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
