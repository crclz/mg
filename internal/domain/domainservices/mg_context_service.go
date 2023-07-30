package domainservices

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/crclz/mg/internal/domain/domainmodels"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

type MgContextService struct {
}

// constructor
func NewMgContextService() *MgContextService {
	return &MgContextService{}
}

// wire
var singletonMgContextService *MgContextService = initSingletonMgContextService()

func GetSingletonMgContextService() *MgContextService {
	return singletonMgContextService
}

func initSingletonMgContextService() *MgContextService {
	return NewMgContextService()
}

// methods

func (p *MgContextService) UsingMgContextNameFile() string {
	return "mg-temp/using-context.txt"
}

func (p *MgContextService) SetUsingMgContextName(ctx context.Context, dir string, contextName string) error {
	var usingContextNameFile = filepath.Join(dir, p.UsingMgContextNameFile())

	var err = os.MkdirAll(filepath.Dir(usingContextNameFile), 0644)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	err = os.WriteFile(usingContextNameFile, []byte(contextName), 0644)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (p *MgContextService) GetUsingMgContextName(
	ctx context.Context, dir string,
) (string, error) {
	var contentBytes, err = os.ReadFile(filepath.Join(dir, p.UsingMgContextNameFile()))
	if os.IsNotExist(err) {
		return "default", nil
	}

	var lines = strings.Split(strings.ReplaceAll(string(contentBytes), "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return "", xerrors.Errorf("Parse error: %v", p.UsingMgContextNameFile())
	}

	var contextName = strings.TrimSpace(lines[0])

	return contextName, nil
}

func (p *MgContextService) GetUsingMgContext(
	ctx context.Context, dir string,
) (*domainmodels.MgContextConfig, error) {
	var usingContextName, err = p.GetUsingMgContextName(ctx, dir)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	var yamlFile = fmt.Sprintf("mg-context.%v.yaml", usingContextName)
	content, err := os.ReadFile(yamlFile)
	if err != nil {
		if os.IsNotExist(err) && usingContextName == "default" {
			return &domainmodels.MgContextConfig{
				Go: &domainmodels.MgContextConfigGoConfig{
					GoTestPrefix:   []string{},
					GoBuildNoOptim: false,
				},
			}, nil
		}
		return nil, xerrors.Errorf(": %w", err)
	}

	var result = &domainmodels.MgContextConfig{}
	err = yaml.Unmarshal(content, result)
	if err != nil {
		return nil, xerrors.Errorf("Yaml unmarshal error when reading %v: %w", yamlFile, err)
	}

	return result, nil
}
