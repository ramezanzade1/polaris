package keeper

import (
	"polaris/x/blog/types"
)

var _ types.QueryServer = Keeper{}
