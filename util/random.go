package util

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder

	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetLength := len(alphabet)

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(alphabetLength)])
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(10)
}

func RandomBalance() int64 {
	return RandomInt(500, 2000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CHF"}
	idx := len(currencies)
	return currencies[rand.Intn(idx)]
}
