package domainmodels

type MgContextConfig struct {
	Help *MgContextConfigHelp     `yaml:"Help"`
	Go   *MgContextConfigGoConfig `yaml:"Go"`
}

func NewMgContextConfig() *MgContextConfig {
	return &MgContextConfig{
		Help: NewMgContextConfigHelp(),
		Go: &MgContextConfigGoConfig{
			GoTestPrefix:   []string{},
			GoBuildNoOptim: false,
		},
	}
}

type MgContextConfigHelp struct {
	GoInstallCommand string `yaml:"GoInstallCommand"`
}

func NewMgContextConfigHelp() *MgContextConfigHelp {
	return &MgContextConfigHelp{
		GoInstallCommand: "go install github.com/crclz/mg@latest",
	}
}

type MgContextConfigGoConfig struct {
	GoTestPrefix   []string `yaml:"GoTestPrefix"`
	GoBuildNoOptim bool     `yaml:"GoBuildNoOptim"`
}
