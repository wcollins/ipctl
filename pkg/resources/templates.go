// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// TemplateResource provides business logic for template operations.
type TemplateResource struct {
	BaseResource
	service services.TemplateServicer
}

// NewTemplateResource creates a new TemplateResource with the given service.
func NewTemplateResource(svc services.TemplateServicer) TemplateResourcer {
	return &TemplateResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetByName retrieves a template by name using client-side filtering.
// It fetches all templates and searches for an exact name match.
func (r *TemplateResource) GetByName(name string) (*services.Template, error) {
	logging.Trace()

	templates, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return FindByName(templates, "template", name, func(t services.Template) string {
		return t.Name
	})
}

// GetAll retrieves all templates from the API.
// This is a pass-through to the service layer for pure API access.
func (r *TemplateResource) GetAll() ([]services.Template, error) {
	return r.service.GetAll()
}

// Get retrieves a template by its ID from the API.
// This is a pass-through to the service layer for pure API access.
func (r *TemplateResource) Get(id string) (*services.Template, error) {
	return r.service.Get(id)
}

// Create creates a new template.
// This is a pass-through to the service layer for pure API access.
func (r *TemplateResource) Create(in services.Template) (*services.Template, error) {
	return r.service.Create(in)
}

// Delete removes a template by its ID.
// This is a pass-through to the service layer for pure API access.
func (r *TemplateResource) Delete(id string) error {
	return r.service.Delete(id)
}

// Export exports a template by its ID.
// This is a pass-through to the service layer for pure API access.
func (r *TemplateResource) Export(id string) (*services.Template, error) {
	return r.service.Export(id)
}

// Import imports a template into the system.
// This is a pass-through to the service layer for pure API access.
func (r *TemplateResource) Import(in services.Template) (*services.Template, error) {
	return r.service.Import(in)
}
