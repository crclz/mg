package application_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestAAA_asdasd(t *testing.T) {

}

func TestAAA_fail(t *testing.T) {
	t.FailNow()
}

func TestAAA_script_safe(t *testing.T) {
	t.Logf("GoScriptName: %v", os.Getenv("GoScriptName"))
}

func TestAAA_LongRunning(t *testing.T) {
	var assert = require.New(t)
	var outputFile = "long-running.tmp.txt"

	// go run . t -c1 TestAAA_LongRunning
	// 然后ctrl+c，不会打断测试用例，long-running.tmp.txt还是会输出 >5000
	signal.Ignore(syscall.SIGINT)

	var err = os.RemoveAll(outputFile)
	assert.NoError(err)

	var t0 = time.Now()

	var eg = errgroup.Group{}
	eg.SetLimit(2)

	for i := 0; i < 10; i++ {
		eg.Go(func() error {
			time.Sleep(time.Second * 1)
			return nil
		})
	}
	err = eg.Wait()
	assert.NoError(err)

	var testRunningMilliseconds = int64(time.Since(t0).Milliseconds())

	var outputContent = fmt.Sprintf("%v", testRunningMilliseconds)

	err = os.WriteFile(outputFile, []byte(outputContent), 0644)
	assert.NoError(err)
}
