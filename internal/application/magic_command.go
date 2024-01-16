package application

import (
	"context"
	"flag"
	"regexp"

	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/google/subcommands"
	"golang.org/x/xerrors"
)

type MagicCommand struct {
	mgContextService     *domainservices.MgContextService
	fileDiscoveryService *domainservices.FileDiscoveryService
}

func NewMagicCommand(
	mgContextService *domainservices.MgContextService,
	fileDiscoveryService *domainservices.FileDiscoveryService,
) *MagicCommand {
	return &MagicCommand{
		mgContextService:     mgContextService,
		fileDiscoveryService: fileDiscoveryService,
	}
}

func (*MagicCommand) Name() string { return "magic" }
func (*MagicCommand) Synopsis() string {
	return "magic command"
}
func (*MagicCommand) Usage() string {
	return "refer to readme.md\n"
}

func (p *MagicCommand) SetFlags(f *flag.FlagSet) {
}

func (p *MagicCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	panic("NotImplemented")
}

func (p *MagicCommand) ExecuteInteral(ctx context.Context) error {
	var pattern = regexp.MustCompile("//[ ]!")
	files, err := p.fileDiscoveryService.Discover(ctx, ".", pattern, 3*60)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if len(files) == 0 ||  {
	}

}
