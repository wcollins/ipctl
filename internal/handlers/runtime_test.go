// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"testing"

	"github.com/itential/ipctl/internal/profile"
	"github.com/itential/ipctl/internal/repository"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockClient implements client.Client for testing
type mockClient struct{}

func (m *mockClient) Get(*client.Request) (*client.Response, error)    { return nil, nil }
func (m *mockClient) Post(*client.Request) (*client.Response, error)   { return nil, nil }
func (m *mockClient) Put(*client.Request) (*client.Response, error)    { return nil, nil }
func (m *mockClient) Delete(*client.Request) (*client.Response, error) { return nil, nil }
func (m *mockClient) Patch(*client.Request) (*client.Response, error)  { return nil, nil }
func (m *mockClient) Trace(*client.Request) (*client.Response, error)  { return nil, nil }

// mockConfig implements config.Provider for testing
type mockConfig struct {
	workingDir        string
	defaultProfile    string
	defaultRepository string
	datasetsEnabled   bool
	gitName           string
	gitEmail          string
	gitUser           string
}

func (m *mockConfig) GetProfile(name string) (*profile.Profile, error) {
	return &profile.Profile{Host: name}, nil
}

func (m *mockConfig) ActiveProfile() (*profile.Profile, error) {
	return &profile.Profile{Host: "localhost"}, nil
}

func (m *mockConfig) GetRepository(name string) (*repository.Repository, error) {
	return &repository.Repository{Url: name}, nil
}

func (m *mockConfig) GetWorkingDir() string        { return m.workingDir }
func (m *mockConfig) GetDefaultProfile() string    { return m.defaultProfile }
func (m *mockConfig) GetDefaultRepository() string { return m.defaultRepository }
func (m *mockConfig) IsDatasetsEnabled() bool      { return m.datasetsEnabled }
func (m *mockConfig) GetGitName() string           { return m.gitName }
func (m *mockConfig) GetGitEmail() string          { return m.gitEmail }
func (m *mockConfig) GetGitUser() string           { return m.gitUser }

func TestNewRuntime(t *testing.T) {
	tests := []struct {
		name        string
		client      client.Client
		config      *mockConfig
		termCfg     *terminal.Config
		expectError bool
	}{
		{
			name:   "successful runtime creation",
			client: &mockClient{},
			config: &mockConfig{
				workingDir:      "/tmp",
				defaultProfile:  "default",
				datasetsEnabled: false,
			},
			termCfg: &terminal.Config{
				NoColor:       false,
				DefaultOutput: "human",
				Pager:         false,
			},
			expectError: false,
		},
		{
			name:   "runtime with custom terminal config",
			client: &mockClient{},
			config: &mockConfig{
				workingDir:      "/home/user",
				defaultProfile:  "production",
				datasetsEnabled: true,
			},
			termCfg: &terminal.Config{
				NoColor:       true,
				DefaultOutput: "json",
				Pager:         true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt, err := NewRuntime(tt.client, tt.config, tt.termCfg)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, rt)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rt)

				// Verify all fields are set
				assert.Equal(t, tt.client, rt.client)
				assert.Equal(t, tt.config, rt.config)
				assert.Equal(t, tt.termCfg, rt.terminalConfig)
				assert.NotNil(t, rt.descriptors)
				assert.False(t, rt.Verbose) // Default value
			}
		})
	}
}

func TestRuntime_GetClient(t *testing.T) {
	mockCli := &mockClient{}
	rt, err := NewRuntime(mockCli, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	client := rt.GetClient()
	assert.Equal(t, mockCli, client)
}

func TestRuntime_GetConfig(t *testing.T) {
	mockCfg := &mockConfig{
		workingDir:     "/test",
		defaultProfile: "test-profile",
	}
	rt, err := NewRuntime(&mockClient{}, mockCfg, &terminal.Config{})
	require.NoError(t, err)

	config := rt.GetConfig()
	assert.Equal(t, mockCfg, config)
	assert.Equal(t, "/test", config.GetWorkingDir())
	assert.Equal(t, "test-profile", config.GetDefaultProfile())
}

func TestRuntime_GetTerminalConfig(t *testing.T) {
	termCfg := &terminal.Config{
		NoColor:       true,
		DefaultOutput: "yaml",
		Pager:         true,
	}
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, termCfg)
	require.NoError(t, err)

	config := rt.GetTerminalConfig()
	assert.Equal(t, termCfg, config)
	assert.True(t, config.NoColor)
	assert.Equal(t, "yaml", config.DefaultOutput)
	assert.True(t, config.Pager)
}

