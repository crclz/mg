package application

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

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

	var err = p.ExecuteInteral(ctx, f)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		// fmt.Printf("Error: %v")
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (p *MagicCommand) ExecuteInteral(ctx context.Context, f *flag.FlagSet) error {
	var pattern = regexp.MustCompile(`//[ ]?![ ]?magic`)
	const timeRange = 3 * 60
	const fileLimit = 5

	files, err := p.fileDiscoveryService.Discover(ctx, ".", "*.go", pattern, 3*60)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if len(files) == 0 {
		return xerrors.Errorf("zero file discovered. timeRange(sec): %v, pattern: %v", timeRange, pattern)
	}

	if len(files) > fileLimit {
		return xerrors.Errorf("too many file discovered. actual: %v, limit: %v, timeRange(sec): %v, pattern: %v",
			len(files), timeRange)
	}

	for _, filename := range files {
		var err = p.ApplyMagicToFile(ctx, filename)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	return nil
}

func (p *MagicCommand) PackageName(goFileName string) (string, error) {
	var packageName = filepath.Dir(goFileName)
	if packageName != "" {
		return packageName, nil
	}

	var err error

	goFileName, err = filepath.Abs(goFileName)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}

	packageName = filepath.Dir(goFileName)

	if packageName == "" {
		return "", xerrors.Errorf("Empty package name detected")
	}

	return packageName, nil
}

func (p *MagicCommand) ApplyMagicToFile(ctx context.Context, filename string) error {
	packageName, err := p.PackageName(filename)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	workingDirectoryAbs, err := filepath.Abs(".")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	filename, err = filepath.Abs(filename)
	if err != nil {
		return xerrors.Errorf("call Abs failure. filename: %v, error: %w", filename, err)
	}

	filename, err = filepath.Rel(workingDirectoryAbs, filename)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	fmt.Printf("Applying magic to file: %v\n", filename)

	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	var contentLines = strings.Split(strings.ReplaceAll(string(contentBytes), "\r\n", "\n"), "\n")

	contentLines, err = p.ApplyMagicTestClassMethod(ctx, packageName, contentLines)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// save
	err = os.WriteFile(filename, []byte(strings.Join(contentLines, "\n")), 0644)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

// replace !magic test SomeClass.SomeMethod one two three
//
// with func TestSomeClass_SomeMethod_oneTwoThree(t*testing.T){arrange, act, assert}
func (p *MagicCommand) ApplyMagicTestClassMethod(
	ctx context.Context, packageName string, contentLines []string,
) ([]string, error) {
	var result []string

	// regex tool: https://www.jyshare.com/front-end/854/
	var pattern = regexp.MustCompile(
		`[\s]*//[ ]?![ ]?magic[\s]+test[\s]+([^\.\s]+)[\.|\)][\s]*([^\.\s]+)[\s]+(.+)`,
	)

	for _, line := range contentLines {

		var matches = pattern.FindAllStringSubmatch(line, 1)
		if len(matches) == 0 {
			result = append(result, line)
			continue
		}

		var groups = matches[0]

		var serviceName = strings.TrimSpace(groups[1])
		var methodName = strings.TrimSpace(groups[2])
		var behaviorClause = groups[3]
		behaviorClause = p.FormatBehaviorCaluse(behaviorClause)

		var serviceVariableName = p.WithFirstLetter(serviceName, false)

		var templateText = `{{define "a"}}
func Test{{.ServiceName}}_{{.MethodName}}_{{.BehaviorClause}}(t* testing.T) {
	// arrange
	var assert = utils.ProdAssert(t)
	var ctx = context.Background()
	var {{.ServiceVariableName}} = {{.PackageName}}.GetSingleton{{.ServiceName}}()
	
	assert.NotNil(ctx)
	assert.NotNil({{.ServiceVariableName}})

	// act
	// {{.ServiceVariableName}}.{{.MethodName}}()

	// assert

}
{{end}}`

		templateObject, err := template.New("").Parse(templateText)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		var generatedCode = bytes.NewBufferString("")

		err = templateObject.ExecuteTemplate(generatedCode, "a", map[string]any{
			"PackageName":         packageName,
			"ServiceName":         serviceName,
			"MethodName":          methodName,
			"BehaviorClause":      behaviorClause,
			"ServiceVariableName": serviceVariableName,
		})

		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		result = append(result, strings.Split(generatedCode.String(), "\n")...)
	}

	return result, nil
}

// input: one two three, output: oneTwoThree
func (p *MagicCommand) FormatBehaviorCaluse(clause string) string {
	var builder = strings.Builder{}

	for i, item := range regexp.MustCompile(`[\s]`).Split(clause, -1) {
		if item == "" {
			continue
		}

		var capitalize = i != 0
		item = p.WithFirstLetter(item, capitalize)

		builder.Write([]byte(item))
	}

	return builder.String()
}

func (p *MagicCommand) WithFirstLetter(s string, capitalize bool) string {
	var firstLetter = s[:1]
	var otherletters = s[1:]

	if capitalize {
		firstLetter = strings.ToUpper(firstLetter)
	} else {
		firstLetter = strings.ToLower(firstLetter)
	}

	return firstLetter + otherletters
}
