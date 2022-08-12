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

func NewClient(fgaAPI *openfga.APIClient) Client {
	return Client{
		fgaAPI: fgaAPI,
	}
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

func (c *Client) SetTypeDefinitions(ctx context.Context, typeDefinitionsContent string) error {
	var typeDefinitions openfga.TypeDefinitions
	err := json.Unmarshal([]byte(typeDefinitionsContent), &typeDefinitions)
	if err != nil {
		return err
	}

	_, _, err = c.fgaAPI.OpenFgaApi.WriteAuthorizationModel(ctx).TypeDefinitions(typeDefinitions).Execute()
	if err != nil {
		return err
	}

	return nil
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
