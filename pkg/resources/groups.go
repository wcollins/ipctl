// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// GroupResource provides business logic for authorization group operations.
type GroupResource struct {
	BaseResource
	service services.GroupServicer
}

// NewGroupResource creates a new GroupResource with the given service.
func NewGroupResource(svc services.GroupServicer) GroupResourcer {
	return &GroupResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetAll retrieves all groups from the API.
// This is a pass-through to the service layer for pure API access.
func (r *GroupResource) GetAll() ([]services.Group, error) {
	return r.service.GetAll()
}

// Get retrieves a specific group by name from the API.
// This is a pass-through to the service layer for pure API access.
func (r *GroupResource) Get(name string) (*services.Group, error) {
	return r.service.Get(name)
}

// Create creates a new group.
// This is a pass-through to the service layer for pure API access.
func (r *GroupResource) Create(in services.Group) (*services.Group, error) {
	return r.service.Create(in)
}

// Delete removes a group by its identifier.
// This is a pass-through to the service layer for pure API access.
func (r *GroupResource) Delete(id string) error {
	return r.service.Delete(id)
}

// GetByName retrieves a group by name using client-side filtering.
// It fetches all groups and searches for a matching name.
func (r *GroupResource) GetByName(name string) (*services.Group, error) {
	logging.Trace()

	groups, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return FindByName(groups, "group", name, func(g services.Group) string {
		return g.Name
	})
}
