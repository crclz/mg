package domainservices

import (
	"context"
	"regexp"
)

type FileDiscoveryService struct {
}

// Constructor of FileDiscoveryService
func NewFileDiscoveryService() *FileDiscoveryService {
	return &FileDiscoveryService{}
}

// wire

var singletonFileDiscoveryService *FileDiscoveryService = initSingletonFileDiscoveryService()

func GetSingletonFileDiscoveryService() *FileDiscoveryService {
	return singletonFileDiscoveryService
}

func initSingletonFileDiscoveryService() *FileDiscoveryService {
	return NewFileDiscoveryService()
}

// methods

// returns absolute paths
func (p *FileDiscoveryService) Discover(
	ctx context.Context, rootDirectory string, pattern *regexp.Regexp, modifyTimeRange int64,
) ([]string, error) {
	panic("NotImplemented")
}
