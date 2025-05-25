package core

import (
	"math/rand"
	"time"
)

type ShorteningStrategy interface {
	// shortens the provided str to the desired length by following a
	// shortening strategy, and returns the short string.
	Shorten(str string, length int) string
}

// character pool containing 0-9, A-Z, a-z (62 characters total)
const charPool = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"


type simpleShorteningStrategy struct{}

func NewSimpleShorteningStrategy() ShorteningStrategy {
	return &simpleShorteningStrategy{}
}

// generates a random string of the given length - the random string is in no way
// related to the supplied string. 
// the same str would generate different outputs every single time.
func (s *simpleShorteningStrategy) Shorten(str string, length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	short_str := ""
	for i := 0; i < length; i++ {
		idx := r.Intn(len(charPool))
		short_str = short_str + string(charPool[idx])
	}
	return short_str
}

