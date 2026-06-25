// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"testing"

	"github.com/itential/ipctl/internal/cmdutils"
	"github.com/itential/ipctl/internal/runners"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockFlagger implements flags.Flagger for testing
type mockFlagger struct{}

func (m *mockFlagger) Flags(cmd *cobra.Command) {}

// mockAssetRunner implements multiple runner interfaces for testing
type mockAssetRunner struct {
	supportsReader     bool
	supportsWriter     bool
	supportsCopier     bool
	supportsEditor     bool
	supportsImporter   bool
	supportsExporter   bool
	supportsController bool
	supportsInspector  bool
	supportsDumper     bool
	supportsLoader     bool
}

// Implement runners.Reader
func (m *mockAssetRunner) Get(req runners.Request) (*runners.Response, error) {
	if !m.supportsReader {
		return nil, nil
	}
	return &runners.Response{Text: "get"}, nil
}

func (m *mockAssetRunner) Describe(req runners.Request) (*runners.Response, error) {
	if !m.supportsReader {
		return nil, nil
	}
	return &runners.Response{Text: "describe"}, nil
}

// Implement runners.Writer
func (m *mockAssetRunner) Create(req runners.Request) (*runners.Response, error) {
	if !m.supportsWriter {
		return nil, nil
	}
	return &runners.Response{Text: "create"}, nil
}

func (m *mockAssetRunner) Delete(req runners.Request) (*runners.Response, error) {
	if !m.supportsWriter {
		return nil, nil
	}
	return &runners.Response{Text: "delete"}, nil
}

func (m *mockAssetRunner) Clear(req runners.Request) (*runners.Response, error) {
	if !m.supportsWriter {
		return nil, nil
	}
	return &runners.Response{Text: "clear"}, nil
}

// Implement runners.Copier
func (m *mockAssetRunner) Copy(req runners.Request) (*runners.Response, error) {
	if !m.supportsCopier {
		return nil, nil
	}
	return &runners.Response{Text: "copy"}, nil
}

// Implement runners.Editor
func (m *mockAssetRunner) Edit(req runners.Request) (*runners.Response, error) {
	if !m.supportsEditor {
		return nil, nil
	}
	return &runners.Response{Text: "edit"}, nil
}

// Implement runners.Importer
func (m *mockAssetRunner) Import(req runners.Request) (*runners.Response, error) {
	if !m.supportsImporter {
		return nil, nil
	}
	return &runners.Response{Text: "import"}, nil
}

// Implement runners.Exporter
func (m *mockAssetRunner) Export(req runners.Request) (*runners.Response, error) {
	if !m.supportsExporter {
		return nil, nil
	}
	return &runners.Response{Text: "export"}, nil
}

// Implement runners.Controller
func (m *mockAssetRunner) Start(req runners.Request) (*runners.Response, error) {
	if !m.supportsController {
		return nil, nil
	}
	return &runners.Response{Text: "start"}, nil
}

func (m *mockAssetRunner) Stop(req runners.Request) (*runners.Response, error) {
	if !m.supportsController {
		return nil, nil
	}
	return &runners.Response{Text: "stop"}, nil
}

func (m *mockAssetRunner) Restart(req runners.Request) (*runners.Response, error) {
	if !m.supportsController {
		return nil, nil
	}
	return &runners.Response{Text: "restart"}, nil
}

// Implement runners.Inspector
func (m *mockAssetRunner) Inspect(req runners.Request) (*runners.Response, error) {
	if !m.supportsInspector {
		return nil, nil
	}
	return &runners.Response{Text: "inspect"}, nil
}

// Implement runners.Dumper
func (m *mockAssetRunner) Dump(req runners.Request) (*runners.Response, error) {
	if !m.supportsDumper {
		return nil, nil
	}
	return &runners.Response{Text: "dump"}, nil
}

// Implement runners.Loader
func (m *mockAssetRunner) Load(req runners.Request) (*runners.Response, error) {
	if !m.supportsLoader {
		return nil, nil
	}
	return &runners.Response{Text: "load"}, nil
}

func createTestDescriptors() DescriptorMap {
	return DescriptorMap{
		"get": cmdutils.Descriptor{
			Use:         "resources",
			Description: "get resources",
		},
		"describe": cmdutils.Descriptor{
			Use:         "resource",
			Description: "describe resource",
		},
		"create": cmdutils.Descriptor{
			Use:         "resource",
			Description: "create resource",
		},
		"delete": cmdutils.Descriptor{
			Use:         "resource",
			Description: "delete resource",
		},
		"clear": cmdutils.Descriptor{
			Use:         "resources",
			Description: "clear resources",
		},
		"copy": cmdutils.Descriptor{
			Use:         "resource",
			Description: "copy resource",
		},
		"edit": cmdutils.Descriptor{
			Use:         "resource",
			Description: "edit resource",
		},
		"import": cmdutils.Descriptor{
			Use:         "resources",
			Description: "import resources",
		},
		"export": cmdutils.Descriptor{
			Use:         "resources",
			Description: "export resources",
		},
		"start": cmdutils.Descriptor{
			Use:         "resource",
			Description: "start resource",
		},
		"stop": cmdutils.Descriptor{
			Use:         "resource",
			Description: "stop resource",
		},
		"restart": cmdutils.Descriptor{
			Use:         "resource",
			Description: "restart resource",
		},
		"inspect": cmdutils.Descriptor{
			Use:         "resource",
			Description: "inspect resource",
		},
		"dump": cmdutils.Descriptor{
			Use:         "resources",
			Description: "dump resources",
		},
		"load": cmdutils.Descriptor{
			Use:         "resources",
			Description: "load resources",
		},
	}
}

