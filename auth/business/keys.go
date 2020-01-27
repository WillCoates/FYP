package business

import "crypto/ecdsa"

func (logic *Logic) SetKeyConfig(currentKey string, keys *map[string]*ecdsa.PrivateKey) {
	logic.keyConfigLock.Lock()
	defer logic.keyConfigLock.Unlock()

	logic.currentKey = currentKey
	logic.keys = *keys

	logic.pubKeys = make(map[string]*ecdsa.PublicKey)
	for name, key := range logic.keys {
		logic.pubKeys[name] = &key.PublicKey
	}
}
