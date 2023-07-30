package domainmodels

type MgContextConfig struct {
	Go *MgContextConfigGoConfig `yaml:"Go"`
}

type MgContextConfigGoConfig struct {
	GoTestPrefix   []string `yaml:"GoTestPrefix"`
	GoBuildNoOptim bool     `yaml:"GoBuildNoOptim"`
}
