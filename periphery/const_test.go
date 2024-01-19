package periphery

import (
	"math/big"

	"github.com/KyberNetwork/int256"
	"github.com/KyberNetwork/uniswapv3-sdk-uint256/constants"
	"github.com/KyberNetwork/uniswapv3-sdk-uint256/entities"
	"github.com/KyberNetwork/uniswapv3-sdk-uint256/utils"
	core "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

var (
	ether  = core.EtherOnChain(1)
	token0 = core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "t0", "token0")
	token1 = core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000002"), 18, "t1", "token1")
	token2 = core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000003"), 18, "t2", "token2")
	token3 = core.NewToken(1, common.HexToAddress("0x0000000000000000000000000000000000000004"), 18, "t2", "token3")

	weth = ether.Wrapped()

	pool_0_1_medium, _ = entities.NewPool(token0, token1, constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	pool_1_2_low, _    = entities.NewPool(token1, token2, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	pool_0_weth, _     = entities.NewPool(token0, weth, constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	pool_1_weth, _     = entities.NewPool(token1, weth, constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)

	route_0_1, _   = entities.NewRoute([]*entities.Pool{pool_0_1_medium}, token0, token1)
	route_0_1_2, _ = entities.NewRoute([]*entities.Pool{pool_0_1_medium, pool_1_2_low}, token0, token2)

	route_0_weth, _   = entities.NewRoute([]*entities.Pool{pool_0_weth}, token0, weth)
	route_0_1_weth, _ = entities.NewRoute([]*entities.Pool{pool_0_1_medium, pool_1_weth}, token0, weth)
	route_weth_0, _   = entities.NewRoute([]*entities.Pool{pool_0_weth}, weth, token0)
	route_weth_0_1, _ = entities.NewRoute([]*entities.Pool{pool_0_weth, pool_0_1_medium}, weth, token1)

	liquidityGross  = uint256.NewInt(1_000_000)
	liquidityNet    = int256.NewInt(1_000_000)
	liquidityNetNeg = int256.NewInt(-1_000_000)

	feeAmount    = constants.FeeMedium
	sqrtRatioX96 = utils.EncodeSqrtRatioX96(big.NewInt(1), big.NewInt(1))
	liquidity    = big.NewInt(1_000_000)
	tick, _      = utils.GetTickAtSqrtRatio(sqrtRatioX96)
	ticks        = []entities.Tick{
		{
			Index:          entities.NearestUsableTick(utils.MinTick, constants.TickSpacings[feeAmount]),
			LiquidityNet:   liquidityNet,
			LiquidityGross: liquidityGross,
		},
		{
			Index:          entities.NearestUsableTick(utils.MaxTick, constants.TickSpacings[feeAmount]),
			LiquidityNet:   liquidityNetNeg,
			LiquidityGross: liquidityGross,
		},
	}

	p, _     = entities.NewTickListDataProvider(ticks, constants.TickSpacings[feeAmount])
	makePool = func(token0, token1 *core.Token) *entities.Pool {
		pool, _ := entities.NewPool(token0, token1, feeAmount, sqrtRatioX96, liquidity, tick, p)
		return pool
	}
)
