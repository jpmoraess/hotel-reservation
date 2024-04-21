package types

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
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

type UpdateUserInput struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=40"`
	LastName  string `json:"lastName" validate:"required,min=2,max=40"`
}

func (p *UpdateUserInput) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}
	return m
}

func (input CreateUserInput) Validate(ctx context.Context) map[string]string {
	validate := validator.New()
	if err := validate.StructCtx(ctx, input); err != nil {
		errors := make(map[string]string)
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			field := strings.ToLower(validationError.StructField())
			switch validationError.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "min":
				errors[field] = fmt.Sprintf("%s is required with min %s", field, validationError.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s is required with max %s", field, validationError.Param())
			case "email":
				errors[field] = fmt.Sprintf("%s is invalid email", field)
			}
		}
		return errors
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
