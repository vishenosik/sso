package validator

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	passDefautlLen = 10
)

type TestingRequest struct {
	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty" validate:"required,email"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty" validate:"required"`
	UserId   int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" validate:"gte=0"`
}

func (x *TestingRequest) Reset() {}

func (x *TestingRequest) String() string { return "" }

func (*TestingRequest) ProtoMessage() {}

func (x *TestingRequest) ProtoReflect() protoreflect.Message { return nil }

func TestValidateRequest(t *testing.T) {

	t.Helper()
	t.Parallel()

	RegReq := new(TestingRequest)

	RegReq.Email = gofakeit.Email()
	RegReq.Password = randomPassword()

	err := Validate(RegReq)
	require.NoError(t, err)
}

func TestValidateRequest_Err(t *testing.T) {

	t.Helper()
	t.Parallel()

	RegReq := new(TestingRequest)

	RegReq.Email = "gofakeit.Email()"
	RegReq.Password = randomPassword()

	err := Validate(RegReq)
	require.Error(t, err)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefautlLen)
}
