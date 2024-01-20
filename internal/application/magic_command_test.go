package application_test

import (
	"context"
	"strings"
	"testing"

	"github.com/crclz/mg/internal/application"
	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/stretchr/testify/require"
)

// go test -v ./internal/application --run TestMagicCommand_ApplyMagicTestClassMethod_happyCase1
func TestMagicCommand_ApplyMagicTestClassMethod_happyCase1(t *testing.T) {
	// arrange
	var assert = require.New(t)
	var ctx = context.Background()
	var magicCommand = application.NewMagicCommand(
		domainservices.GetSingletonMgContextService(),
		domainservices.GetSingletonFileDiscoveryService(),
	)

	var sourceCode = `//!magic test NetworkService.Detect return false  when not connected`
	var sourceCodeLines = strings.Split(sourceCode, "\n")

	// act
	result, err := magicCommand.ApplyMagicTestClassMethod(ctx, "some_package", sourceCodeLines)
	assert.NoError(err)

	assert.Contains(strings.Join(result, "\n"), "TestNetworkService_Detect_returnFalseWhenNotConnected")

	// assert
	t.Logf("Result:")
	for _, line := range result {
		t.Log(line)
	}
}
