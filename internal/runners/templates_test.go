// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"testing"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	templatesGetAllEmpty    = testlib.Fixture("testdata/templates/getall_empty.json")
	templatesImportResponse = testlib.Fixture("testdata/templates/import_response.json")
)

// TestImportTemplateReplaceNotExists covers the bug where importing with --replace
// failed with "item with name '...' not found" when the template did not yet exist
// on the server. The fix uses errors.Is(err, ErrNotFound) so a missing template is
// treated as a no-op and the import proceeds.
func TestImportTemplateReplaceNotExists(t *testing.T) {
	runner := NewTemplateRunner(
		testlib.Setup(),
		testlib.DefaultConfig(),
	)
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/automation-studio/templates", templatesGetAllEmpty, 0)
	testlib.AddPostResponseToMux("/automation-studio/templates/import", templatesImportResponse, 200)

	res, err := runner.Import(Request{
		Args: []string{"testdata/templates/template.json"},
		Common: &flags.AssetImportCommon{
			Replace: true,
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Text, "test-template")
}

// TestImportTemplateNoReplace verifies a plain import (no --replace) succeeds
// when the template does not exist on the server.
func TestImportTemplateNoReplace(t *testing.T) {
	runner := NewTemplateRunner(
		testlib.Setup(),
		testlib.DefaultConfig(),
	)
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/automation-studio/templates", templatesGetAllEmpty, 0)
	testlib.AddPostResponseToMux("/automation-studio/templates/import", templatesImportResponse, 200)

	res, err := runner.Import(Request{
		Args: []string{"testdata/templates/template.json"},
		Common: &flags.AssetImportCommon{
			Replace: false,
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res.Text, "test-template")
}
