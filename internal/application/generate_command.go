package application

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/crclz/mg/internal/domain/domainmodels"
	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/google/subcommands"
	"golang.org/x/xerrors"
)

type GenerateCommand struct {
	mgContextService *domainservices.MgContextService
}

func NewGenerateCommand(
	mgContextService *domainservices.MgContextService,
) *GenerateCommand {
	return &GenerateCommand{
		mgContextService: mgContextService,
	}
}

func (*GenerateCommand) Name() string { return "g" }
func (*GenerateCommand) Synopsis() string {
	return "code generation"
}
func (*GenerateCommand) Usage() string {
	return "refer to readme.md\n"
}

func (p *GenerateCommand) SetFlags(f *flag.FlagSet) {
}

type ParseClassNameResult struct {
	Dir            string
	GoFileName     string
	GoTestFileName string
	PackageName    string
	ClassName      string
}

func (p *GenerateCommand) Parse(inputPath string) (*ParseClassNameResult, error) {
	// inputPath example: ./some_package/network_service
	// dir: ./some_package; file: network_service
	var dir, className = filepath.Split(inputPath)

	// className
	// support: network-service, NetworkService, network_service
	var classNameSplit = regexp.MustCompile("[-_]").Split(className, -1)
	className = ""
	for _, s := range classNameSplit {
		if s != "" {
			className += strings.ToUpper(s[0:1]) + s[1:]
		}
	}

	var goFileName = inputPath + ".go"
	var goTestFileName = inputPath + "_test.go"

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	var _, packageName = filepath.Split(absDir)
	if packageName == "" {
		return nil, xerrors.Errorf("Package name is empty, please check input path")
	}

	return &ParseClassNameResult{
		Dir:            dir,
		GoFileName:     goFileName,
		GoTestFileName: goTestFileName,
		PackageName:    packageName,
		ClassName:      className,
	}, nil
}

func (p *GenerateCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	switch arg0 := f.Arg(0); arg0 {
	case "s":
		// generate service
		var inputPath = f.Arg(1)
		var tmpl, err = template.New("").Parse(domainmodels.GetAssetTemplateSimpleSingletonService())
		if err != nil {
			panic(err)
		}

		parseResult, err := p.Parse(inputPath)
		if err != nil {
			fmt.Printf("Parse input path error: %v\n", err)
			return subcommands.ExitFailure
		}

		// check existence
		if _, err := os.Stat(parseResult.GoFileName); !os.IsNotExist(err) {
			fmt.Printf("Already exist: %v\n", parseResult.GoFileName)
			return subcommands.ExitFailure
		}

		var generatedCode = bytes.NewBufferString("")

		err = tmpl.ExecuteTemplate(generatedCode, "SimpleSingletonService", map[string]any{
			"PackageName":      parseResult.PackageName,
			"ServiceClassName": parseResult.ClassName,
		})
		if err != nil {
			panic(err)
		}

		// write go source file
		err = os.MkdirAll(parseResult.Dir, 0700)
		if err != nil {
			fmt.Printf("MkdirAll error. Dir: %v, error: %v\n", parseResult.Dir, err)
			return subcommands.ExitFailure
		}

		err = os.WriteFile(parseResult.GoFileName, generatedCode.Bytes(), 0644)
		if err != nil {
			fmt.Printf("Write file error: %v\n", parseResult.Dir)
			return subcommands.ExitFailure
		}

		fmt.Printf("Created: %v => %v\n", parseResult.ClassName, parseResult.GoFileName)

		return subcommands.ExitSuccess
	default:
		fmt.Printf("Unknown arg0: %v", arg0)
		return subcommands.ExitFailure
	}
}
