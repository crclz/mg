package mgtesting_test

import (
	"testing"

	"github.com/crclz/mg/mgtesting"
)

func TestScriptSafeAssert_case_1(t *testing.T) {
	var assert = mgtesting.ScriptSafeAssert(t)

	assert.False(true)
}

func TestGetCallerFunctionName_happy_1(t *testing.T) {
	var result = mgtesting.GetCallerFunctionName(1)

	t.Logf("result: %v", result)
}
