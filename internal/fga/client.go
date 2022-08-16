package fga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	openfga "github.com/openfga/go-sdk"
)

type Client struct {
	fgaAPI *openfga.APIClient
}

type Authorizer interface {
	Check(ctx context.Context, callerID uuid.UUID, relation Relation, document Document, objectID string) (bool, error)
	SetTypeDefinitions(ctx context.Context, typeDefinitionsContent string) (string, error)
	RunAssertions(ctx context.Context, typeDefinitionsContent string) (bool, error)

	RemoveUser(ctx context.Context, userID uuid.UUID) error
	UserTuples(ctx context.Context, userID uuid.UUID, document string) ([]openfga.TupleKey, error)

	CheckUserAlreadyExistsInOrganization(ctx context.Context, userID uuid.UUID) (bool, error)
	AddUserToOrganization(ctx context.Context, userID uuid.UUID) error
	RemoveUserFromOrganization(ctx context.Context, userID uuid.UUID) error

	CheckUserAlreadyExistsInGroup(ctx context.Context, userID, groupID uuid.UUID) (bool, error)
	AddUsersToGroup(ctx context.Context, userIDs []uuid.UUID, groupID uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) error
	RemoveUsersInGroup(ctx context.Context, groupID uuid.UUID) error
	ReplaceUsersInGroup(ctx context.Context, userIDs []uuid.UUID, groupID uuid.UUID) error

	CheckAccountAlreadyExistsInOrganization(ctx context.Context, accountID uuid.UUID) (bool, error)
	AddAccountToOrganization(ctx context.Context, accountID uuid.UUID) error

	CheckOrganizationalUnitRelationship(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) (bool, error)
	AddOrganizationalUnit(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) error
	RemoveOrganizationalUnitRelationships(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) error
}

func NewClient(fgaAPI *openfga.APIClient) Authorizer {
	return &Client{
		fgaAPI: fgaAPI,
	}
}

func (c *Client) RunAssertions(ctx context.Context, typeDefinitionsContent string) (bool, error) {
	resp, _, err := c.fgaAPI.OpenFgaApi.CreateStore(ctx).Body(openfga.CreateStoreRequest{
		Name: openfga.PtrString("test"),
	}).Execute()
	if err != nil {
		return false, err
	}

	storeID := resp.GetId()

	c.fgaAPI.SetStoreId(storeID)

	defer func(deleteStore openfga.ApiDeleteStoreRequest) {
		_, _ = deleteStore.Execute()
	}(c.fgaAPI.OpenFgaApi.DeleteStore(ctx))

	typeDefinitionID, err := c.SetTypeDefinitions(ctx, typeDefinitionsContent)
	if err != nil {
		return false, err
	}

	err = c.writeAssertionTuples(ctx)
	if err != nil {
		return false, err
	}

	err = c.writeAssertions(ctx, typeDefinitionID)
	if err != nil {
		return false, err
	}

	return c.readAssertions(ctx, typeDefinitionID)
}

func (c *Client) Check(ctx context.Context, callerID uuid.UUID, relation Relation, document Document, objectID string) (bool, error) {
	body := openfga.CheckRequest{
		TupleKey: &openfga.TupleKey{
			User:     openfga.PtrString(callerID.String()),
			Relation: openfga.PtrString(string(relation)),
			Object:   openfga.PtrString(fmt.Sprintf("%s:%s", document, objectID)),
		},
	}

	data, _, err := c.fgaAPI.OpenFgaApi.Check(ctx).Body(body).Execute()
	if err != nil {
		return false, err
	}

	return data.GetAllowed(), nil
}

func (c *Client) SetTypeDefinitions(ctx context.Context, typeDefinitionsContent string) (string, error) {
	var typeDefinitions openfga.TypeDefinitions
	err := json.Unmarshal([]byte(typeDefinitionsContent), &typeDefinitions)
	if err != nil {
		return "", err
	}

	resp, _, err := c.fgaAPI.OpenFgaApi.WriteAuthorizationModel(ctx).TypeDefinitions(typeDefinitions).Execute()
	if err != nil {
		return "", err
	}

	return resp.GetAuthorizationModelId(), nil
}

