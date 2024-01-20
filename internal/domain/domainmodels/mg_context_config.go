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
	MgOfficial string `yaml:"MgOfficial"`
}

func NewMgContextConfigHelp() *MgContextConfigHelp {
	return &MgContextConfigHelp{
		MgOfficial: "The official page is: https://github.com/crclz/mg",
	}
}

type MgContextConfigGoConfig struct {
	GoTestPrefix    []string `yaml:"GoTestPrefix"`
	GoBuildNoOptim  bool     `yaml:"GoBuildNoOptim"`
	MeshTestCommand []string `yaml:"MeshTestCommand"`

	Magic MagicConfig `yaml:"Magic"`
}

type MagicConfig struct {
	TestArrangePart string `yaml:"TestArrangePart"`
}
