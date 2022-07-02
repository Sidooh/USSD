package utils

import (
	"encoding/binary"
	"fmt"
	"math"
)

func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Float64frombytes(bytes []byte) float64 {
	fmt.Println(bytes, string(bytes))
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func CalculateAirtimeEarnings(amount int, subscribed bool) float64 {
	if subscribed {
		return CalculateSubscribedAirtimeEarnings(amount)
	}

	// Get discount
	discount := .06

	// Get users earning ratio
	ratio := .6

	// Get ripples
	ripples := 6

	// Calculate earnings
	return float64(amount) * discount * ratio / float64(ripples)
}

func CalculateSubscribedAirtimeEarnings(amount int) float64 {
	// TODO: Get these values from products api

	// Get discount
	discount := .06

	// Get subscribed users earning ratio
	ratio := 1.0

	// Get ripples (Subscribed users earn 100% passthru at the moment)
	ripples := 1

	// Calculate earnings
	return float64(amount) * discount * ratio / float64(ripples)
}

// TODO: create currency formatter function
// 	picks currency from env, parameters, amount, precision, thousand sep... override prefix..etc
