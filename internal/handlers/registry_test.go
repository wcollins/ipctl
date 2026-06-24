// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockReader implements the Reader interface for testing
type mockReader struct {
	name string
}

func (m *mockReader) Get(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-get"}
}

func (m *mockReader) Describe(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-describe"}
}

// mockWriter implements the Writer interface for testing
type mockWriter struct {
	name string
}

func (m *mockWriter) Create(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-create"}
}

func (m *mockWriter) Delete(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-delete"}
}

func (m *mockWriter) Clear(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-clear"}
}

// mockCopier implements the Copier interface for testing
type mockCopier struct {
	name string
}

func (m *mockCopier) Copy(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-copy"}
}

// mockController implements the Controller interface for testing
type mockController struct {
	name string
}

func (m *mockController) Start(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-start"}
}

func (m *mockController) Stop(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-stop"}
}

func (m *mockController) Restart(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-restart"}
}

// mockEditor implements the Editor interface for testing
type mockEditor struct {
	name string
}

func (m *mockEditor) Edit(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-edit"}
}

// mockImporter implements the Importer interface for testing
type mockImporter struct {
	name string
}

func (m *mockImporter) Import(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-import"}
}

// mockExporter implements the Exporter interface for testing
type mockExporter struct {
	name string
}

func (m *mockExporter) Export(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-export"}
}

// mockInspector implements the Inspector interface for testing
type mockInspector struct {
	name string
}

func (m *mockInspector) Inspect(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-inspect"}
}

// mockDumper implements the Dumper interface for testing
type mockDumper struct {
	name string
}

func (m *mockDumper) Dump(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-dump"}
}

// mockLoader implements the Loader interface for testing
type mockLoader struct {
	name string
}

func (m *mockLoader) Load(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-load"}
}

// mockMultiHandler implements multiple interfaces
type mockMultiHandler struct {
	name string
}

func (m *mockMultiHandler) Get(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-get"}
}

func (m *mockMultiHandler) Describe(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-describe"}
}

func (m *mockMultiHandler) Create(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-create"}
}

func (m *mockMultiHandler) Delete(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-delete"}
}

func (m *mockMultiHandler) Clear(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-clear"}
}

func (m *mockMultiHandler) Copy(rt *Runtime) *cobra.Command {
	return &cobra.Command{Use: m.name + "-copy"}
}

func TestNewRegistry_EmptyHandlers(t *testing.T) {
	registry := NewRegistry([]any{})

	assert.NotNil(t, registry)
	assert.Empty(t, registry.Readers())
	assert.Empty(t, registry.Writers())
	assert.Empty(t, registry.Copiers())
	assert.Empty(t, registry.Controllers())
	assert.Empty(t, registry.Editors())
	assert.Empty(t, registry.Importers())
	assert.Empty(t, registry.Exporters())
	assert.Empty(t, registry.Inspectors())
	assert.Empty(t, registry.Dumpers())
	assert.Empty(t, registry.Loaders())
}

func TestNewRegistry_SingleInterface(t *testing.T) {
	tests := []struct {
		name    string
		handler any
		checkFn func(*testing.T, *Registry)
	}{
		{
			name:    "Reader interface",
			handler: &mockReader{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Readers(), 1)
				assert.Empty(t, r.Writers())
			},
		},
		{
			name:    "Writer interface",
			handler: &mockWriter{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Writers(), 1)
				assert.Empty(t, r.Readers())
			},
		},
		{
			name:    "Copier interface",
			handler: &mockCopier{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Copiers(), 1)
				assert.Empty(t, r.Readers())
			},
		},
		{
			name:    "Controller interface",
			handler: &mockController{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Controllers(), 1)
				assert.Empty(t, r.Readers())
			},
		},
		{
			name:    "Editor interface",
			handler: &mockEditor{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Editors(), 1)
				assert.Empty(t, r.Readers())
			},
		},
		{
			name:    "Importer interface",
			handler: &mockImporter{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Importers(), 1)
				assert.Empty(t, r.Readers())
			},
		},
		{
			name:    "Exporter interface",
			handler: &mockExporter{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Exporters(), 1)
				assert.Empty(t, r.Readers())
			},
		},
		{
			name:    "Inspector interface",
			handler: &mockInspector{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Inspectors(), 1)
				assert.Empty(t, r.Readers())
			},
		},
		{
			name:    "Dumper interface",
			handler: &mockDumper{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Dumpers(), 1)
				assert.Empty(t, r.Readers())
			},
		},
		{
			name:    "Loader interface",
			handler: &mockLoader{name: "test"},
			checkFn: func(t *testing.T, r *Registry) {
				assert.Len(t, r.Loaders(), 1)
				assert.Empty(t, r.Readers())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewRegistry([]any{tt.handler})
			require.NotNil(t, registry)
			tt.checkFn(t, registry)
		})
	}
}

