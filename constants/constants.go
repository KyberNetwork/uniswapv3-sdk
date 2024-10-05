package constants

import (
	"math/big"

	"github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
)

const PoolInitCodeHash = "0xe34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b54"

var (
	FactoryAddress = common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
	AddressZero    = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

// The default factory enabled fee amounts, denominated in hundredths of bips.
type FeeAmount uint64

const (
	FeeLowest FeeAmount = 100
	Fee200    FeeAmount = 200
	Fee300    FeeAmount = 300
	Fee400    FeeAmount = 400
	FeeLow    FeeAmount = 500
	FeeMedium FeeAmount = 3000
	FeeHigh   FeeAmount = 10000

	// These FeeTiers are used in other dexes such as BaseSwapV3, ArbiDexV3.
	Fee1    FeeAmount = 1
	Fee80   FeeAmount = 80
	Fee350  FeeAmount = 350
	Fee450  FeeAmount = 450
	Fee750  FeeAmount = 750
	Fee2500 FeeAmount = 2500

	FeeMax FeeAmount = 1000000
)

// The default factory tick spacings by fee amount.
var TickSpacings = map[FeeAmount]int{
	Fee1:      1,
	Fee80:     1,
	FeeLowest: 1,
	Fee200:    4,
	Fee300:    6,
	Fee350:    10,
	Fee400:    8,
	FeeLow:    10,
	Fee450:    10,
	Fee750:    15,
	FeeMedium: 60,
	Fee2500:   60,
	FeeHigh:   200,
}

var (
	NegativeOne = big.NewInt(-1)
	Zero        = big.NewInt(0)
	One         = big.NewInt(1)

	// used in liquidity amount math
	Q96  = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil)
	Q192 = new(big.Int).Exp(Q96, big.NewInt(2), nil)

	PercentZero = entities.NewFraction(big.NewInt(0), big.NewInt(1))
)
