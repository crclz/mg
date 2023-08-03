package application

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"text/template"

	"github.com/crclz/mg/internal/domain/domainmodels"
	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/google/subcommands"
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

func (p *GenerateCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	switch arg0 := f.Arg(0); arg0 {
	case "s":
		// generate service
		var servicePath = f.Arg(1)
		var tmpl, err = template.New("").Parse(domainmodels.GetAssetTemplateSimpleSingletonService())
		if err != nil {
			panic(err)
		}

		var generatedCode = bytes.NewBufferString("")

		err = tmpl.ExecuteTemplate(generatedCode, "SimpleSingletonService", map[string]any{
			"PackageName":      "some_package",
			"ServiceClassName": "SomeService",
		})

		if err != nil {
			panic(err)
		}

		fmt.Printf("%v: \n%v\n", servicePath, generatedCode.String())
		return subcommands.ExitSuccess
	default:
		fmt.Printf("Unknown arg0: %v", arg0)
		return subcommands.ExitFailure
	}
}
