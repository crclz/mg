package application

import (
	"context"
	"flag"
	"fmt"
	"regexp"

	"github.com/crclz/mg/internal/domain/domainmodels"
	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/crclz/mg/internal/domain/domainutils"
	"github.com/google/subcommands"
	"golang.org/x/xerrors"
)

type CreateContextCommand struct {
	mgContextService *domainservices.MgContextService
}

func NewCreateContextCommand(
	mgContextService *domainservices.MgContextService,
) *CreateContextCommand {
	return &CreateContextCommand{
		mgContextService: mgContextService,
	}
}

func (*CreateContextCommand) Name() string     { return "create-context" }
func (*CreateContextCommand) Synopsis() string { return "create mg context" }
func (*CreateContextCommand) Usage() string {
	return "create-context $contextName"
}

func (p *CreateContextCommand) SetFlags(f *flag.FlagSet) {
}

func (p *CreateContextCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	exitStatus, err := p.ExecuteInteral(ctx, f)
	if err != nil {
		domainutils.ShowErrorToUser(err)

		exitStatus = subcommands.ExitFailure
	}

	return exitStatus
}

// return: success, error.
//   - error will be shown to user and overrides success
//   - success affects exit status
func (p *CreateContextCommand) ExecuteInteral(
	ctx context.Context, f *flag.FlagSet,
) (subcommands.ExitStatus, error) {

	var positionalArgs = f.Args()
	if len(positionalArgs) >= 2 {
		fmt.Printf("Expecting 0/1 positional argument, but got %v.\n", len(positionalArgs))
		return subcommands.ExitFailure, nil
	}

	var contextName = "default"

	if len(positionalArgs) == 1 {
		contextName = positionalArgs[0]
	}

	if !regexp.MustCompile(p.mgContextService.PatternOfContextName()).MatchString(contextName) {
		fmt.Printf("ContextName does not match: %v\n", p.mgContextService.PatternOfContextName())
		return subcommands.ExitFailure, nil
	}

	// check existence
	config, err := p.mgContextService.ReadContextConfigFromDisk(ctx, ".", contextName)
	if err != nil {
		fmt.Printf("ReadContextConfigFromDisk error: %v\n", err)
		return subcommands.ExitFailure, nil
	}

	if config != nil {
		fmt.Printf("Context already exist: %v\n", contextName)
		return subcommands.ExitFailure, nil
	}

	// create context
	config = domainmodels.NewMgContextConfig()

	err = p.mgContextService.WriteContextConfigToDisk(ctx, ".", contextName, config)
	if err != nil {
		return subcommands.ExitFailure, xerrors.Errorf(": %w", err)
	}

	fmt.Printf("Created file: %v\n", p.mgContextService.ContextConfigFilePath(".", contextName))

	return subcommands.ExitSuccess, nil
}
