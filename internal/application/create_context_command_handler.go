package application

import (
	"context"
	"flag"
	"fmt"
	"regexp"

	"github.com/crclz/mg/internal/domain/domainmodels"
	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/google/subcommands"
)

type CreateContextCommandHandler struct {
	mgContextService *domainservices.MgContextService
}

func NewCreateContextCommandHandler(
	mgContextService *domainservices.MgContextService,
) *CreateContextCommandHandler {
	return &CreateContextCommandHandler{
		mgContextService: mgContextService,
	}
}

func (*CreateContextCommandHandler) Name() string     { return "create-context" }
func (*CreateContextCommandHandler) Synopsis() string { return "create mg context" }
func (*CreateContextCommandHandler) Usage() string {
	return "create-context $contextName"
}

func (p *CreateContextCommandHandler) SetFlags(f *flag.FlagSet) {
}

func (p *CreateContextCommandHandler) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var positionalArgs = f.Args()
	if len(positionalArgs) >= 2 {
		fmt.Printf("Expecting 0/1 positional argument, but got %v.\n", len(positionalArgs))
		return subcommands.ExitFailure
	}

	var contextName = "default"

	if len(positionalArgs) == 1 {
		contextName = positionalArgs[0]
	}

	if !regexp.MustCompile(p.mgContextService.PatternOfContextName()).MatchString(contextName) {
		fmt.Printf("ContextName does not match: %v\n", p.mgContextService.PatternOfContextName())
		return subcommands.ExitFailure
	}

	// check existence
	config, err := p.mgContextService.ReadContextConfigFromDisk(ctx, ".", contextName)
	if err != nil {
		fmt.Printf("ReadContextConfigFromDisk error: %v\n", err)
		return subcommands.ExitFailure
	}

	if config != nil {
		fmt.Printf("Context already exist: %v\n", contextName)
		return subcommands.ExitFailure
	}

	// create context
	config = domainmodels.NewMgContextConfig()

	err = p.mgContextService.WriteContextConfigToDisk(ctx, ".", contextName, config)
	if err != nil {
		fmt.Printf("WriteContextConfigToDisk error: %v\n", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("Created file: %v\n", p.mgContextService.ContextConfigFilePath(".", contextName))

	return subcommands.ExitSuccess
}
