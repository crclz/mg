package application_test

import (
	"testing"

	"github.com/crclz/mg/internal/application"
	"github.com/stretchr/testify/require"
)

func TestGenerateCommand_Parse_happy_case_1(t *testing.T) {
	// arrange
	var assert = require.New(t)
	var generateCommand = application.NewGenerateCommand(nil)

	// act
	var result, err = generateCommand.Parse("./network_service")
	assert.NoError(err)

	// assert
	assert.Equal("./", result.Dir)
	assert.Equal("./network_service.go", result.GoFileName)
	assert.Equal("./network_service_test.go", result.GoTestFileName)
	assert.Equal("application", result.PackageName)
	assert.Equal("NetworkService", result.ClassName)
}

func TestGenerateCommand_Parse_happy_case_2(t *testing.T) {
	// arrange
	var assert = require.New(t)
	var generateCommand = application.NewGenerateCommand(nil)

	// act
	var result, err = generateCommand.Parse("./some123/network-service")
	assert.NoError(err)

	// assert
	assert.Equal("./some123/", result.Dir)
	assert.Equal("./some123/network_service.go", result.GoFileName)
	assert.Equal("./some123/network_service_test.go", result.GoTestFileName)
	assert.Equal("some123", result.PackageName)
	assert.Equal("NetworkService", result.ClassName)
}