func TestNewAssetHandler_WithFlags(t *testing.T) {
	runner := &mockAssetRunner{supportsReader: true}
	desc := createTestDescriptors()

	flags := &AssetHandlerFlags{
		Get: &mockFlagger{},
	}

	handler := NewAssetHandler(runner, desc, flags)

	assert.NotNil(t, handler.flags)
	assert.Equal(t, flags, handler.flags)
}

func TestNewAssetHandler_NilFlags(t *testing.T) {
	runner := &mockAssetRunner{supportsReader: true}
	desc := createTestDescriptors()

	handler := NewAssetHandler(runner, desc, nil)

	// Should create empty flags if nil is passed
	assert.NotNil(t, handler.flags)
}

func TestAssetHandler_Get_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsReader: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Get(rt)

	assert.NotNil(t, cmd)
	assert.Equal(t, "resources", cmd.Use)
}

func TestAssetHandler_Describe_ExactArgsSet(t *testing.T) {
	runner := &mockAssetRunner{
		supportsReader: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Describe(rt)

	require.NotNil(t, cmd)
	assert.NotNil(t, cmd.Args)

	// Verify exact args validator is set
	err = cmd.Args(cmd, []string{"arg1"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"arg1", "arg2"})
	assert.Error(t, err)
}

func TestAssetHandler_Create_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsWriter: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Create(rt)

	assert.NotNil(t, cmd)
	assert.Equal(t, "resource", cmd.Use)
}

func TestAssetHandler_Delete_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsWriter: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Delete(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Clear_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsWriter: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Clear(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Edit_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsEditor: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Edit(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Import_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsImporter: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Import(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Export_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsExporter: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Export(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Start_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsController: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Start(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Stop_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsController: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Stop(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Restart_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsController: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Restart(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Inspect_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsInspector: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Inspect(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Dump_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsDumper: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Dump(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_Load_WithSupport(t *testing.T) {
	runner := &mockAssetRunner{
		supportsLoader: true,
	}

	desc := createTestDescriptors()
	handler := NewAssetHandler(runner, desc, nil)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	cmd := handler.Load(rt)

	assert.NotNil(t, cmd)
}

func TestAssetHandler_WithCustomFlags(t *testing.T) {
	runner := &mockAssetRunner{
		supportsReader: true,
		supportsWriter: true,
	}

	desc := createTestDescriptors()

	customFlags := &AssetHandlerFlags{
		Get:    &mockFlagger{},
		Create: &mockFlagger{},
	}

	handler := NewAssetHandler(runner, desc, customFlags)

	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	// Get command should have flags applied
	getCmd := handler.Get(rt)
	assert.NotNil(t, getCmd)

	// Create command should have flags applied
	createCmd := handler.Create(rt)
	assert.NotNil(t, createCmd)
}

func TestAssetHandlerFlags_AllFields(t *testing.T) {
	flags := &AssetHandlerFlags{
		Create:   &mockFlagger{},
		Delete:   &mockFlagger{},
		Get:      &mockFlagger{},
		Describe: &mockFlagger{},
		Copy:     &mockFlagger{},
		Clear:    &mockFlagger{},
		Edit:     &mockFlagger{},
		Import:   &mockFlagger{},
		Export:   &mockFlagger{},
		Start:    &mockFlagger{},
		Stop:     &mockFlagger{},
		Restart:  &mockFlagger{},
		Inspect:  &mockFlagger{},
		Dump:     &mockFlagger{},
		Load:     &mockFlagger{},
	}

	assert.NotNil(t, flags.Create)
	assert.NotNil(t, flags.Delete)
	assert.NotNil(t, flags.Get)
	assert.NotNil(t, flags.Describe)
	assert.NotNil(t, flags.Copy)
	assert.NotNil(t, flags.Clear)
	assert.NotNil(t, flags.Edit)
	assert.NotNil(t, flags.Import)
	assert.NotNil(t, flags.Export)
	assert.NotNil(t, flags.Start)
	assert.NotNil(t, flags.Stop)
	assert.NotNil(t, flags.Restart)
	assert.NotNil(t, flags.Inspect)
	assert.NotNil(t, flags.Dump)
	assert.NotNil(t, flags.Load)
}
