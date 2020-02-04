package auth

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
)

var ErrBadSignature = errors.New("Bad signature")

type KeyBundle struct {
	Keys map[string]*ecdsa.PublicKey
}

type EcdsaSignature struct {
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

func NewKeyBundle() *KeyBundle {
	bundle := new(KeyBundle)
	bundle.Keys = make(map[string]*ecdsa.PublicKey)
	return bundle
}

func LoadBundleHTTP(endpoint string, master *ecdsa.PublicKey) (*KeyBundle, error) {
	sigChan := make(chan *EcdsaSignature)
	bunChan := make(chan *KeyBundle)
	sigError := make(chan error)
	bunError := make(chan error)

	// Download bundle
	go func() {
		checkErr := func(err error) bool {
			if err != nil {
				// Ensure signature routine finishes
				select {
				case <-sigError:
				case <-sigChan:
				}
				bunError <- err
				return true
			}
			return false
		}

		resp, err := http.Get(endpoint)
		if checkErr(err) {
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if checkErr(err) {
			return
		}

		keys := make(map[string]string)
		err = json.Unmarshal(body, &keys)
		if checkErr(err) {
			return
		}

		bundle := NewKeyBundle()

		for name, keyBase64 := range keys {
			decoded, err := base64.RawURLEncoding.DecodeString(keyBase64)
			if checkErr(err) {
				return
			}
			key, err := x509.ParsePKIXPublicKey(decoded)
			if checkErr(err) {
				return
			}
			bundle.Keys[name] = key.(*ecdsa.PublicKey)
		}

		select {
		case err := <-sigError:
			bunError <- err

		case sig := <-sigChan:
			hash := sha256.Sum256(body)
			if ecdsa.Verify(master, hash[:], sig.R, sig.S) {
				bunChan <- bundle
			} else {
				bunError <- ErrBadSignature
			}
		}
	}()

	// Download signature
	go func() {
		resp, err := http.Get(endpoint + ".sig")
		if err != nil {
			sigError <- err
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			sigError <- err
			return
		}

		sig := new(EcdsaSignature)

		err = json.Unmarshal(body, sig)
		if err != nil {
			sigError <- err
			return
		}

		sigChan <- sig
	}()

	select {
	case err := <-bunError:
		return nil, err

	case bun := <-bunChan:
		return bun, nil
	}
}

func (bundle *KeyBundle) Encode() ([]byte, error) {
	keys := make(map[string]string)
	for name, key := range bundle.Keys {
		encoded, err := x509.MarshalPKIXPublicKey(key)
		if err != nil {
			return nil, err
		}

		keys[name] = base64.RawURLEncoding.EncodeToString(encoded)
	}

	return json.Marshal(&keys)
}
