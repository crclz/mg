package domainservices

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/xerrors"
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
//
// limit: -1 for no limit
func (p *FileDiscoveryService) Discover(
	ctx context.Context, rootDirectory string, filenamePattern string,
	linePattern *regexp.Regexp, modifyTimeRange time.Duration, limit int64,
) ([]string, error) {
	var err error

	if limit == 0 {
		return nil, nil
	}

	if limit < 0 {
		limit = 9999999999
	}

	rootDirectory, err = filepath.Abs(rootDirectory)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	type FileEntry struct {
		Path string
		Info fs.FileInfo
	}

	var fileEntries []*FileEntry

	filepath.Walk(rootDirectory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("filepath.Walk error. path: %v, error: %v", path, err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		matched, err := filepath.Match(filenamePattern, info.Name())
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}

		if matched {
			if info.ModTime().After(time.Now().Add(-modifyTimeRange)) {
				fileEntries = append(fileEntries, &FileEntry{
					Path: path,
					Info: info,
				})
			}
		}

		return nil
	})

	// log.Printf("fileEntries: %v", len(fileEntries))

	sort.Slice(fileEntries, func(i, j int) bool {
		// order by mod time desc
		return fileEntries[i].Info.ModTime().After(fileEntries[j].Info.ModTime())
	})

	var result []string

	for _, entry := range fileEntries {
		if len(result) >= int(limit) {
			break
		}

		contentBytes, err := os.ReadFile(entry.Path)
		if err != nil {
			fmt.Printf("Read file error: %v, path: %v", err, entry.Path)
		}

		var lines = strings.Split(strings.ReplaceAll(string(contentBytes), "\r\n", "\n"), "\n")

		for _, line := range lines {
			if linePattern.MatchString(line) {
				result = append(result, entry.Path)
				continue
			}
		}
	}

	return result, nil
}
