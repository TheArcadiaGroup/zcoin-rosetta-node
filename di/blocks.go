//+build wireinject

package dependency

import (
	"github.com/google/wire"
	"github.com/arcadiamediagroup/Zcoin-rosetta-node/provider"
)

func InitBadgerDb() (*provider.BadgerDB, error) {
	wire.Build(provider.DatabaseSet)
	return nil, nil // These return values are ignored.
}
