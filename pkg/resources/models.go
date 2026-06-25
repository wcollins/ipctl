// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package resources

import (
	"fmt"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/services"
)

// ModelResource provides business logic for lifecycle manager model operations.
type ModelResource struct {
	BaseResource
	service               services.ModelServicer
	workflowService       services.WorkflowServicer
	transformationService services.TransformationServicer
	instanceService       services.InstanceServicer
}

// NewModelResource creates a new ModelResource with the given services.
func NewModelResource(
	svc services.ModelServicer,
	wfSvc services.WorkflowServicer,
	jstSvc services.TransformationServicer,
	instSvc services.InstanceServicer,
) ModelResourcer {
	return &ModelResource{
		BaseResource:          NewBaseResource(),
		service:               svc,
		workflowService:       wfSvc,
		transformationService: jstSvc,
		instanceService:       instSvc,
	}
}

// DeleteOptions contains options for model deletion.
type DeleteOptions struct {
	DeleteInstances bool
	DeleteRelated   bool
}

// GetByName retrieves a model by name using client-side filtering.
// It fetches all models and searches for a matching name.
func (r *ModelResource) GetByName(name string) (*services.Model, error) {
	logging.Trace()

	models, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return FindByName(models, "model", name, func(m services.Model) string {
		return m.Name
	})
}

// GetAll retrieves all models from the API.
// This is a pass-through to the service layer for pure API access.
func (r *ModelResource) GetAll() ([]services.Model, error) {
	return r.service.GetAll()
}

// Create creates a new model.
// This is a pass-through to the service layer for pure API access.
func (r *ModelResource) Create(in services.Model) (*services.Model, error) {
	return r.service.Create(in)
}

// Delete removes a model by its ID.
// This is a pass-through to the service layer for pure API access.
func (r *ModelResource) Delete(id string, deleteInstances bool) error {
	return r.service.Delete(id, deleteInstances)
}

// DeleteWithOptions removes a model with advanced options including checking for instances
// and deleting related workflows and transformations.
func (r *ModelResource) DeleteWithOptions(model *services.Model, opts DeleteOptions) error {
	logging.Trace()

	// Check for attached instances if not forcing deletion
	if !opts.DeleteInstances {
		instances, err := r.GetInstances(model.Id)
		if err != nil {
			return fmt.Errorf("checking for instances: %w", err)
		}

		if len(instances) > 0 {
			return fmt.Errorf("model `%s` has %d attached instances, use --delete-instances to delete all instances", model.Name, len(instances))
		}
	}

	// Delete related workflows and transformations if requested
	if opts.DeleteRelated {
		if err := r.deleteRelatedResources(model); err != nil {
			return fmt.Errorf("deleting related resources: %w", err)
		}
	}

	// Finally delete the model
	return r.service.Delete(model.Id, opts.DeleteInstances)
}

// deleteRelatedResources deletes workflows and transformations associated with a model.
func (r *ModelResource) deleteRelatedResources(model *services.Model) error {
	for _, action := range model.Actions {
		// Delete associated workflow
		if action.Workflow != nil && *action.Workflow != "" {
			if err := r.deleteWorkflowIfExists(*action.Workflow); err != nil {
				return fmt.Errorf("deleting workflow %s: %w", *action.Workflow, err)
			}
		}

		// Delete pre-workflow transformation
		if action.PreWorkflowJst != nil && *action.PreWorkflowJst != "" {
			if err := r.deleteTransformationIfExists(*action.PreWorkflowJst); err != nil {
				logging.Warn("error deleting pre-workflow transformation %s: %v", *action.PreWorkflowJst, err)
			}
		}

		// Delete post-workflow transformation
		if action.PostWorkflowJst != nil && *action.PostWorkflowJst != "" {
			if err := r.deleteTransformationIfExists(*action.PostWorkflowJst); err != nil {
				logging.Warn("error deleting post-workflow transformation %s: %v", *action.PostWorkflowJst, err)
			}
		}
	}

	return nil
}

// deleteWorkflowIfExists deletes a workflow by ID if it exists.
func (r *ModelResource) deleteWorkflowIfExists(workflowId string) error {
	workflow, err := r.workflowService.GetById(workflowId)
	if err != nil {
		if err.Error() == "workflow not found" {
			return nil // Already deleted or doesn't exist
		}
		return err
	}

	if workflow != nil {
		return r.workflowService.Delete(workflow.Name)
	}

	return nil
}

// deleteTransformationIfExists deletes a transformation by ID if it exists.
func (r *ModelResource) deleteTransformationIfExists(transformationId string) error {
	transformation, err := r.transformationService.Get(transformationId)
	if err != nil {
		if err.Error() == "transformation not found" {
			return nil // Already deleted or doesn't exist
		}
		return err
	}

	if transformation != nil {
		return r.transformationService.Delete(transformation.Id)
	}

	return nil
}

// GetInstances retrieves all instances for a given model ID.
func (r *ModelResource) GetInstances(modelId string) ([]services.Instance, error) {
	return r.instanceService.GetAll(modelId)
}

// Export exports a model by its ID.
// This is a pass-through to the service layer for pure API access.
func (r *ModelResource) Export(id string) (*services.Model, error) {
	return r.service.Export(id)
}

// Import imports a model into the system.
// This is a pass-through to the service layer for pure API access.
func (r *ModelResource) Import(in services.Model) (*services.Model, error) {
	return r.service.Import(in)
}
