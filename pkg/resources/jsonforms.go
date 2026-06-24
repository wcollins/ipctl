// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// JsonFormResource provides business logic for JSON form operations.
type JsonFormResource struct {
	BaseResource
	service services.JsonFormServicer
}

// NewJsonFormResource creates a new JsonFormResource with the given service.
func NewJsonFormResource(svc services.JsonFormServicer) JsonFormResourcer {
	return &JsonFormResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetByName retrieves a JSON form by name using client-side filtering.
// It fetches all forms and searches for a matching name.
func (r *JsonFormResource) GetByName(name string) (*services.JsonForm, error) {
	logging.Trace()

	forms, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return FindByName(forms, "json form", name, func(f services.JsonForm) string {
		return f.Name
	})
}

// Clear deletes all JSON forms from the server.
// This is a bulk operation that collects all form IDs and performs a single delete.
func (r *JsonFormResource) Clear() error {
	logging.Trace()

	forms, err := r.service.GetAll()
	if err != nil {
		return err
	}

	if len(forms) == 0 {
		return nil
	}

	var ids []string
	for _, form := range forms {
		ids = append(ids, form.Id)
	}

	return r.service.Delete(ids)
}

// GetAll retrieves all JSON forms from the API.
// This is a pass-through to the service layer for pure API access.
func (r *JsonFormResource) GetAll() ([]services.JsonForm, error) {
	return r.service.GetAll()
}

// Get retrieves a JSON form by its ID from the API.
// This is a pass-through to the service layer for pure API access.
func (r *JsonFormResource) Get(id string) (*services.JsonForm, error) {
	return r.service.Get(id)
}

// Create creates a new JSON form.
// This is a pass-through to the service layer for pure API access.
func (r *JsonFormResource) Create(in services.JsonForm) (*services.JsonForm, error) {
	return r.service.Create(in)
}

// Delete removes one or more JSON forms by their IDs.
// This is a pass-through to the service layer for pure API access.
func (r *JsonFormResource) Delete(ids []string) error {
	return r.service.Delete(ids)
}

// Import imports a JSON form into the system.
// This is a pass-through to the service layer for pure API access.
func (r *JsonFormResource) Import(in services.JsonForm) (*services.JsonForm, error) {
	return r.service.Import(in)
}
