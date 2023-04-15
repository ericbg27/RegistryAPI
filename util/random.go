package util

import (
	"math/rand"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvxwyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(letters)
	for i := 0; i < n; i++ {
		c := letters[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
