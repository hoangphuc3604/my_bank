package util

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandOwner() string {
	n := rand.Intn(6) + 5 
	name := make([]byte, n)
	for i := range name {
		name[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(name)
}

func RandBalance(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandCurrency() string {
	currencies := []string{USD, EUR, GBP, JPY, AUD}
	return currencies[rand.Intn(len(currencies))]
}

func RandEmail() string {
	return RandOwner() + "@gmail.com"
}
