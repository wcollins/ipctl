// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/resources"
	"github.com/itential/ipctl/pkg/services"
)

const (
	jsonFormUrlTemplate = "/automation-studio/#/edit?tab=0&json-form=%s"
)

type JsonFormRunner struct {
	BaseRunner
	resource resources.JsonFormResourcer
}

func NewJsonFormRunner(c client.Client, cfg config.Provider) *JsonFormRunner {
	return &JsonFormRunner{
		resource:   resources.NewJsonFormResource(services.NewJsonFormService(c)),
		BaseRunner: NewBaseRunner(c, cfg),
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get implements the `get json_forms` command
func (r *JsonFormRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	var options flags.WorkflowGetOptions
	utils.LoadObject(in.Options, &options)

	res, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	var jsonforms []services.JsonForm

	for _, ele := range res {
		if strings.HasPrefix(ele.Name, "@") && options.All {
			jsonforms = append(jsonforms, ele)
		} else if !strings.HasPrefix(ele.Name, "@") {
			jsonforms = append(jsonforms, ele)
		}
	}

	return &Response{
		Keys:   []string{"name"},
		Object: jsonforms,
	}, nil

}

// Describe implements the `describe json_form <name>` command
func (r *JsonFormRunner) Describe(in Request) (*Response, error) {
	logging.Trace()

	res, err := r.resource.Get(in.Args[0])
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Name: %s (%s)", res.Name, res.Id),
		Object: res,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

// Create implements the `create jsonform <name>` command
func (r *JsonFormRunner) Create(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	options := in.Options.(*flags.JsonFormCreateOptions)

	if options.Replace {
		existing, err := r.resource.GetByName(name)

		if existing != nil {
			if err := r.resource.Delete([]string{existing.Id}); err != nil {
				return nil, err
			}
		} else if err != nil {
			if !errors.Is(err, resources.ErrNotFound) {
				return nil, err
			}
		}
	}

	jf, err := r.resource.Create(services.NewJsonForm(name, options.Description))
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully created jsonform `%s` (%s)", jf.Name, jf.Id),
		Object: jf,
	}, nil
}

// Delete implements the `delete jsonform <name>` command
func (r *JsonFormRunner) Delete(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	elements, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	var jf *services.JsonForm

	for _, ele := range elements {
		if ele.Name == name {
			jf = &ele
			break
		}
	}

	if jf == nil {
		return nil, errors.New(fmt.Sprintf("JSON form `%s` not found", name))
	}

	if err := r.resource.Delete([]string{jf.Id}); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted jsonform `%s`", name),
	}, nil
}

// Clear implements the `clear jsonforms` command
func (r *JsonFormRunner) Clear(in Request) (*Response, error) {
	logging.Trace()

	jsonforms, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	var ids []string

	for _, ele := range jsonforms {
		ids = append(ids, ele.Id)
	}

	if err := r.resource.Delete(ids); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v jsonform(s)", len(ids)),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

// Copy implements the `copy jsonform <name>` command
func (r *JsonFormRunner) Copy(in Request) (*Response, error) {
	logging.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "jsonform"}, r)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied jsonform `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil
}

func (r *JsonFormRunner) CopyFrom(profile, name string) (any, error) {
	logging.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewJsonFormService(client)
	res := resources.NewJsonFormResource(svc)

	jsonform, err := res.GetByName(name)
	if err != nil {
		return nil, err
	}

	return *jsonform, err
}

func (r *JsonFormRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logging.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	resource := resources.NewJsonFormResource(services.NewJsonFormService(client))

	name := in.(services.JsonForm).Name

	if exists, err := resource.GetByName(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("jsonform `%s` exists on the destination server, use --replace to overwrite", name))
		} else if err != nil {
			return nil, err
		}
	}

	return resource.Import(in.(services.JsonForm))
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import implements the command `import jsonform <path>`
func (r *JsonFormRunner) Import(in Request) (*Response, error) {
	logging.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var res services.JsonForm

	if err := importUnmarshalFromRequest(in, &res); err != nil {
		return nil, err
	}

	jf, err := r.importJsonForm(res, common.Replace)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully imported jsonform `%s` (%s)", jf.Name, jf.Id),
		Object: jf,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export is the implementation of the command `export jsonform <name>`
func (r *JsonFormRunner) Export(in Request) (*Response, error) {
	logging.Trace()

	var options *flags.AssetExportCommon
	utils.LoadObject(in.Common, &options)

	name := in.Args[0]

	jsonform, err := r.resource.GetByName(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.jsonform.json", name)

	if err := utils.WriteJsonToDisk(jsonform, fn, options.Path); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported jsonform `%s`", jsonform.Name),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

func (r JsonFormRunner) importJsonForm(in services.JsonForm, replace bool) (*services.JsonForm, error) {
	logging.Trace()

	p, err := r.resource.Get(in.Name)
	if err == nil {
		if replace {
			if err := r.resource.Delete([]string{p.Id}); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(fmt.Sprintf("jsonform with name `%s` already exists", p.Name))
		}
	}

	return r.resource.Import(in)
}
