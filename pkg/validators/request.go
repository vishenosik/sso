package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrNilRequest = errors.New("request is nil")
)

type Request interface {
	ProtoReflect() protoreflect.Message
	ProtoMessage()
	Descriptor() ([]byte, []int)
}

func ValidateRequest(x Request) error {

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
