package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const COST int = 12

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	EmailAddress string             `bson:"email"`
	Password     []byte             `bson:"password"`
	Name         string             `bson:"name"`
}

func (user *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), COST)
	if err != nil {
		return err
	}
	user.Password = hash
	return nil
}

func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