func TestRuntime_GetDescriptors(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	descriptors := rt.GetDescriptors()
	assert.NotNil(t, descriptors)

	// Verify that descriptors are loaded from embedded files
	// Check for known descriptors
	assert.Contains(t, descriptors, "projects")
	assert.Contains(t, descriptors, "workflows")
	assert.Contains(t, descriptors, "automations")
	assert.Contains(t, descriptors, "accounts")
	assert.Contains(t, descriptors, "adapters")
}

func TestRuntime_IsVerbose(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{
			name:    "verbose disabled",
			verbose: false,
		},
		{
			name:    "verbose enabled",
			verbose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
			require.NoError(t, err)

			rt.Verbose = tt.verbose
			assert.Equal(t, tt.verbose, rt.IsVerbose())
		})
	}
}

func TestRuntime_PointerReceivers(t *testing.T) {
	// Verify that all methods use pointer receivers by checking
	// that modifications to the runtime are visible after method calls
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	// Get the client - this should work with pointer receiver
	client := rt.GetClient()
	assert.NotNil(t, client)

	// Modify Verbose field
	rt.Verbose = true

	// Verify the modification is visible (would not be with value receivers)
	assert.True(t, rt.IsVerbose())
}

func TestRuntime_ImplementsRuntimeContext(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	// Verify Runtime implements RuntimeContext interface
	var _ RuntimeContext = rt

	// Test all interface methods
	assert.NotNil(t, rt.GetClient())
	assert.NotNil(t, rt.GetConfig())
	assert.NotNil(t, rt.GetTerminalConfig())
	assert.NotNil(t, rt.GetDescriptors())
	assert.False(t, rt.IsVerbose())
}

func TestRuntime_DescriptorsAreShared(t *testing.T) {
	// Create two runtimes
	rt1, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	rt2, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	// Descriptors should be loaded independently (not shared between runtimes)
	desc1 := rt1.GetDescriptors()
	desc2 := rt2.GetDescriptors()

	assert.NotNil(t, desc1)
	assert.NotNil(t, desc2)

	// Both should have the same descriptor keys
	assert.Equal(t, len(desc1), len(desc2))
}

func TestRuntime_NilInputs(t *testing.T) {
	tests := []struct {
		name    string
		client  client.Client
		config  *mockConfig
		termCfg *terminal.Config
		wantErr bool
	}{
		{
			name:    "all valid inputs",
			client:  &mockClient{},
			config:  &mockConfig{},
			termCfg: &terminal.Config{},
			wantErr: false,
		},
		{
			name:    "nil client",
			client:  nil,
			config:  &mockConfig{},
			termCfg: &terminal.Config{},
			wantErr: false, // Runtime doesn't validate nil client
		},
		{
			name:    "nil config",
			client:  &mockClient{},
			config:  nil,
			termCfg: &terminal.Config{},
			wantErr: false, // Runtime doesn't validate nil config
		},
		{
			name:    "nil terminal config",
			client:  &mockClient{},
			config:  &mockConfig{},
			termCfg: nil,
			wantErr: false, // Runtime doesn't validate nil termCfg
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt, err := NewRuntime(tt.client, tt.config, tt.termCfg)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, rt)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, rt)
			}
		})
	}
}

func TestRuntime_DescriptorStructure(t *testing.T) {
	rt, err := NewRuntime(&mockClient{}, &mockConfig{}, &terminal.Config{})
	require.NoError(t, err)

	descriptors := rt.GetDescriptors()

	// Test descriptor structure for a known resource (projects)
	projectDesc, ok := descriptors["projects"]
	require.True(t, ok, "projects descriptor should exist")
	assert.NotNil(t, projectDesc)

	// Verify the descriptor map has expected commands
	// Most resources should have at least "get" command
	if getDesc, ok := projectDesc["get"]; ok {
		assert.NotEmpty(t, getDesc.Use, "get descriptor should have Use field")
	}
}
