package application_test

import (
	"os"
	"testing"
)

func TestAAA_asdasd(t *testing.T) {

}

func TestAAA_fail(t *testing.T) {
	t.FailNow()
}

func TestAAA_script_safe(t *testing.T) {
	t.Logf("GoScriptName: %v", os.Getenv("GoScriptName"))
}
