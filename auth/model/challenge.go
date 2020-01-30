package model

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"github.com/WillCoates/FYP/common/util"
)

type Challenge struct {
	Challenge []byte `bson:"chal"`
	Algorithm string `bson:"algo"`
}

func NewChallenge(challenge, algorithm string) (*Challenge, error) {
	rawChallenge, err := base64.RawStdEncoding.DecodeString(challenge)
	if err != nil {
		return nil, err
	}

	result := new(Challenge)
	result.Challenge = rawChallenge
	result.Algorithm = algorithm

	return result, nil
}

var ErrBadAlgorithm error = errors.New("Bad algorithm")
var ErrWrongSolution error = errors.New("Wrong solution for challenge")

func (challenge *Challenge) Prove(solution string) error {
	var result []byte
	solutionRaw, err := base64.RawURLEncoding.DecodeString(solution)
	if err != nil {
		return err
	}

	switch challenge.Algorithm {
	case "S256":
		result = make([]byte, 32)
		temp := sha256.Sum256(solutionRaw)
		copy(result, temp[:])

	case "S512":
		result = make([]byte, 64)
		temp := sha512.Sum512(solutionRaw)
		copy(result, temp[:])

	default:
		return ErrBadAlgorithm
	}

	if !util.SecureEquals(result, challenge.Challenge) {
		return ErrWrongSolution
	}

	return nil
}