func (c *Client) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	documents := []string{"organization", "group", "account"}
	for _, document := range documents {
		tuples, err := c.UserTuples(ctx, userID, document)
		if err != nil {
			return err
		} else if len(tuples) == 0 {
			continue
		}

		body := openfga.WriteRequest{
			Deletes: &openfga.TupleKeys{
				TupleKeys: tuples,
			},
		}

		_, _, err = c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) UserTuples(ctx context.Context, userID uuid.UUID, document string) ([]openfga.TupleKey, error) {
	body := openfga.ReadRequest{
		TupleKey: &openfga.TupleKey{
			User:   openfga.PtrString(userID.String()),
			Object: openfga.PtrString(fmt.Sprintf("%s:", document)),
		},
	}

	resp, _, err := c.fgaAPI.OpenFgaApi.Read(ctx).Body(body).Execute()
	if err != nil {
		return nil, err
	}

	tuples := resp.GetTuples()
	tupleKeys := make([]openfga.TupleKey, 0, len(tuples))
	for _, tuple := range tuples {
		tupleKeys = append(tupleKeys, tuple.GetKey())
	}

	return tupleKeys, nil
}

func (c *Client) CheckUserAlreadyExistsInOrganization(ctx context.Context, userID uuid.UUID) (bool, error) {
	body := openfga.ReadRequest{
		TupleKey: &openfga.TupleKey{
			User:     openfga.PtrString(userID.String()),
			Relation: openfga.PtrString("member"),
			Object:   openfga.PtrString("organization:default"),
		},
	}

	resp, _, err := c.fgaAPI.OpenFgaApi.Read(ctx).Body(body).Execute()
	if err != nil {
		return false, err
	}

	return len(*resp.Tuples) > 0, nil
}

