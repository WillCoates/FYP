package business

import (
	"context"
	"errors"
	"time"

	"github.com/WillCoates/FYP/auth/model"
	"github.com/WillCoates/FYP/common/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrInvalidUserPassword is returned by Authenticate when it fails to authenticate a user
var ErrInvalidUserPassword error = errors.New("Incorrect email address or password")

// Authenticate attempts to authenticate a user, returning a token if successful
func (logic *Logic) Authenticate(ctx context.Context, email, password, aud string, duration uint32) (*auth.Token, error) {
	usersCollection := logic.db.Collection("users")
	tokensCollection := logic.db.Collection("tokens")

	audience, err := logic.GetAudience(ctx, aud)

	if err != nil {
		return nil, err
	}

	var user model.User
	err = usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInvalidUserPassword
		}
		return nil, err
	}

	if duration == 0 {
		duration = audience.DefaultDuration
	}

	if duration > audience.MaxDuration && audience.MaxDuration != 0 {
		duration = audience.MaxDuration
	}

	if user.CheckPassword(password) != nil {
		return nil, ErrInvalidUserPassword
	}

	token := new(auth.Token)

	token.Header.Type = "JWT"
	token.Header.Algorithm = "ES256"
	token.Header.KeyID = logic.currentKey

	token.Payload.Issuer = "fyp.willtc.uk"
	token.Payload.Subject = user.ID.Hex()
	token.Payload.Audience = aud
	token.Payload.Issued = time.Now().Unix()
	token.Payload.Expires = token.Payload.Issued + int64(duration)
	token.Payload.JwtID = GenerateJwtID()

	registeredToken, err := model.CreateRegisteredToken(token)
	if err != nil {
		return nil, err
	}

	_, err = tokensCollection.InsertOne(ctx, &registeredToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
