package domainmodels

import _ "embed"

//go:embed Template_simple_singleton_service.txt
var asset_template_simple_singleton_service string

func GetAssetTemplateSimpleSingletonService() string {
	return asset_template_simple_singleton_service
}
