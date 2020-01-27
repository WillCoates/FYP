package business

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateJwtID generates a random string which should be unpredictable and almost guarenteed to be unique
func GenerateJwtID() string {
	var bytes [18]byte
	generated := 0

	for generated < 18 {
		n, _ := rand.Reader.Read(bytes[generated:])
		generated += n
	}

	return base64.RawURLEncoding.EncodeToString(bytes[:])
}
