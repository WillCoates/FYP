package auth

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
)

type Token struct {
	Header  TokenHeader
	Payload TokenPayload
}

var tokenSeperator = []byte{'.'}

func ParseToken(raw []byte, keys map[string]*ecdsa.PublicKey) (*Token, error) {
	token := new(Token)
	var headerRaw, payloadRaw, signatureRaw []byte
	var header, payload, signature []byte

	headerEnd := bytes.Index(raw, tokenSeperator)
	if headerEnd == -1 {
		return nil, errors.New("Missing seperator (header)")
	}

	payloadEnd := bytes.Index(raw[headerEnd+1:], tokenSeperator)

	if payloadEnd == -1 {
		return nil, errors.New("Missing seperator (payload)")
	}

	payloadEnd += headerEnd + 1

	headerRaw = raw[:headerEnd]
	payloadRaw = raw[headerEnd+1 : payloadEnd]
	signatureRaw = raw[payloadEnd+1:]

	header = make([]byte, base64.RawURLEncoding.DecodedLen(len(headerRaw)))
	payload = make([]byte, base64.RawURLEncoding.DecodedLen(len(payloadRaw)))
	signature = make([]byte, base64.RawURLEncoding.DecodedLen(len(signatureRaw)))

	_, err := base64.RawURLEncoding.Decode(header, headerRaw)
	if err != nil {
		return nil, err
	}

	_, err = base64.RawURLEncoding.Decode(payload, payloadRaw)
	if err != nil {
		return nil, err
	}

	_, err = base64.RawURLEncoding.Decode(signature, signatureRaw)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(header, &token.Header)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(payload, &token.Payload)
	if err != nil {
		return nil, err
	}

	key, keyFound := keys[token.Header.KeyID]

	if !keyFound {
		return nil, errors.New("Unknown key ID")
	}

	var hash []byte
	var r big.Int
	var s big.Int

	switch token.Header.Algorithm {
	case "ES256":
		temp := sha256.Sum256(raw[:payloadEnd])
		hash = make([]byte, len(temp))
		copy(hash, temp[:])

		r.SetBytes(signature[:32])
		s.SetBytes(signature[32:])

	case "ES512":
		temp := sha512.Sum512(raw[:payloadEnd])
		hash = make([]byte, len(temp))
		copy(hash, temp[:])

		r.SetBytes(signature[:64])
		s.SetBytes(signature[64:])

	default:
		return nil, errors.New("Unsupported algorithm")
	}

	if !ecdsa.Verify(key, hash, &r, &s) {
		return nil, errors.New("Failed to verify signature")
	}

	return token, nil
}

func (token *Token) Encode(keys map[string]*ecdsa.PrivateKey) ([]byte, error) {
	header, err := json.Marshal(&token.Header)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(&token.Payload)
	if err != nil {
		return nil, err
	}

	headerLength := base64.RawURLEncoding.EncodedLen(len(header))
	payloadLength := base64.RawURLEncoding.EncodedLen(len(payload))

	var raw []byte

	switch token.Header.Algorithm {
	case "ES256":
		raw = make([]byte, headerLength+payloadLength+2+base64.RawURLEncoding.EncodedLen(64))

	case "ES512":
		raw = make([]byte, headerLength+payloadLength+2+base64.RawURLEncoding.EncodedLen(128))
	}

	var headerEnd = headerLength
	var payloadEnd = headerLength + payloadLength + 1

	raw[headerEnd] = '.'
	raw[payloadEnd] = '.'

	base64.RawURLEncoding.Encode(raw[:headerEnd], header)
	base64.RawURLEncoding.Encode(raw[headerEnd+1:payloadEnd], payload)

	key, found := keys[token.Header.KeyID]

	if !found {
		return nil, errors.New("Unknown key ID")
	}

	var hash []byte

	switch token.Header.Algorithm {
	case "ES256":
		temp := sha256.Sum256(raw[:payloadEnd])
		hash = make([]byte, len(temp))
		copy(hash, temp[:])
	case "ES512":
		temp := sha512.Sum512(raw[:payloadEnd])
		hash = make([]byte, len(temp))
		copy(hash, temp[:])
	default:
		return nil, errors.New("Unsupported algorithm")
	}

	r, s, err := ecdsa.Sign(rand.Reader, key, hash)
	if err != nil {
		return nil, err
	}

	rBytes := r.Bytes()
	sBytes := s.Bytes()

	signature := make([]byte, len(rBytes)+len(sBytes))
	copy(signature[:len(rBytes)], rBytes)
	copy(signature[len(rBytes):], sBytes)

	base64.RawURLEncoding.Encode(raw[payloadEnd+1:], signature)

	return raw, nil
}

// TokenHeader JSON Web Token header
type TokenHeader struct {
	// Aglorithm used to sign the token
	Algorithm string `json:"alg"`

	// Type of the payload, should be "JWT"
	Type string `json:"typ"`

	// Key identifier
	KeyID string `json:"kid"`
}

// TokenPayload JSON Web Token payload
type TokenPayload struct {
	// Who issued the JWT
	Issuer string `json:"iss,omitempty"`

	// Who is the JWT describing
	Subject string `json:"sub,omitempty"`

	// Who is the intended user of the JWT
	Audience string `json:"aud,omitempty"`

	// When the token expires
	Expires int64 `json:"exp"`

	// When the token becomes valid
	NotBefore int64 `json:"nbf,omitempty"`

	// When the token was issued
	Issued int64 `json:"iat,omitempty"`

	// A unique ID for the token
	JwtID string `json:"jti,omitempty"`
}
