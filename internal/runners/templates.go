// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/resources"
	"github.com/itential/ipctl/pkg/services"
)

type TemplateRunner struct {
	resource resources.TemplateResourcer
	BaseRunner
}

func NewTemplateRunner(c client.Client, cfg config.Provider) *TemplateRunner {
	return &TemplateRunner{
		resource:   resources.NewTemplateResource(services.NewTemplateService(c)),
		BaseRunner: NewBaseRunner(c, cfg),
	}
}

/*
*******************************************************************************
Reader interface
*******************************************************************************
*/

// Get implements the `get command-templates` command
func (r *TemplateRunner) Get(in Request) (*Response, error) {
	logging.Trace()

	options := in.Options.(*flags.TemplateGetOptions)

	res, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	var templates []services.Template

	for _, ele := range res {
		if strings.HasPrefix(ele.Name, "@") && options.All {
			templates = append(templates, ele)
		} else if !strings.HasPrefix(ele.Name, "@") {
			templates = append(templates, ele)
		}
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: templates,
	}, nil

}

// Describe implements the `describe command-template <name>` command
func (r *TemplateRunner) Describe(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	res, err := r.resource.GetByName(name)
	if err != nil {
		return nil, err
	}

	output := []string{
		fmt.Sprintf("Name: %s (%s)", res.Name, res.Id),
		fmt.Sprintf("Description: %s", res.Description),
		fmt.Sprintf("Type: %s", res.Type),
		fmt.Sprintf("Group: %s, Command: %s", res.Group, res.Command),
		fmt.Sprintf("Created: %s", res.Created),
		fmt.Sprintf("Updated: %s", res.LastUpdated),
	}

	return &Response{
		Text:   strings.Join(output, "\n"),
		Object: res,
	}, nil
}

/*
*******************************************************************************
Writer interface
*******************************************************************************
*/

// Create implements the `create template ...` command
func (r *TemplateRunner) Create(in Request) (*Response, error) {
	logging.Trace()

	options := in.Options.(*flags.TemplateCreateOptions)

	name := in.Args[0]

	if options.Replace {
		existing, err := r.resource.GetByName(name)

		if existing != nil {
			if err := r.resource.Delete(existing.Id); err != nil {
				return nil, err
			}
		} else if err != nil {
			if !errors.Is(err, resources.ErrNotFound) {
				return nil, err
			}
		}
	}

	res, err := r.resource.Create(services.NewTemplate(
		name,
		options.Group,
		options.Description,
		options.Type,
	))
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   fmt.Sprintf("Successfully created template `%s` (%s)", res.Name, res.Id),
		Object: res,
	}, nil
}

func (r *TemplateRunner) Delete(in Request) (*Response, error) {
	logging.Trace()

	t, err := r.resource.GetByName(in.Args[0])
	if err != nil {
		return nil, err
	}

	if err := r.resource.Delete(t.Id); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted template `%s` (%s)", t.Name, t.Id),
	}, nil
}

func (r *TemplateRunner) Clear(in Request) (*Response, error) {
	logging.Trace()

	elements, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	for _, ele := range elements {
		terminal.Display("Deleting template `%s`  (%s)", ele.Name, ele.Id)
		if err := r.resource.Delete(ele.Id); err != nil {
			logging.Debug("failed to delete template `%s` (%s)", ele.Name, ele.Id)
			return nil, err
		}
	}

	return &Response{
		Text: fmt.Sprintf("\nDeleted %v template(s)", len(elements)),
	}, nil
}

/*
*******************************************************************************
Copier interface
*******************************************************************************
*/

