package main

import (
	"context"
	"flag"
	"os"

	"github.com/crclz/mg/internal/application"
	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(application.NewTeCommand(domainservices.GetSingletonMgContextService()),
		"Testing")

	subcommands.Register(application.NewCreateContextCommand(domainservices.GetSingletonMgContextService()),
		"Context Management")

	subcommands.Register(application.NewUseContextCommand(domainservices.GetSingletonMgContextService()),
		"Context Management")

	// generate
	subcommands.Register(application.NewGenerateCommand(domainservices.GetSingletonMgContextService()),
		"Code Generation")

	subcommands.Register(application.NewMagicCommand(
		domainservices.GetSingletonMgContextService(),
		domainservices.GetSingletonFileDiscoveryService(),
	), "Magic")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
