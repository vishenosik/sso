package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RequestContext(t *testing.T) {
	ctx := context.Background()

	requestID := "12345"

	requestCtx, err := NewRequestContext(requestID)
	require.NoError(t, err)

	ctx = With(ctx, requestCtx)

	actual, ok := RequestContext(ctx)
	require.True(t, ok)

	app, err := AppContext(ctx)
	require.Error(t, err)
	require.Nil(t, app)

	require.Equal(t, requestID, actual.RequestID)
}
