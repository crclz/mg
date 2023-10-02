package application

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/crclz/mg/internal/domain/domainutils"
	"github.com/google/subcommands"
	"golang.org/x/xerrors"
)

type TeCommand struct {
	mgContextService *domainservices.MgContextService

	// flags
	script   bool
	countOne bool
}

func NewTeCommand(
	mgContextService *domainservices.MgContextService,
) *TeCommand {
	return &TeCommand{
		mgContextService: mgContextService,
	}
}

func (*TeCommand) Name() string     { return "t" }
func (*TeCommand) Synopsis() string { return "Run tests with convenience" }
func (*TeCommand) Usage() string {
	return "refer to readme.md"
}

func (p *TeCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.script, "script", false, "Will set environment variable $GoScriptName")
	f.BoolVar(&p.countOne, "c1", false, "will add --count=1 to go test command")
}

func (p *TeCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	ctx, cancelCause := context.WithCancelCause(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var sig = <-sigChan
		cancelCause(xerrors.Errorf("Received signal: %v", sig))
	}()

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

	// TODO: prefix
	mgContext, err := p.mgContextService.GetUsingMgContext(ctx, ".")
	if err != nil {
		fmt.Printf("GetUsingMgContext error: %v\n", err)
		return subcommands.ExitFailure
	}

	// TODO: optim build flags
	var goTestCommand = []string{}
	goTestCommand = append(goTestCommand, mgContext.Go.GoTestPrefix...)
	goTestCommand = append(goTestCommand, "go", "test")

	if mgContext.Go.GoBuildNoOptim {
		goTestCommand = append(goTestCommand, `--gcflags`, `all=-l -N`)
	}

	goTestCommand = append(goTestCommand, "-v", matchDir, "--run", "^"+testName+"$")

	if p.countOne {
		goTestCommand = append(goTestCommand, "--count=1")
	}

	var commandString = ""

	for _, part := range goTestCommand {
		if strings.Contains(part, " ") {
			part = "\"" + part + "\""
		}

		commandString += " " + part
	}

	commandString = strings.TrimSpace(commandString)

	fmt.Printf("Command array: %v\n", domainutils.ToJson(goTestCommand))
	fmt.Printf("Command string: %v\n", commandString)

	var commandObject = exec.CommandContext(ctx, goTestCommand[0], goTestCommand[1:]...)
	commandObject.Stdout = os.Stdout
	commandObject.Stderr = os.Stderr
	commandObject.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}

	if p.script {
		commandObject.Env = append(os.Environ(), "GoScriptName="+testName)
	}

	defer func() {

		if !commandObject.ProcessState.Exited() {
			var pgid = -1 * commandObject.Process.Pid

			// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
			// Solution: In addition to sending a signal to a single PID,
			// kill(2) also supports sending a signal to a Process Group by passing
			// the process group id (PGID) as a negative number.
			var err = syscall.Kill(pgid, syscall.SIGKILL)
			if err != nil {
				fmt.Printf("Kill pgid %v error: %+v\n", pgid, err)
			} else {
				fmt.Printf("Kill pgid %v success\n", pgid)
			}
		}
	}()

	err = commandObject.Run()
	if err != nil {
		fmt.Printf("go test failure status. error: %v\n", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
