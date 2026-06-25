// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/resources"
	"github.com/itential/ipctl/pkg/services"
)

type ModelRunner struct {
	BaseRunner
	resource resources.ModelResourcer
	client   client.Client
}

func NewModelRunner(client client.Client, cfg config.Provider) *ModelRunner {
	return &ModelRunner{
		BaseRunner: NewBaseRunner(client, cfg),
		resource: resources.NewModelResource(
			services.NewModelService(client),
			services.NewWorkflowService(client),
			services.NewTransformationService(client),
			services.NewInstanceService(client),
		),
		client: client,
	}
}

/*
*******************************************************************************
Reader interface
*******************************************************************************
*/

// Get implements the `get model ...` command
func (r *ModelRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	models, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Object: models,
		Keys:   []string{"name", "description"},
	}, nil
}

// Describe implements the `describe model ....` command
func (r *ModelRunner) Describe(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	model, err := r.resource.GetByName(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Object:   model,
		Template: "Name: {{.Name}} ({{.Id}})",
	}, nil
}

/*
*******************************************************************************
Writer interface
*******************************************************************************
*/

// Create implements the `create model ...` command
func (r *ModelRunner) Create(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	options := in.Options.(*flags.ModelCreateOptions)

	if options.Replace {
		existing, err := r.resource.GetByName(name)

		if existing != nil {
			if err := r.resource.Delete(existing.Id, false); err != nil {
				return nil, err
			}
		} else if err != nil {
			if !errors.Is(err, resources.ErrNotFound) {
				return nil, err
			}
		}
	}

	model := services.NewModel(name, options.Description)

	if options.Schema != "" {
		data, err := os.ReadFile(options.Schema)
		if err != nil {
			return nil, err
		}

		var schema map[string]interface{}

		if err := json.Unmarshal(data, &schema); err != nil {
			return nil, err
		}

		model.Schema = schema
	}

	res, err := r.resource.Create(model)
	if err != nil {
		return nil, err
	}

	return &Response{
		Template: "Successfully created model `{{.Name}}`",
		Object:   res,
	}, nil
}

// Delete implements the `delete model ...` command
func (r *ModelRunner) Delete(in Request) (*Response, error) {
	logging.Trace()

	options := in.Options.(*flags.ModelDeleteOptions)
	name := in.Args[0]

	// Find the model by name
	model, err := r.resource.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("model `%s` not found", name)
	}

	// Delete with options (moved to resource layer)
	deleteOpts := resources.DeleteOptions{
		DeleteInstances: options.DeleteInstances,
		DeleteRelated:   options.All,
	}

	if err := r.resource.DeleteWithOptions(model, deleteOpts); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted model `%s`", name),
	}, nil
}

// Clear implements the `clear models` command
func (r *ModelRunner) Clear(in Request) (*Response, error) {
	logging.Trace()

	models, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range models {
		if err := r.resource.Delete(ele.Id, false); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v model(s)", len(models)),
	}, nil
}

/*
******************************************************************************
Copier interface
******************************************************************************
*/

