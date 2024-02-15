package util

import (
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateID(n uint) string {
	res := make([]byte, n)
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}

	return string(res)
}
