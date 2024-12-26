package derrors

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestEquality(t *testing.T) {
	err1 := ErrDatabaseIssue
	assert.False(t, errors.Is(err1, ErrOrderProductNotFound))
}
