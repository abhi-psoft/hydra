// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package nsql

import (
	"context"

	"github.com/gobuffalo/pop/v6"
	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/x/events"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/sqlcon"
)

func (p *PersisterNSQL) GetConcreteClient(ctx context.Context, id string) (c *client.Client, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetConcreteClient")
	defer otelx.End(span, &err)

	var cl client.Client
	if err := p.QueryWithNetwork(ctx).Where("id = ?", id).First(&cl); err != nil {
		return nil, sqlcon.HandleError(err)
	}
	return &cl, nil
}

func (p *PersisterNSQL) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return p.GetConcreteClient(ctx, id)
}

func (p *PersisterNSQL) UpdateClient(ctx context.Context, cl *client.Client) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateClient")
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		o, err := p.GetConcreteClient(ctx, cl.GetID())
		if err != nil {
			return err
		}

		if cl.Secret == "" {
			cl.Secret = string(o.GetHashedSecret())
		} else {
			h, err := p.r.ClientHasher().Hash(ctx, []byte(cl.Secret))
			if err != nil {
				return errorsx.WithStack(err)
			}
			cl.Secret = string(h)
		}

		// Ensure ID is the same
		cl.ID = o.ID

		if err = cl.BeforeSave(c); err != nil {
			return sqlcon.HandleError(err)
		}

		count, err := p.UpdateWithNetwork(ctx, cl)
		if err != nil {
			return sqlcon.HandleError(err)
		} else if count == 0 {
			return sqlcon.HandleError(sqlcon.ErrNoRows)
		}

		events.Trace(ctx, events.ClientUpdated,
			events.WithClientID(cl.ID),
			events.WithClientName(cl.Name))

		return sqlcon.HandleError(err)
	})
}

func (p *PersisterNSQL) AuthenticateClient(ctx context.Context, id string, secret []byte) (_ *client.Client, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AuthenticateClient")
	defer otelx.End(span, &err)

	c, err := p.GetConcreteClient(ctx, id)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	if err := p.r.ClientHasher().Compare(ctx, c.GetHashedSecret(), secret); err != nil {
		return nil, errorsx.WithStack(err)
	}

	return c, nil
}

func (p *PersisterNSQL) CreateClient(ctx context.Context, c *client.Client) (err error) {
	return nil
}

func (p *PersisterNSQL) DeleteClient(ctx context.Context, id string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteClient")
	defer otelx.End(span, &err)

	c, err := p.GetConcreteClient(ctx, id)
	if err != nil {
		return err
	}

	if err := sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("id = ?", id).Delete(&client.Client{})); err != nil {
		return err
	}

	events.Trace(ctx, events.ClientDeleted,
		events.WithClientID(c.ID),
		events.WithClientName(c.Name))

	return nil
}

func (p *PersisterNSQL) GetClients(ctx context.Context, filters client.Filter) (_ []client.Client, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetClients")
	defer otelx.End(span, &err)

	cs := make([]client.Client, 0)

	query := p.QueryWithNetwork(ctx).
		Paginate(filters.Offset/filters.Limit+1, filters.Limit).
		Order("id")

	if filters.Name != "" {
		query.Where("client_name = ?", filters.Name)
	}
	if filters.Owner != "" {
		query.Where("owner = ?", filters.Owner)
	}

	if err := query.All(&cs); err != nil {
		return nil, sqlcon.HandleError(err)
	}
	return cs, nil
}

func (p *PersisterNSQL) CountClients(ctx context.Context) (n int, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CountClients")
	defer otelx.End(span, &err)

	n, err = p.QueryWithNetwork(ctx).Count(&client.Client{})
	return n, sqlcon.HandleError(err)
}