func (r *TemplateRunner) Copy(in Request) (*Response, error) {
	logging.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "template"}, r)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied template `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil
}

func (r *TemplateRunner) CopyFrom(profile, name string) (any, error) {
	logging.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewTemplateService(client)

	template, err := svc.Export(name)
	if err != nil {
		return nil, err
	}
	return *template, nil

}

func (r *TemplateRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logging.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewTemplateService(client)

	name := in.(services.Template).Name

	if exists, err := svc.Get(name); exists != nil {
		if !replace {
			return nil, errors.New(fmt.Sprintf("template `%s` exists on the destination server", name))
		} else if err != nil {
			return nil, err
		}
		logging.Info("Deleting existing template `%s` from `%s`", name, profile)
		if err := svc.Delete(name); err != nil {
			return nil, err
		}
	}

	res, err := svc.Import(in.(services.Template))
	if err != nil {
		return nil, err
	}

	return res, nil

}

/*
*******************************************************************************
Importer interface
*******************************************************************************
*/

func (r *TemplateRunner) Import(in Request) (*Response, error) {
	logging.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var res services.Template

	if err := importUnmarshalFromRequest(in, &res); err != nil {
		return nil, err
	}

	if err := r.importTemplate(res, common.Replace); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported command template `%s`", res.Name),
	}, nil
}

/*
*******************************************************************************
Exporter interface
*******************************************************************************
*/

func (r *TemplateRunner) Export(in Request) (*Response, error) {
	logging.Trace()

	name := in.Args[0]

	res, err := r.resource.GetByName(name)
	if err != nil {
		return nil, err
	}

	exported, err := r.resource.Export(res.Id)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.template.json", exported.Name)

	if err := exportAssetFromRequest(in, exported, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported template `%s` (%s)", exported.Name, exported.Id),
	}, nil
}

/*
*******************************************************************************
Dumper interface
*******************************************************************************
*/

// Dump implements the `dump templates...` command
func (r *TemplateRunner) Dump(in Request) (*Response, error) {
	logging.Trace()

	res, err := r.resource.GetAll()
	if err != nil {
		return nil, err
	}

	var assets = map[string]interface{}{}

	for _, ele := range res {
		if !strings.HasPrefix(ele.Name, "@") {
			key := fmt.Sprintf("%s.template.json", ele.Name)
			assets[key] = ele
		}
	}

	if err := dumpAssets(in, assets); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Dumped %v template(s)", len(assets)),
	}, nil
}

/*
*******************************************************************************
Loader interface
*******************************************************************************
*/

// Load implements the `load template ...` command
func (r *TemplateRunner) Load(in Request) (*Response, error) {
	logging.Trace()

	options := in.Options.(*flags.TemplateLoadOptions)

	var elements map[string]interface{}
	var err error

	if options.Type == "textfsm" {
		elements, err = loadStringAssets(in, LoadOptions{})
	} else {
		elements, err = loadAssets(in)
		if err != nil {
			return nil, err
		}
	}

	var loaded int
	var skipped int

	for fn, ele := range elements {
		var template services.Template
		var err error

		if options.Type == "textfsm" || options.Type == "jinja2" {
			name := strings.TrimSuffix(fn, filepath.Ext(fn))
			template = services.NewTemplate(name, "Imported", "", options.Type)
			template.Template = ele.(string)
		} else {
			err = loadUnmarshalAsset(ele, &template)
			if err != nil {
				terminal.Display("Failed to load template from `%s`, skipping", fn)
				skipped++
			}
		}

		if err == nil {
			if err := r.importTemplate(template, false); err != nil {
				if !strings.HasPrefix(err.Error(), "template with name") {
					return nil, err
				}
				terminal.Display("Skipping `%s`, template `%s` already exists", fn, template.Name)
				skipped++
			} else {
				terminal.Display("Loaded template `%s` successfully from `%s`", template.Name, fn)
				loaded++
			}
		}
	}

	output := fmt.Sprintf("\nSuccessfully loaded %v and skipped %v files from `%s`", loaded, skipped, in.Args[0])

	return &Response{
		Text: output,
	}, nil

}

/*
*******************************************************************************
Private functions
*******************************************************************************
*/

func (r TemplateRunner) importTemplate(in services.Template, replace bool) error {
	logging.Trace()
	logging.Debug("attempting to import template `%s`", in.Name)

	p, err := r.resource.GetByName(in.Name)
	if err != nil {
		if !errors.Is(err, resources.ErrNotFound) {
			return err
		}
	}
	if p != nil {
		if replace {
			logging.Debug("template exists, deleting it")
			r.resource.Delete(p.Id)
		} else {
			return errors.New(fmt.Sprintf("template with name `%s` already exists, use `--replace` to overwrite", p.Name))
		}
	}

	_, err = r.resource.Import(in)
	if err != nil {
		return err
	}

	logging.Debug("successfully imported template `%s`", in.Name)

	return nil
}
