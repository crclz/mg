package domainservices_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/crclz/mg/internal/domain/domainservices"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// go test ./internal/domain/domainservices --run TestFileDiscoveryService_Discover_happyCase1
func TestFileDiscoveryService_Discover_happyCase1(t *testing.T) {
	// arrange
	var assert = require.New(t)
	var ctx = context.Background()
	var fileDiscoveryService = domainservices.GetSingletonFileDiscoveryService()

	var err = os.WriteFile(".discover/2.txt", []byte(fmt.Sprintf("a_uuid=%v", string(uuid.NewString()))), 0644)
	assert.NoError(err)

	err = os.WriteFile(".discover/3.txt", []byte(fmt.Sprintf("b_uuid=%v", string(uuid.NewString()))), 0644)
	assert.NoError(err)

	// act
	result, err := fileDiscoveryService.Discover(ctx, "..", "*.txt", regexp.MustCompile("a_uuid="), time.Second*5, 10)
	assert.NoError(err)

	// assert
	assert.Len(result, 1)
	assert.True(strings.HasSuffix(result[0], "2.txt"), "result: %v", result)
}
