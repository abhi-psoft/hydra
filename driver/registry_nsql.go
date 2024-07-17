package driver

import "github.com/ory/hydra/v2/client"

type RegistryNSQLImpl struct {
	*RegistrySQL
}

func NewRegistryNSQL() *RegistryNSQLImpl {
	r := &RegistryNSQLImpl{
		RegistrySQL: NewRegistrySQL(),
	}
	return r
}

// Implement the new XClientManager method
func (r *RegistryNSQLImpl) XClientManager() client.Manager {
	return r.nsqlPersister
}
