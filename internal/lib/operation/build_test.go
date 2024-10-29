package operation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildOperation(t *testing.T) {

	t.Helper()
	t.Parallel()

	expect := "layer.service.method"
	result := buildOperation("layer", "service", "method")
	assert.Equal(t, expect, result)
}

func Test_ServicesOperation(t *testing.T) {

	t.Helper()
	t.Parallel()

	expect := "Services.service.method"
	result := ServicesOperation("service", "method")
	assert.Equal(t, expect, result)
}