func (c *Client) AddUserToOrganization(ctx context.Context, userID uuid.UUID) error {
	alreadyExists, err := c.CheckUserAlreadyExistsInOrganization(ctx, userID)
	if err != nil {
		return err
	} else if alreadyExists {
		return nil
	}

	body := openfga.WriteRequest{
		Writes: &openfga.TupleKeys{
			TupleKeys: []openfga.TupleKey{
				{
					User:     openfga.PtrString(userID.String()),
					Relation: openfga.PtrString("member"),
					Object:   openfga.PtrString("organization:default"),
				},
			},
		},
	}

	_, _, err = c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveUserFromOrganization(ctx context.Context, userID uuid.UUID) error {
	alreadyExists, err := c.CheckUserAlreadyExistsInOrganization(ctx, userID)
	if err != nil {
		return err
	} else if alreadyExists {
		return nil
	}

	body := openfga.WriteRequest{
		Deletes: &openfga.TupleKeys{
			TupleKeys: []openfga.TupleKey{
				{
					User:     openfga.PtrString(userID.String()),
					Relation: openfga.PtrString("member"),
					Object:   openfga.PtrString("organization:default"),
				},
			},
		},
	}

	_, _, err = c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CheckUserAlreadyExistsInGroup(ctx context.Context, userID, groupID uuid.UUID) (bool, error) {
	body := openfga.ReadRequest{
		TupleKey: &openfga.TupleKey{
			User:     openfga.PtrString(userID.String()),
			Relation: openfga.PtrString("member"),
			Object:   openfga.PtrString(fmt.Sprintf("group:%s", groupID.String())),
		},
	}

	resp, _, err := c.fgaAPI.OpenFgaApi.Read(ctx).Body(body).Execute()
	if err != nil {
		return false, err
	}

	return len(*resp.Tuples) > 0, nil
}

func (c *Client) AddUsersToGroup(ctx context.Context, userIDs []uuid.UUID, groupID uuid.UUID) error {
	memberTuples := make([]openfga.TupleKey, 0, len(userIDs))
	for _, member := range userIDs {
		alreadyExists, err := c.CheckUserAlreadyExistsInGroup(ctx, member, groupID)
		if err != nil {
			return err
		} else if alreadyExists {
			continue
		}

		memberTuples = append(memberTuples, openfga.TupleKey{
			User:     openfga.PtrString(member.String()),
			Relation: openfga.PtrString("member"),
			Object:   openfga.PtrString(fmt.Sprintf("group:%s", groupID.String())),
		})
	}

	body := openfga.WriteRequest{
		Writes: &openfga.TupleKeys{
			TupleKeys: memberTuples,
		},
	}

	_, _, err := c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveUserFromGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) error {
	alreadyExists, err := c.CheckUserAlreadyExistsInGroup(ctx, userID, groupID)
	if err != nil {
		return err
	} else if alreadyExists {
		return nil
	}

	body := openfga.WriteRequest{
		Deletes: &openfga.TupleKeys{
			TupleKeys: []openfga.TupleKey{
				{
					User:     openfga.PtrString(userID.String()),
					Relation: openfga.PtrString("member"),
					Object:   openfga.PtrString(fmt.Sprintf("group:%s", groupID.String())),
				},
			},
		},
	}

	_, _, err = c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveUsersInGroup(ctx context.Context, groupID uuid.UUID) error {
	body := openfga.ReadRequest{
		TupleKey: &openfga.TupleKey{
			Relation: openfga.PtrString("member"),
			Object:   openfga.PtrString(fmt.Sprintf("group:%s", groupID.String())),
		},
	}

	for {
		resp, _, err := c.fgaAPI.OpenFgaApi.Read(ctx).Body(body).Execute()
		if err != nil {
			return err
		}

		for _, tuple := range *resp.Tuples {
			userID, err := uuid.Parse(*tuple.Key.User)
			if err != nil {
				return err
			}

			err = c.RemoveUserFromGroup(ctx, userID, groupID)
			if err != nil {
				return err
			}
		}

		if resp.ContinuationToken == nil || *resp.ContinuationToken == "" {
			break
		}

		body.ContinuationToken = resp.ContinuationToken
	}

	return nil
}

func (c *Client) ReplaceUsersInGroup(ctx context.Context, userIDs []uuid.UUID, groupID uuid.UUID) error {
	err := c.RemoveUsersInGroup(ctx, groupID)
	if err != nil {
		return err
	}

	return c.AddUsersToGroup(ctx, userIDs, groupID)
}

func (c *Client) CheckAccountAlreadyExistsInOrganization(ctx context.Context, accountID uuid.UUID) (bool, error) {
	body := openfga.ReadRequest{
		TupleKey: &openfga.TupleKey{
			User:     openfga.PtrString("organization:default"),
			Relation: openfga.PtrString("parent"),
			Object:   openfga.PtrString(fmt.Sprintf("account:%s", accountID.String())),
		},
	}

	resp, _, err := c.fgaAPI.OpenFgaApi.Read(ctx).Body(body).Execute()
	if err != nil {
		return false, err
	}

	return len(*resp.Tuples) > 0, nil
}

func (c *Client) AddAccountToOrganization(ctx context.Context, accountID uuid.UUID) error {
	alreadyExists, err := c.CheckAccountAlreadyExistsInOrganization(ctx, accountID)
	if err != nil {
		return err
	} else if alreadyExists {
		return nil
	}

	memberTuple := openfga.TupleKey{
		User:     openfga.PtrString("organization:default"),
		Relation: openfga.PtrString("parent"),
		Object:   openfga.PtrString(fmt.Sprintf("account:%s", accountID.String())),
	}

	body := openfga.WriteRequest{
		Writes: &openfga.TupleKeys{
			TupleKeys: []openfga.TupleKey{memberTuple},
		},
	}

	_, _, err = c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CheckOrganizationalUnitRelationship(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) (bool, error) {
	body := openfga.ReadRequest{
		TupleKey: &openfga.TupleKey{
			User:     openfga.PtrString(c.organizationalUnitParentID(parentID)),
			Object:   openfga.PtrString(fmt.Sprintf("organizational_unit:%s", id.String())),
			Relation: openfga.PtrString("parent"),
		},
	}

	resp, _, err := c.fgaAPI.OpenFgaApi.Read(ctx).Body(body).Execute()
	if err != nil {
		return false, err
	}

	return len(*resp.Tuples) > 0, nil
}

func (c *Client) AddOrganizationalUnit(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) error {
	alreadyExists, err := c.CheckOrganizationalUnitRelationship(ctx, id, parentID)
	if err != nil {
		return err
	} else if alreadyExists {
		return nil
	}

	memberTuple := openfga.TupleKey{
		User:     openfga.PtrString(c.organizationalUnitParentID(parentID)),
		Object:   openfga.PtrString(fmt.Sprintf("organizational_unit:%s", id.String())),
		Relation: openfga.PtrString("parent"),
	}

	body := openfga.WriteRequest{
		Writes: &openfga.TupleKeys{
			TupleKeys: []openfga.TupleKey{memberTuple},
		},
	}

	_, _, err = c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveOrganizationalUnitRelationships(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) error {
	alreadyExists, err := c.CheckOrganizationalUnitRelationship(ctx, id, parentID)
	if err != nil {
		return err
	} else if !alreadyExists {
		return nil
	}

	body := openfga.WriteRequest{
		Deletes: &openfga.TupleKeys{
			TupleKeys: []openfga.TupleKey{
				{
					User:     openfga.PtrString(c.organizationalUnitParentID(parentID)),
					Object:   openfga.PtrString(fmt.Sprintf("organizational_unit:%s", id.String())),
					Relation: openfga.PtrString("parent"),
				},
			},
		},
	}

	_, _, err = c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) organizationalUnitParentID(parentID uuid.NullUUID) string {
	if parentID.Valid {
		return fmt.Sprintf("organizational_unit:%s", parentID.UUID.String())
	}

	return "organization:default"
}

func (c *Client) readAssertions(ctx context.Context, authorizationModelID string) (bool, error) {
	resp, _, err := c.fgaAPI.OpenFgaApi.ReadAssertions(ctx, authorizationModelID).Execute()
	if err != nil {
		return false, err
	}

	for _, assertion := range *resp.Assertions {
		body := openfga.CheckRequest{
			TupleKey: assertion.TupleKey,
		}

		data, _, err := c.fgaAPI.OpenFgaApi.Check(ctx).Body(body).Execute()
		if err != nil {
			return false, err
		} else if data.GetAllowed() != assertion.Expectation {
			fmt.Println("Assertion failed:", assertion.TupleKey.GetUser(), assertion.TupleKey.GetRelation(), assertion.TupleKey.GetObject())
			return false, nil
		}
	}

	return true, nil
}

func (c *Client) writeAssertionTuples(ctx context.Context) error {
	body := openfga.WriteRequest{
		Writes: &openfga.TupleKeys{
			TupleKeys: []openfga.TupleKey{
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
			},
		},
	}

	_, _, err := c.fgaAPI.OpenFgaApi.Write(ctx).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) writeAssertions(ctx context.Context, authorizationModelID string) error {
	body := openfga.WriteAssertionsRequest{
		Assertions: []openfga.Assertion{
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("org-test-editor"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("org-test-editor"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("org-test-viewer"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("org-test-viewer"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-admin-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-admin-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-viewer-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-viewer-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-viewer-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-viewer-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-editor-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-editor-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("organization:default"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-editor-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("organizational_unit:ou-1"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-editor-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("organizational_unit:ou-1"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-viewer-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("organizational_unit:ou-1"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-viewer-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("organizational_unit:ou-1"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-viewer-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("account:account-1"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-viewer-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("account:account-1"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-editor-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("account:account-1"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-editor-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("account:account-1"),
				},
				Expectation: true,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-viewer-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("account:account-2"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-viewer-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("account:account-2"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-editor-user"),
					Relation: openfga.PtrString("editor"),
					Object:   openfga.PtrString("account:account-2"),
				},
				Expectation: false,
			},
			{
				TupleKey: &openfga.TupleKey{
					User:     openfga.PtrString("test-ou-1-editor-user"),
					Relation: openfga.PtrString("viewer"),
					Object:   openfga.PtrString("account:account-2"),
				},
				Expectation: false,
			},
		},
	}

	_, err := c.fgaAPI.OpenFgaApi.WriteAssertions(ctx, authorizationModelID).Body(body).Execute()
	if err != nil {
		return err
	}

	return nil
}