func TestNewRegistry_MultipleInterfaces(t *testing.T) {
	handler := &mockMultiHandler{name: "multi"}
	registry := NewRegistry([]any{handler})

	require.NotNil(t, registry)

	// Handler implements Reader, Writer, and Copier
	assert.Len(t, registry.Readers(), 1)
	assert.Len(t, registry.Writers(), 1)
	assert.Len(t, registry.Copiers(), 1)

	// Handler does not implement other interfaces
	assert.Empty(t, registry.Controllers())
	assert.Empty(t, registry.Editors())
	assert.Empty(t, registry.Importers())
	assert.Empty(t, registry.Exporters())
	assert.Empty(t, registry.Inspectors())
	assert.Empty(t, registry.Dumpers())
	assert.Empty(t, registry.Loaders())
}

func TestNewRegistry_MultipleHandlers(t *testing.T) {
	handlers := []any{
		&mockReader{name: "handler1"},
		&mockReader{name: "handler2"},
		&mockWriter{name: "handler3"},
		&mockCopier{name: "handler4"},
	}

	registry := NewRegistry(handlers)
	require.NotNil(t, registry)

	assert.Len(t, registry.Readers(), 2)
	assert.Len(t, registry.Writers(), 1)
	assert.Len(t, registry.Copiers(), 1)
}

func TestNewRegistry_NonHandlerTypes(t *testing.T) {
	// Test that non-handler types are safely ignored
	handlers := []any{
		&mockReader{name: "valid"},
		"string",   // Not a handler
		42,         // Not a handler
		nil,        // Nil value
		struct{}{}, // Empty struct
	}

	registry := NewRegistry(handlers)
	require.NotNil(t, registry)

	// Only the valid handler should be registered
	assert.Len(t, registry.Readers(), 1)
}

func TestRegistry_DefensiveCopies(t *testing.T) {
	handler := &mockReader{name: "test"}
	registry := NewRegistry([]any{handler})

	// Get the readers slice
	readers1 := registry.Readers()
	assert.Len(t, readers1, 1)

	// Modify the returned slice
	readers1 = append(readers1, &mockReader{name: "injected"})

	// Get readers again - should not include the injected handler
	readers2 := registry.Readers()
	assert.Len(t, readers2, 1, "registry should return defensive copy")
}

func TestRegistry_AllInterfaces(t *testing.T) {
	// Create handlers for all 10 interfaces
	handlers := []any{
		&mockReader{name: "reader"},
		&mockWriter{name: "writer"},
		&mockCopier{name: "copier"},
		&mockController{name: "controller"},
		&mockEditor{name: "editor"},
		&mockImporter{name: "importer"},
		&mockExporter{name: "exporter"},
		&mockInspector{name: "inspector"},
		&mockDumper{name: "dumper"},
		&mockLoader{name: "loader"},
	}

	registry := NewRegistry(handlers)
	require.NotNil(t, registry)

	// Verify all interfaces are registered
	assert.Len(t, registry.Readers(), 1)
	assert.Len(t, registry.Writers(), 1)
	assert.Len(t, registry.Copiers(), 1)
	assert.Len(t, registry.Controllers(), 1)
	assert.Len(t, registry.Editors(), 1)
	assert.Len(t, registry.Importers(), 1)
	assert.Len(t, registry.Exporters(), 1)
	assert.Len(t, registry.Inspectors(), 1)
	assert.Len(t, registry.Dumpers(), 1)
	assert.Len(t, registry.Loaders(), 1)
}

func TestRegistry_TypeAssertionSafety(t *testing.T) {
	// Test that type assertions don't panic with various input types
	handlers := []any{
		nil,
		42,
		"string",
		[]string{"slice"},
		map[string]string{"map": "value"},
		func() {},
		&mockReader{name: "valid"},
	}

	// This should not panic
	registry := NewRegistry(handlers)
	require.NotNil(t, registry)

	// Only the valid handler should be registered
	assert.Len(t, registry.Readers(), 1)
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	handler := &mockReader{name: "test"}
	registry := NewRegistry([]any{handler})

	// Test concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			readers := registry.Readers()
			assert.Len(t, readers, 1)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestRegistry_SameHandlerMultipleTimes(t *testing.T) {
	handler := &mockReader{name: "test"}

	// Register the same handler multiple times
	handlers := []any{handler, handler, handler}
	registry := NewRegistry(handlers)

	// All three instances should be registered
	readers := registry.Readers()
	assert.Len(t, readers, 3)
}
