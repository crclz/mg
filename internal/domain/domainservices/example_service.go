package domainservices

import (
	"context"
	"strings"

	"github.com/bytedance/sonic"
	"golang.org/x/xerrors"
)

type ExampleService struct {
}

// Constructor of ExampleService
func NewExampleService() *ExampleService {
	return &ExampleService{}
}

// wire

var singletonExampleService *ExampleService = initSingletonExampleService()

func GetSingletonExampleService() *ExampleService {
	return singletonExampleService
}

func initSingletonExampleService() *ExampleService {
	return NewExampleService()
}

// methods

func (p *ExampleService) HelloWorlds(ctx context.Context, worldsJson string) (string, error) {
	var worlds []string

	var err = sonic.UnmarshalString(worldsJson, &worlds)
	if err != nil {
		return "", xerrors.Errorf("Json unmarshal error: %w", err)
	}

	var message = strings.Builder{}

	_, err = message.WriteString("Hello: ")
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}

	for _, world := range worlds {
		message.WriteString(world + ", ")
	}

	return strings.TrimSpace(message.String()) + "!", nil
}
