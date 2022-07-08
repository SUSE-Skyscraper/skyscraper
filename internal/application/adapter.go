package application

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type Adapter struct {
	app      *App
	filtered bool
}

func NewAdapter(app *App) *Adapter {
	return &Adapter{
		app: app,
	}
}

func (a *Adapter) IsFiltered() bool {
	return a.filtered
}

func (a *Adapter) LoadPolicy(model model.Model) error {
	policies, err := a.app.DB.GetPolicies(context.Background())
	if err != nil {
		return err
	}

	for _, policy := range policies {
		persist.LoadPolicyLine(policyToString(policy), model)
	}

	a.filtered = false

	return nil
}

func (a *Adapter) SavePolicy(model model.Model) error {
	ctx := context.Background()

	tx, err := a.app.PostgresPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	q := db.New(tx)

	err = q.TruncatePolicies(ctx)
	if err != nil {
		return err
	}

	var toInput []db.AddPolicyParams

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			input := addPolicyInput(ptype, rule)
			toInput = append(toInput, input)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			input := addPolicyInput(ptype, rule)
			toInput = append(toInput, input)
		}
	}

	for _, input := range toInput {
		err = q.AddPolicy(ctx, input)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *Adapter) AddPolicy(_ string, ptype string, rule []string) error {
	input := addPolicyInput(ptype, rule)
	return a.app.DB.AddPolicy(context.Background(), input)
}

func (a *Adapter) RemovePolicy(_ string, ptype string, rule []string) error {
	input := removePolicyInput(ptype, rule)
	return a.app.DB.RemovePolicy(context.Background(), input)
}

func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}

func removePolicyInput(ptype string, rule []string) db.RemovePolicyParams {
	input := db.RemovePolicyParams{
		Ptype: ptype,
		V0:    rule[0],
		V1:    rule[1],
	}

	l := len(rule)
	if l > 2 {
		input.V2 = sql.NullString{
			String: rule[2],
			Valid:  true,
		}
	}
	if l > 3 {
		input.V3 = sql.NullString{
			String: rule[3],
			Valid:  true,
		}
	}
	if l > 4 {
		input.V4 = sql.NullString{
			String: rule[4],
			Valid:  true,
		}
	}
	if l > 5 {
		input.V5 = sql.NullString{
			String: rule[5],
			Valid:  true,
		}
	}

	return input
}

func addPolicyInput(ptype string, rule []string) db.AddPolicyParams {
	input := db.AddPolicyParams{
		Ptype: ptype,
		V0:    rule[0],
		V1:    rule[1],
	}

	l := len(rule)
	if l > 2 {
		input.V2 = rule[2]
	}
	if l > 3 {
		input.V3 = rule[3]
	}
	if l > 4 {
		input.V4 = rule[4]
	}
	if l > 5 {
		input.V5 = rule[5]
	}

	return input
}

func policyToString(policy db.Policy) string {
	const prefixLine = ", "
	var sb strings.Builder

	length := len(policy.Ptype) + len(policy.V0) + len(policy.V1) + len(policy.V2.String) +
		len(policy.V3.String) + len(policy.V4.String) + len(policy.V5.String)

	sb.Grow(length)

	sb.WriteString(policy.Ptype)
	sb.WriteString(prefixLine)
	sb.WriteString(policy.V0)
	sb.WriteString(prefixLine)
	sb.WriteString(policy.V1)

	if policy.V2.Valid {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V2.String)
	}

	if policy.V3.Valid {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V3.String)
	}

	if policy.V4.Valid {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V4.String)
	}

	if policy.V5.Valid {
		sb.WriteString(prefixLine)
		sb.WriteString(policy.V5.String)
	}

	return sb.String()
}
