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

func (p *MgContextService) PatternOfContextName() string {
	return `^[a-z][a-z0-9-]{0,63}$`
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

	contextConfig, err := p.ReadContextConfigFromDisk(ctx, dir, usingContextName)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if contextConfig == nil && usingContextName == "default" {
		contextConfig = domainmodels.NewMgContextConfig()
	}

	if contextConfig == nil {
		return nil, xerrors.Errorf("context file not exist: %v", usingContextName)
	}

	return contextConfig, nil
}

func (p *MgContextService) ContextConfigFilePath(dir string, contextName string) string {
	return fmt.Sprintf("mg-context.%v.yaml", contextName)
}

func (p *MgContextService) ReadContextConfigFromDisk(
	ctx context.Context, dir string, contextName string,
) (*domainmodels.MgContextConfig, error) {
	var yamlFile = p.ContextConfigFilePath(dir, contextName)

	content, err := os.ReadFile(yamlFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
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

func (p *MgContextService) WriteContextConfigToDisk(
	ctx context.Context, dir string, contextName string, contextConfig *domainmodels.MgContextConfig,
) error {
	if contextConfig == nil {
		panic("contextConfig is nil")
	}

	var yamlFile = p.ContextConfigFilePath(dir, contextName)

	yamlBytes, err := yaml.Marshal(contextConfig)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	err = os.WriteFile(yamlFile, yamlBytes, 0644)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}
