package domainutils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
	"golang.org/x/xerrors"
)

var ErrAlpha = errors.New("alpha error")
var ErrBeta = errors.New("beta error")

// 让error同时有包装、stack trace和类型
func TestErrorTyping_WrapWithBetaType_fail(t *testing.T) {
	var assert = require.New(t)

	var err = xerrors.Errorf(": %w", ErrAlpha)
	err = xerrors.Errorf("Error Beta %w happens because of alpha: %w", ErrBeta, err)

	t.Logf("err: %+v", err)

	assert.True(errors.Is(err, ErrAlpha))
	assert.False(errors.Is(err, ErrBeta))
}

func TestErrorTyping_WrapWithBetaType_success(t *testing.T) {
	var assert = require.New(t)

	var err = xerrors.Errorf(": %w", ErrAlpha)
	err = multierr.Append(ErrBeta, err)

	t.Logf("err: %+v", err)

	assert.True(errors.Is(err, ErrBeta))
	assert.True(errors.Is(err, ErrAlpha))
}
