package application

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
		return subcommands.ExitFailure
	}

	if !strings.HasPrefix(testName, "Test") {
		fmt.Printf("testName should start with Test\n")
		return subcommands.ExitFailure
	}

	if !regexp.MustCompile("[_a-zA-Z][_a-zA-Z0-9]{0,100}").MatchString(testName) {
		fmt.Printf("testName is not a valid identifier\n")
		return subcommands.ExitFailure
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

	var testFunctionNamePattern = "func " + testName + "("

	var matchDirs = map[string]struct{}{}

	for _, goTestSourceFile := range goTestSourceFiles {
		var content, err = os.ReadFile(goTestSourceFile)
		if err != nil {
			fmt.Printf("Read %v error: %v", goTestSourceFile, err)
			continue
		}

		if !strings.Contains(string(content), testFunctionNamePattern) {
			continue
		}

		// match

		var dirRelativePath = "./" + filepath.Dir(goTestSourceFile)
		dirRelativePath = strings.ReplaceAll(dirRelativePath, "\\", "/")

		matchDirs[dirRelativePath] = struct{}{}
	}

	if len(matchDirs) != 1 {
		fmt.Printf("Expecting pattern found 1 packages, actual: %v, pattern: %v\n",
			len(matchDirs), testFunctionNamePattern)
		return subcommands.ExitFailure
	}

	var matchDir string

	for k := range matchDirs {
		matchDir = k
		break
	}

	fmt.Printf("matchDir: %v", matchDir)

	return subcommands.ExitSuccess
}
