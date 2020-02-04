package business

import (
	"crypto/ecdsa"
	"sync"

	"github.com/WillCoates/FYP/common/auth"
	"go.mongodb.org/mongo-driver/mongo"
)

// Logic is a structure which provides all business logic for authenticating
// and managing users
type Logic struct {
	db *mongo.Database

	keyConfigLock sync.RWMutex
	currentKey    string
	keys          map[string]*ecdsa.PrivateKey
	// pubKeys       map[string]*ecdsa.PublicKey
	pubKeys *auth.KeyBundle
}

// Constructs a new logic object
func MakeLogic(db *mongo.Database) *Logic {
	logic := new(Logic)
	logic.db = db
	return logic
}
