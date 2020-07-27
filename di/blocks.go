//+build wireinject

package dependency

import (
	"github.com/google/wire"
	"gopkg.in/ArcadiaMediaGroup/zcoin-rosetta-node.v0/provider"
)

func InitBadgerDb() (*provider.BadgerDB, error) {
	wire.Build(provider.DatabaseSet)
	return nil, nil // These return values are ignored.
}
