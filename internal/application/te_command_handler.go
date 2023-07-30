package application

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/subcommands"
)

type TeCommandHandler struct {
}

func (*TeCommandHandler) Name() string     { return "t" }
func (*TeCommandHandler) Synopsis() string { return "Run tests with convenience" }
func (*TeCommandHandler) Usage() string {
	return "refer to readme.md"
}

func (p *TeCommandHandler) SetFlags(f *flag.FlagSet) {
}

func (p *TeCommandHandler) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var positionalArgs = f.Args()
	if len(positionalArgs) != 1 {
		fmt.Printf("Expecting 1 positional argument, but got %v.\n", len(positionalArgs))
		return subcommands.ExitFailure
	}

	var testName = strings.TrimSpace(positionalArgs[0])

	if testName == "" {
		fmt.Printf("testName is empty\n")
	}

	// find test according to testName

	var goTestSourceFiles []string

	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_test.go") {
			if !f.IsDir() {
				goTestSourceFiles = append(goTestSourceFiles, path)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("goTestSourceFiles: %v", goTestSourceFiles)

	return subcommands.ExitSuccess
}
