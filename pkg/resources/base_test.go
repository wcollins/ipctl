// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindByNameFound(t *testing.T) {
	items := []string{"alpha", "beta", "gamma"}
	result, err := FindByName(items, "item", "beta", func(s string) string { return s })
	assert.NoError(t, err)
	assert.Equal(t, "beta", *result)
}

func TestFindByNameNotFound(t *testing.T) {
	items := []string{"alpha", "beta"}
	result, err := FindByName(items, "item", "missing", func(s string) string { return s })
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestFindByNameNotFoundErrorIs(t *testing.T) {
	items := []string{"alpha"}
	_, err := FindByName(items, "template", "missing", func(s string) string { return s })
	assert.True(t, errors.Is(err, ErrNotFound), "expected errors.Is(err, ErrNotFound) to be true")
}

func TestFindByNameNotFoundErrorMessage(t *testing.T) {
	items := []string{"alpha"}
	_, err := FindByName(items, "template", "my-template", func(s string) string { return s })
	assert.EqualError(t, err, "template 'my-template' not found")
}

func TestFindByNameEmpty(t *testing.T) {
	var items []string
	_, err := FindByName(items, "template", "any", func(s string) string { return s })
	assert.True(t, errors.Is(err, ErrNotFound))
}
