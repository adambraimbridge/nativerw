package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {
	optional := NewOptional(func() (interface{}, error) {
		time.Sleep(500 * time.Millisecond)
		return "hi", nil
	})

	isNil := optional.Nil()
	assert.True(t, isNil)

	result := optional.Get()
	assert.Nil(t, result)

	result, err := optional.Block()
	assert.NoError(t, err)
	assert.Equal(t, "hi", result)

	isNil = optional.Nil()
	assert.False(t, isNil)

	result = optional.Get()
	assert.Equal(t, "hi", result)
}
