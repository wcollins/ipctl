// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// WorkflowResource provides business logic for workflow operations.
type WorkflowResource struct {
	BaseResource
	service services.WorkflowServicer
}

// NewWorkflowResource creates a new WorkflowResource with the given service.
func NewWorkflowResource(svc services.WorkflowServicer) WorkflowResourcer {
	return &WorkflowResource{
		BaseResource: NewBaseResource(),
		service:      svc,
	}
}

// GetAll retrieves all workflows from the API.
// This is a pass-through to the service layer for pure API access.
func (r *WorkflowResource) GetAll() ([]services.Workflow, error) {
	return r.service.GetAll()
}

// Get retrieves a specific workflow by name from the API.
// This is a pass-through to the service layer for pure API access.
func (r *WorkflowResource) Get(name string) (*services.Workflow, error) {
	return r.service.Get(name)
}

// Create creates a new workflow.
// This is a pass-through to the service layer for pure API access.
func (r *WorkflowResource) Create(in services.Workflow) (*services.Workflow, error) {
	return r.service.Create(in)
}

// Delete removes a workflow by its name.
// This is a pass-through to the service layer for pure API access.
func (r *WorkflowResource) Delete(name string) error {
	return r.service.Delete(name)
}

// Import imports a workflow.
// This is a pass-through to the service layer for pure API access.
func (r *WorkflowResource) Import(in services.Workflow) (*services.Workflow, error) {
	return r.service.Import(in)
}

// Export exports a workflow by name.
// This is a pass-through to the service layer for pure API access.
func (r *WorkflowResource) Export(name string) (*services.Workflow, error) {
	return r.service.Export(name)
}

// ExportById exports a workflow by ID.
// This is a pass-through to the service layer for pure API access.
func (r *WorkflowResource) ExportById(id string) (*services.Workflow, error) {
	return r.service.ExportById(id)
}

// Update updates an existing workflow.
// This is a pass-through to the service layer for pure API access.
func (r *WorkflowResource) Update(in services.Workflow) (*services.Workflow, error) {
	return r.service.Update(in)
}

// GetById retrieves a workflow by ID using client-side filtering.
// It fetches all workflows and searches for a matching ID.
func (r *WorkflowResource) GetById(id string) (*services.Workflow, error) {
	logging.Trace()

	workflows, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return FindByName(workflows, "workflow", id, func(w services.Workflow) string {
		return w.Id
	})
}

// Clear deletes all workflows from the server.
// This is a bulk operation that orchestrates multiple delete calls.
func (r *WorkflowResource) Clear() error {
	logging.Trace()

	workflows, err := r.service.GetAll()
	if err != nil {
		return err
	}

	return DeleteAll(workflows, func(w services.Workflow) string {
		return w.Name
	}, r.service.Delete)
}
