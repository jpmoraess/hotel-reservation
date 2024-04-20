package types

import (
	"context"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

type CreateUserInput struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=40"`
	LastName  string `json:"lastName" validate:"required,min=2,max=40"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (input CreateUserInput) Validate(ctx context.Context) error {
	validate := validator.New()
	if err := validate.StructCtx(ctx, input); err != nil {
		return err
	}
	return nil
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromInput(input CreateUserInput) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Email:             input.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
