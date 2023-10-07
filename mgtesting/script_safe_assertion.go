package mgtesting

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

// ScriptSafeAssert skips the test if not found GoScriptName=TestFunctionName
// in environment variable
func ScriptSafeAssert(t *testing.T) *require.Assertions {
	var callerName = GetCallerFunctionName(2)

	t.Logf("callerName: %v", callerName)

	var goScriptName = os.Getenv("GoScriptName")

	if goScriptName == "" {
		t.Skipf("$GoScriptName is empty, skipping test.")
	}

	if goScriptName != callerName {
		t.Skipf("$GoScriptName != callerName, skipping test. $GoScriptName: %v, callerName: %v",
			goScriptName, callerName)
	}

	return require.New(t)
}

func GetCallerFunctionName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		panic(xerrors.Errorf("runtime.Caller(%v) failure", skip))
	}

	details := runtime.FuncForPC(pc)
	if details == nil {
		panic(xerrors.Errorf("runtime.FuncForPC(%v) failure", pc))
	}

	var name = details.Name()

	var idx = strings.LastIndex(name, ".")
	if idx >= 0 {
		name = name[idx+1:]
	}

	return name
}
