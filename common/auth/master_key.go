package auth

import (
	ecdsa "crypto/ecdsa"
	big "math/big"
	elliptic "crypto/elliptic"
)

func MasterKey() *ecdsa.PublicKey {
	key := new(ecdsa.PublicKey)
	key.Curve = elliptic.P256()
	key.X = big.NewInt(0)
	key.X.SetBytes([]byte{0x37, 0x9c, 0x1, 0xf, 0x4d, 0xa1, 0xad, 0x10, 0x81, 0xb3, 0xfb, 0xc7, 0x5a, 0x41, 0xdd, 0xd, 0x40, 0xa6, 0xd6, 0x8e, 0x16, 0x65, 0x66, 0x90, 0xad, 0xe, 0xa2, 0x1c, 0x7d, 0xe3, 0xe, 0xb0})
	key.Y = big.NewInt(0)
	key.Y.SetBytes([]byte{0x60, 0x9f, 0x4f, 0xa7, 0x5c, 0xea, 0x6a, 0x4f, 0x19, 0xcb, 0xbd, 0xd7, 0xf2, 0xba, 0x38, 0xb9, 0x6c, 0xa, 0xd2, 0xb8, 0xe7, 0x9b, 0x84, 0x6, 0x0, 0x33, 0x1d, 0x76, 0xf6, 0xff, 0xb0, 0xb8})
	return key
}
