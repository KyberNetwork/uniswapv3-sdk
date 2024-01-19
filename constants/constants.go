package constants

import (
	"math/big"

	"github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
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
	FeeLow    FeeAmount = 500
	FeeMedium FeeAmount = 3000
	FeeHigh   FeeAmount = 10000

	// These FeeTiers are used in other dexes such as BaseSwapV3, ArbiDexV3.
	Fee80   FeeAmount = 80
	Fee450  FeeAmount = 450
	Fee2500 FeeAmount = 2500

	FeeMax FeeAmount = 1000000
)

// The default factory tick spacings by fee amount.
var TickSpacings = map[FeeAmount]int{
	FeeLowest: 1,
	Fee80:     1,
	FeeLow:    10,
	Fee450:    10,
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

	Q96U256 = new(uint256.Int).Exp(uint256.NewInt(2), uint256.NewInt(96))

	PercentZero = entities.NewFraction(big.NewInt(0), big.NewInt(1))
)
