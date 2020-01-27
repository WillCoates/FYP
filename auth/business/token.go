package business

import "github.com/WillCoates/FYP/common/auth"

// EncodeToken serializes a token into JWT, applying a signature
func (logic *Logic) EncodeToken(token *auth.Token) ([]byte, error) {
	logic.keyConfigLock.RLock()
	defer logic.keyConfigLock.RUnlock()
	return token.Encode(logic.keys)
}

// DecodeToken deseralizes a token from JWT, returning an error if the format
// is invalid or has a bad signature
func (logic *Logic) DecodeToken(token []byte) (*auth.Token, error) {
	logic.keyConfigLock.RLock()
	defer logic.keyConfigLock.RUnlock()
	return auth.ParseToken(token, logic.pubKeys)
}

// EncodeTokenStr serializes a token into JWT, applying a signature
func (logic *Logic) EncodeTokenStr(token *auth.Token) (string, error) {
	logic.keyConfigLock.RLock()
	defer logic.keyConfigLock.RUnlock()
	res, err := token.Encode(logic.keys)
	return string(res), err
}

// DecodeTokenStr deseralizes a token from JWT, returning an error if the
// format is invalid or has a bad signature
func (logic *Logic) DecodeTokenStr(token string) (*auth.Token, error) {
	logic.keyConfigLock.RLock()
	defer logic.keyConfigLock.RUnlock()
	return auth.ParseToken([]byte(token), logic.pubKeys)
}