// Copy implements the `copy model ...` command
func (r *ModelRunner) Copy(in Request) (*Response, error) {
	logging.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "model"}, r)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied model `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil
}

func (r *ModelRunner) CopyFrom(profile, name string) (any, error) {
	logging.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res := resources.NewModelResource(
		services.NewModelService(client),
		services.NewWorkflowService(client),
		services.NewTransformationService(client),
		services.NewInstanceService(client),
	)

	model, err := res.GetByName(name)
	if err != nil {
		return nil, err
	}

	exported, err := res.Export(model.Id)
	if err != nil {
		return nil, err
	}

	return *exported, err
}

func (r *ModelRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logging.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	resource := resources.NewModelResource(
		services.NewModelService(client),
		services.NewWorkflowService(client),
		services.NewTransformationService(client),
		services.NewInstanceService(client),
	)

	name := in.(services.Model).Name

	if exists, err := resource.GetByName(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("model `%s` exists on the destination server, use --replace to overwrite", name))
		} else if err != nil {
			return nil, err
		}
		if err := resource.Delete(exists.Id, false); err != nil {
			return nil, err
		}
	}

	res, err := resource.Import(in.(services.Model))
	if err != nil {
		return nil, err
	}

	return res, nil
}

/*
******************************************************************************
Importer interface
******************************************************************************
*/

// Import implements the `import model ...` command
func (r *ModelRunner) Import(in Request) (*Response, error) {
	logging.Trace()

	common := in.Common.(*flags.AssetImportCommon)
	options := in.Options.(*flags.ModelImportOptions)

	path, err := importGetPathFromRequest(in)
	if err != nil {
		return nil, err
	}

	wd := filepath.Dir(path)

	if common.Repository != "" {
		defer os.RemoveAll(wd)
	}

	var mModel map[string]interface{}

	if err := importLoadFromDisk(path, &mModel); err != nil {
		return nil, err
	}

	// Check the actions defined in the model to validate the model can be
	// imported.  This will also handle reconstructing the model defintion if
	// it was exported using `--expand`
	for _, ele := range mModel["actions"].([]interface{}) {
		if err := r.importActionMap(ele.(map[string]interface{}), wd, options.SkipChecks); err != nil {
			return nil, err
		}
	}

	b, err := json.Marshal(mModel)
	if err != nil {
		return nil, err
	}

	var model services.Model
	if err := json.Unmarshal(b, &model); err != nil {
		return nil, err
	}

	if common.Replace {
		m, err := r.resource.GetByName(model.Name)

		if err != nil {
			if !strings.HasSuffix(err.Error(), "not found") {
				return nil, err
			}
		}

		if m != nil {
			instances, err := r.resource.GetInstances(m.Id)
			if err != nil {
				return nil, err
			}

			if len(instances) > 0 {
				return nil, fmt.Errorf("cannot replace a model that has instances")
			}

			if err := r.resource.Delete(m.Id, false); err != nil {
				return nil, err
			}
		}
	}

	res, err := r.resource.Import(model)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully imported model `%s` (%s)", res.Name, res.Id),
		Object: model,
	}, nil
}

/*
******************************************************************************
Exporter interface
******************************************************************************
*/

// Export implements the `export model ...` command
func (r *ModelRunner) Export(in Request) (*Response, error) {
	logging.Trace()

	common := in.Common.(*flags.AssetExportCommon)
	options := in.Options.(*flags.ModelExportOptions)

	name := in.Args[0]

	res, err := r.resource.GetByName(name)
	if err != nil {
		return nil, err
	}

	model, err := r.resource.Export(res.Id)
	if err != nil {
		return nil, err
	}

	if options.Expand {
		path := common.Path

		var repo *Repository
		var repoPath string

		if common.Repository != "" {
			repo, err = exportNewRepositoryFromRequest(in)
			if err != nil {
				return nil, err
			}

			var e error

			repoPath, e = repo.Clone(
				&FileReaderImpl{},
				&ClonerImpl{},
			)
			if e != nil {
				return nil, e
			}
			defer os.RemoveAll(repoPath)

			path = filepath.Join(repoPath, common.Path)
		}

		if err := r.expandModel(in, model, path); err != nil {
			return nil, err
		}

		if common.Repository != "" {
			logging.Info("commiting %s to %s", repoPath, common.Repository)
			if err := repo.CommitAndPush(repoPath, common.Message); err != nil {
				return nil, err
			}
		}

	} else {
		fn := fmt.Sprintf("%s.model.json", normalizeFilename(model.Name))

		if err := exportAssetFromRequest(in, model, fn); err != nil {
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("Successfull exported model `%s`", model.Name),
	}, nil
}

/*
*******************************************************************************
Private functions
*******************************************************************************
*/

// expandModel will take a model and export the assets associated with the
// actions in the model.
func (r *ModelRunner) expandModel(in Request, model *services.Model, path string) error {
	logging.Trace()

	mModel, err := toMap(model)
	if err != nil {
		return err
	}

	for idx, ele := range model.Actions {

		var res any
		var e error
		var fn string

		mAction := mModel["actions"].([]interface{})[idx].(map[string]interface{})

		// export the action workflow
		if ele.Workflow != nil && *ele.Workflow != "" {
			res, e = services.NewWorkflowService(r.client).Export(*ele.Workflow)
			if e != nil {
				return e
			}

			name := normalizeFilename(res.(*services.Workflow).Name)
			fn = fmt.Sprintf("%s.workflow.json", name)

			delete(mAction, "workflow")
			mAction["workflowFilename"] = fn

			if err := utils.WriteJsonToDisk(res, fn, path); err != nil {
				return err
			}
		}

		// export the preworkflow jst
		if ele.PreWorkflowJst != nil && *ele.PreWorkflowJst != "" {
			res, e = services.NewTransformationService(r.client).Get(*ele.PreWorkflowJst)
			if e != nil {
				return e
			}
			name := normalizeFilename(res.(*services.Transformation).Name)
			fn = fmt.Sprintf("%s.transformation.json", name)

			delete(mAction, "preWorkflowJst")
			mAction["preWorkflowJstFilename"] = fn

			if err := utils.WriteJsonToDisk(res, fn, path); err != nil {
				return err
			}
		}

		// export the postworkflow jst
		if ele.PostWorkflowJst != nil && *ele.PostWorkflowJst != "" {
			res, e = services.NewTransformationService(r.client).Get(*ele.PostWorkflowJst)
			if e != nil {
				return e
			}
			name := normalizeFilename(res.(*services.Transformation).Name)
			fn = fmt.Sprintf("%s.transformation.json", name)

			delete(mAction, "postWorkflowJst")
			mAction["postWorkflowJstFilename"] = fn

			if err := utils.WriteJsonToDisk(res, fn, path); err != nil {
				return err
			}
		}
	}

	fn := fmt.Sprintf("%s.model.json", model.Name)

	if err := utils.WriteJsonToDisk(mModel, fn, path); err != nil {
		return err
	}

	return nil
}

func (r *ModelRunner) importActionMap(action map[string]interface{}, path string, skipChecks bool) error {
	logging.Trace()

	if value, exists := action["workflow"]; exists {
		if !skipChecks {
			wf, err := services.NewWorkflowService(r.client).Get(value.(string))
			if err != nil {
				return errors.New(
					fmt.Sprintf("workflow for action `%s` encountered the following error: %s", action["name"], err.Error()),
				)
			}
			if wf != nil {
				return errors.New(
					fmt.Sprintf("workflow for action `%s` does not exist", action["name"]),
				)
			}
		}

	} else if value, exists := action["workflowFilename"]; exists {
		var wf services.Workflow
		if err := importLoadFromDisk(filepath.Join(path, value.(string)), &wf); err != nil {
			return err
		}
		res, err := services.NewWorkflowService(r.client).Import(wf)
		if err != nil {
			return err
		}
		delete(action, "workflowFilename")
		action["workflow"] = res.Id

	} else if value, exists := action["preWorkflowJst"]; exists {
		if !skipChecks {
			var res *services.Transformation
			res, err := services.NewTransformationService(r.client).Get(value.(string))
			if err != nil {
				return errors.New(
					fmt.Sprintf("pre transformation for action `%s` encountered the following error: %s", action["name"], err.Error()),
				)
			}
			if res != nil {
				return errors.New(
					fmt.Sprintf("pre transformation for action `%s` does not exist", action["name"]),
				)
			}
		}

	} else if value, exists := action["preWorkflowJstFilename"]; exists {
		var jst services.Transformation
		if err := importLoadFromDisk(filepath.Join(path, value.(string)), &jst); err != nil {
			return err
		}
		res, err := services.NewTransformationService(r.client).Import(jst)
		if err != nil {
			return err
		}
		delete(action, "preWorkflowjstFilename")
		action["preWorkflowJst"] = res.Id

	} else if value, exists := action["postWorkflowJst"]; exists {
		if !skipChecks {
			var res *services.Transformation
			res, err := services.NewTransformationService(r.client).Get(value.(string))
			if err != nil {
				return errors.New(
					fmt.Sprintf("post transformation for action `%s` encountered the following error: %s", action["name"], err.Error()),
				)
			}
			if res != nil {
				return errors.New(
					fmt.Sprintf("pre transformation for action `%s` does not exist", action["name"]),
				)
			}
		}

	} else if value, exists := action["postWorkflowJstFilename"]; exists {
		var jst services.Transformation
		if err := importLoadFromDisk(filepath.Join(path, value.(string)), &jst); err != nil {
			return err
		}
		res, err := services.NewTransformationService(r.client).Import(jst)
		if err != nil {
			return err
		}
		delete(action, "postWorkflowjstFilename")
		action["postWorkflowJst"] = res.Id
	}

	return nil
}
