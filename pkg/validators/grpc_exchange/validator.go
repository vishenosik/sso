package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrNilRequest = errors.New("request is nil")
)

type ProtoExchange interface {
	ProtoMessage()
	ProtoReflect() protoreflect.Message
	Reset()
	String() string
}

func Validate(x ProtoExchange) error {

	if x == nil {
		return ErrNilRequest
	}

	validate := validator.New()

	err := validate.Struct(x)
	if err != nil {
		return err
	}

	return nil
}
