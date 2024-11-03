package main

import (
	"errors"
	"fmt"

	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func main() {
	var (
		valid = validator.New(validator.WithPrivateFieldValidation())
	)

	ID1 := uuid.New().String()
	ID2 := "uuid.New().String()"

	request := models.LoginRequest{
		Email:    "gofakeit.Email()",
		Password: "password",
	}

	if err := errors.Join(
		valid.Var(ID2, "required,uuid4"),
		valid.Struct(request),
	); err != nil {
		fmt.Println(err)
		return
	}

	if err := errors.Join(
		valid.Var(ID1, "required,uuid4"),
		valid.Struct(request),
	); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("success")

}
