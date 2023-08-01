package application

import (
	"context"
	"flag"
	"fmt"
	"regexp"

	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/google/subcommands"
	"golang.org/x/xerrors"
)

type UseContextCommand struct {
	mgContextService *domainservices.MgContextService

	// flags
	query bool
}

func NewUseContextCommand(
	mgContextService *domainservices.MgContextService,
) *UseContextCommand {
	return &UseContextCommand{
		mgContextService: mgContextService,
	}
}

func (*UseContextCommand) Name() string { return "use-context" }
func (*UseContextCommand) Synopsis() string {
	return "use a mg context, or query the current context name."
}
func (*UseContextCommand) Usage() string {
	return "use-context $contextName, or use-context --query\n"
}

func (p *UseContextCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.query, "query", false, "pass this argument to query currently using context")
}

func (p *UseContextCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.query {
		var err = p.QueryCurrentContext(ctx)
		if err != nil {
			fmt.Printf("QueryCurrentContext failure: %v", err)
			return subcommands.ExitFailure
		}
		return subcommands.ExitSuccess
	}

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

	if config == nil && contextName != "default" {
		fmt.Printf("Context not exist: %v\n", contextName)
		return subcommands.ExitFailure
	}

	err = p.mgContextService.SetUsingMgContextName(ctx, ".", contextName)
	if err != nil {
		fmt.Printf("SetUsingMgContextName error: %v\n", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("Change using context to: %v\n", contextName)

	return subcommands.ExitSuccess
}

func (p *UseContextCommand) QueryCurrentContext(ctx context.Context) error {
	name, err := p.mgContextService.GetUsingMgContextName(ctx, ".")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	fmt.Printf("Currently using context is: %v\n", name)
	return nil
}
