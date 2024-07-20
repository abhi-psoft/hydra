// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"fmt"

	as "github.com/aerospike/aerospike-client-go/v7"
	"github.com/gofrs/uuid"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/persistence/aerospike"
)

func (p *Persister) CreateClientNSQL(ctx context.Context, c *client.Client) (err error) {
	_, _ = fmt.Println("CreateClient NSQL Hello World: ", aerospike.Client.DefaultPolicy)
	// Define the namespace and set
	namespace := "hydra_org"
	set := "client"

	h, err := p.r.ClientHasher().Hash(ctx, []byte(c.Secret))
	if err != nil {
		return err
	}

	c.Secret = string(h)
	if c.ID == "" {
		c.ID = uuid.Must(uuid.NewV4()).String()
	}

	// Create a unique key for the user
	key, err := as.NewKey(namespace, set, c.ID)
	if err != nil {
		fmt.Println(err)
	}

	// Create bins to hold the user data
	bins := as.BinMap{
		"client_id":     c.ID,
		"client_secret": c.Secret,
	}

	// Write the record to the database
	writePolicy := as.NewWritePolicy(0, 0) // Default write policy
	err = aerospike.Client.Put(writePolicy, key, bins)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
