package business

import (
	"crypto/ecdsa"
)

type Logic struct {
	Config    map[string]string
	MasterKey *ecdsa.PublicKey
}
