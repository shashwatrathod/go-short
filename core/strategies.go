package core

import (
	"math/rand"
	"time"
)

type AliasingStrategy interface {
	// aliases the provided str to the desired length by following a
	// aliasing strategy, and returns the aliased string.
	Alias(str string, length int) string
}

// character pool containing 0-9, A-Z, a-z (62 characters total)
const charPool = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"


type simplAliasingStrategy struct{}

func NewSimpleAliasingStrategy() AliasingStrategy {
	return &simplAliasingStrategy{}
}

// generates a random string of the given length - the random string is in no way
// related to the supplied string. 
// the same str would generate different outputs every single time.
func (s *simplAliasingStrategy) Alias(str string, length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	alias_str := ""
	for i := 0; i < length; i++ {
		idx := r.Intn(len(charPool))
		alias_str = alias_str + string(charPool[idx])
	}
	return alias_str
}

