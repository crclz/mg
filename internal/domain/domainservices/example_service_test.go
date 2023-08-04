package domainservices_test

import (
	"context"
	"testing"

	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/stretchr/testify/require"
)

func TestExampleService_HelloWorlds_whenProvideEmptyWorldsThenReturnOnlyHello(t *testing.T) {
	var assert = require.New(t) // TODO: replace with env check

	// arrange
	var ctx = context.TODO() // TODO: replace with wire
	var exampleService = domainservices.GetSingletonExampleService()

	// act
	var result, err = exampleService.HelloWorlds(ctx, "[]")
	assert.NoError(err)

	// assert
	assert.Equal("Hello!", result)
}
