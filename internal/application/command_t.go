package application

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/google/subcommands"
)

type PrintCmd struct {
	capitalize bool
}

func (*PrintCmd) Name() string     { return "print" }
func (*PrintCmd) Synopsis() string { return "Print args to stdout." }
func (*PrintCmd) Usage() string {
	return `print [-capitalize] <some text>:
	Print args to stdout.
  `
}

func (p *PrintCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.capitalize, "capitalize", false, "capitalize output")
}

func (p *PrintCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	for _, arg := range f.Args() {
		if p.capitalize {
			arg = strings.ToUpper(arg)
		}
		fmt.Printf("%s ", arg)
	}
	fmt.Println()
	return subcommands.ExitSuccess
}
