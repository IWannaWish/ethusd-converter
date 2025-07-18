package mapper

import (
	"github.com/IWannaWish/ethusd-converter/cmd/cli/display"
	"github.com/IWannaWish/ethusd-converter/internal/core"
	"math/big"
)

type AssetMapper interface {
	Map(coreAssets []core.Asset) ([]display.AssetInfo, *big.Float, error)
}
