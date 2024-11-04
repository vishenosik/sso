package operation

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_ReturnFailWithError(t *testing.T) {

	op := "op"
	Err := errors.New("test error")

	String := "string"
	result1, err := ReturnFailWithError(String, op)(Err)
	require.Equal(t, String, result1)
	require.ErrorIs(t, err, Err)

	Bool := false
	result2, _ := ReturnFailWithError(Bool, op)(Err)
	require.Equal(t, Bool, result2)

	Int := 9
	result3, _ := ReturnFailWithError(Int, op)(Err)
	require.Equal(t, Int, result3)

}
