package test

import (
	"context"
	"os"
	"testing"

	openfga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/assert"
)

func TestOpenFGAAssertions(t *testing.T) {
	ctx := context.Background()
	typeDefinitions, err := os.ReadFile("../cmd/app/fga/type-definition.json")
	assert.Nil(t, err)

	typeDefinitionID, err := app.FGAClient.SetTypeDefinitions(ctx, string(typeDefinitions))
	assert.Nil(t, err)

	tuples := fixtureOpenFGATuples()
	err = app.FGAClient.WriteTuples(ctx, tuples)
	assert.Nil(t, err)

	assertions := fixtureOpenFGAAssertions()
	err = app.FGAClient.WriteAssertions(ctx, typeDefinitionID, assertions)
	assert.Nil(t, err)

	_, err = app.FGAClient.RunAssertions(ctx, typeDefinitionID)
	assert.Nil(t, err)
}

func fixtureOpenFGAAssertions() []openfga.Assertion {
	return []openfga.Assertion{
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("org-test-editor"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("org-test-editor"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("org-test-viewer"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("org-test-viewer"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-admin-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-admin-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-viewer-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-viewer-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-viewer-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-viewer-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-editor-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-editor-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("organization:default"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-editor-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("organizational_unit:ou-1"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-editor-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("organizational_unit:ou-1"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-viewer-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("organizational_unit:ou-1"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-viewer-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("organizational_unit:ou-1"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-viewer-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("account:account-1"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-viewer-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("account:account-1"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-editor-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("account:account-1"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-editor-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("account:account-1"),
			},
			Expectation: true,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-viewer-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("account:account-2"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-viewer-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("account:account-2"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-editor-user"),
				Relation: openfga.PtrString("editor"),
				Object:   openfga.PtrString("account:account-2"),
			},
			Expectation: false,
		},
		{
			TupleKey: openfga.TupleKey{
				User:     openfga.PtrString("test-ou-1-editor-user"),
				Relation: openfga.PtrString("viewer"),
				Object:   openfga.PtrString("account:account-2"),
			},
			Expectation: false,
		},
	}
}

func fixtureOpenFGATuples() []openfga.TupleKey {
	return []openfga.TupleKey{
		{
			User:     openfga.PtrString("org-test-editor"),
			Relation: openfga.PtrString("editor"),
			Object:   openfga.PtrString("organization:default"),
		},
		{
			User:     openfga.PtrString("org-test-viewer"),
			Relation: openfga.PtrString("viewer"),
			Object:   openfga.PtrString("organization:default"),
		},
		{
			User:     openfga.PtrString("test-admin-user"),
			Relation: openfga.PtrString("member"),
			Object:   openfga.PtrString("group:group-org-admins"),
		},
		{
			User:     openfga.PtrString("group:group-org-admins#member"),
			Relation: openfga.PtrString("editor"),
			Object:   openfga.PtrString("organization:default"),
		},
		{
			User:     openfga.PtrString("test-viewer-user"),
			Relation: openfga.PtrString("member"),
			Object:   openfga.PtrString("group:group-org-viewers"),
		},
		{
			User:     openfga.PtrString("group:group-org-viewers#member"),
			Relation: openfga.PtrString("viewer"),
			Object:   openfga.PtrString("organization:default"),
		},
		{
			User:     openfga.PtrString("test-ou-1-editor-user"),
			Relation: openfga.PtrString("member"),
			Object:   openfga.PtrString("group:group-ou-1-editors"),
		},
		{
			User:     openfga.PtrString("group:group-ou-1-editors#member"),
			Relation: openfga.PtrString("editor"),
			Object:   openfga.PtrString("organizational_unit:ou-1"),
		},
		{
			User:     openfga.PtrString("test-ou-1-viewer-user"),
			Relation: openfga.PtrString("member"),
			Object:   openfga.PtrString("group:group-ou-1-viewers"),
		},
		{
			User:     openfga.PtrString("group:group-ou-1-viewers#member"),
			Relation: openfga.PtrString("viewer"),
			Object:   openfga.PtrString("organizational_unit:ou-1"),
		},
		{
			User:     openfga.PtrString("organizational_unit:ou-1"),
			Relation: openfga.PtrString("parent"),
			Object:   openfga.PtrString("account:account-1"),
		},
		{
			User:     openfga.PtrString("organizational_unit:ou-2"),
			Relation: openfga.PtrString("parent"),
			Object:   openfga.PtrString("account:account-2"),
		},
	}
}
