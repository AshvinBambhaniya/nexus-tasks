package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	email := "xyz@nexustask.com"

	valid, err := ValidateEmail(email)
	assert.NoError(t, err)
	assert.True(t, valid)

	email = "xyz@nexustask.com"
	valid, err = ValidateEmail(email)
	assert.NoError(t, err)
	assert.False(t, valid)
}
