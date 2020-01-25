package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestTokenSigning(t *testing.T) {
	var token Token

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
		return
	}

	pubKeyMap := make(map[string]*ecdsa.PublicKey)
	privKeyMap := make(map[string]*ecdsa.PrivateKey)

	pubKeyMap["test"] = &key.PublicKey
	privKeyMap["test"] = key

	token.Header.Algorithm = "ES256"
	token.Header.Type = "JWT"
	token.Header.KeyID = "test"

	token.Payload.Issuer = "foo"
	token.Payload.Subject = "bar"
	token.Payload.Audience = "baz"

	raw, err := token.Encode(privKeyMap)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(string(raw))

	decoded, err := ParseToken(raw, pubKeyMap)
	if err != nil {
		t.Error(err)
		return
	}

	if token.Header.Algorithm != decoded.Header.Algorithm {
		t.Error("Algo not match")
	}
}
