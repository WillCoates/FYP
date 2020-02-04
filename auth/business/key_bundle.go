package business

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"

	"github.com/WillCoates/FYP/common/auth"
)

var ErrNoMasterKey = errors.New("No master key present")

func (logic *Logic) GetKeyBundle() ([]byte, error) {
	logic.keyConfigLock.RLock()
	defer logic.keyConfigLock.RUnlock()
	return logic.pubKeys.Encode()
}

func (logic *Logic) GetKeyBundleSignature() ([]byte, error) {
	logic.keyConfigLock.RLock()
	defer logic.keyConfigLock.RUnlock()
	bundle, err := logic.pubKeys.Encode()

	if err != nil {
		return nil, err
	}

	master, ok := logic.keys["master"]
	if !ok {
		return nil, ErrNoMasterKey
	}

	hash := sha256.Sum256(bundle)

	r, s, err := ecdsa.Sign(rand.Reader, master, hash[:])

	signature := auth.EcdsaSignature{R: r, S: s}

	return json.Marshal(&signature)
}
