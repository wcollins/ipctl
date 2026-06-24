// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"errors"
	"fmt"
)

// ErrNotFound is the sentinel error returned when a resource cannot be found by name.
// Callers should use errors.Is(err, resources.ErrNotFound) to check for this condition.
var ErrNotFound = errors.New("not found")

// NotFoundError is returned by FindByName when no item matches the given name.
type NotFoundError struct {
	Kind string
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s '%s' not found", e.Kind, e.Name)
}

func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// BaseResource provides common functionality for resource operations.
type BaseResource struct{}

// NewBaseResource creates a new BaseResource instance.
func NewBaseResource() BaseResource {
	return BaseResource{}
}

// FindByName is a generic helper that searches through a slice of items
// and returns the first item where the getName function returns a matching name.
// Returns a *NotFoundError (which satisfies errors.Is(err, ErrNotFound)) when not found.
func FindByName[T any](items []T, kind, name string, getName func(T) string) (*T, error) {
	for i := range items {
		if getName(items[i]) == name {
			return &items[i], nil
		}
	}
	return nil, &NotFoundError{Kind: kind, Name: name}
}

// DeleteAll is a generic helper that deletes all items using the provided delete function.
func DeleteAll[T any](items []T, getId func(T) string, deleteFn func(string) error) error {
	for _, item := range items {
		if err := deleteFn(getId(item)); err != nil {
			return err
		}
	}
	return nil
}

// ValidateNotEmpty validates that a string field is not empty.
func ValidateNotEmpty(value, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidateGbacRules validates GBAC (Group-Based Access Control) rules.
// Returns an error if read groups are configured but write groups are not.
func ValidateGbacRules(readGroups, writeGroups []interface{}) error {
	if len(readGroups) > 0 && len(writeGroups) == 0 {
		return errors.New("write group must be configured, when read group present")
	}
	return nil
}
