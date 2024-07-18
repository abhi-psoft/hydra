// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/v2/client"
)

func (p *Persister) CreateClientNSQL(ctx context.Context, c *client.Client) (err error) {
	_, _ = fmt.Println("CreateClient NSQL Hello World: ", uuid.Must(uuid.NewV4()).String())
	return nil
}
