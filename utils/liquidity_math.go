package utils

import (
	"math/big"

	"github.com/KyberNetwork/uniswapv3-sdk-uint256/constants"
)

func AddDelta(x, y *big.Int) *big.Int {
	if y.Cmp(constants.Zero) < 0 {
		return new(big.Int).Sub(x, new(big.Int).Mul(y, constants.NegativeOne))
	} else {
		return new(big.Int).Add(x, y)
	}
}
