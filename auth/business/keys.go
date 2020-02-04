package business

import (
	"crypto/ecdsa"

	"github.com/WillCoates/FYP/common/auth"
)

func (logic *Logic) SetKeyConfig(currentKey string, keys map[string]*ecdsa.PrivateKey) {
	logic.keyConfigLock.Lock()
	defer logic.keyConfigLock.Unlock()

	logic.currentKey = currentKey
	logic.keys = keys

	logic.pubKeys = auth.NewKeyBundle()
	for name, key := range logic.keys {
		logic.pubKeys.Keys[name] = &key.PublicKey
	}
}
