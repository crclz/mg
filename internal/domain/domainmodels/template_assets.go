package domainmodels

import _ "embed"

//go:embed template_simple_singleton_service.txt
var asset_template_simple_singleton_service string

func GetAssetTemplateSimpleSingletonService() string {
	return asset_template_simple_singleton_service
}

//go:embed template_testing.txt
var template_testing string

func GetAssetTemplateTesting() string {
	return template_testing
}
