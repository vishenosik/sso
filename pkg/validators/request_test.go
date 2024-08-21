package validators

import (
	"testing"

	ssov1 "github.com/blacksmith-vish/sso/protos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

const (
	passDefautlLen = 10
)

func TestValidateRequest(t *testing.T) {

	t.Helper()
	t.Parallel()

	RegReq := new(ssov1.RegisterRequest)

	RegReq.Email = gofakeit.Email()
	RegReq.Password = randomPassword()

	err := ValidateRequest(RegReq)
	require.NoError(t, err)
}

func TestValidateRequest_Err(t *testing.T) {

	t.Helper()
	t.Parallel()

	RegReq := new(ssov1.RegisterRequest)

	RegReq.Email = "gofakeit.Email()"
	RegReq.Password = randomPassword()

	err := ValidateRequest(RegReq)
	require.Error(t, err)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefautlLen)
}
