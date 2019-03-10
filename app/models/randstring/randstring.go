package randstring

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	"1234567890"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// Randstring constructor
func Randstring() {
}

// StringWithCharset return a random string with a specific length and using a specified charset
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// String return a random string with a specific length
func String(length int) string {
	return StringWithCharset(length, charset)
}
