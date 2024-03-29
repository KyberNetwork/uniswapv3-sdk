package entities

import (
	"math/big"
	"testing"

	"github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var (
	USDC     = entities.NewToken(1, common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), 6, "USDC", "USD Coin")
	USDT     = entities.NewToken(1, common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"), 6, "USDT", "Tether USD")
	DAI      = entities.NewToken(1, common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), 18, "DAI", "Dai Stablecoin")
	WETH     = entities.NewToken(1, common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), 18, "WETH", "Wrapped Ether")
	OneEther = big.NewInt(1e18)
)

func TestNewPool(t *testing.T) {
	_, err := NewPool(USDC, entities.WETH9[3], constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.ErrorIs(t, err, entities.ErrDifferentChain, "cannot be used for tokens on different chains")

	_, err = NewPool(USDC, entities.WETH9[1], 1e6, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.ErrorIs(t, err, ErrFeeTooHigh, "fee cannot be more than 1e6'")

	_, err = NewPool(USDC, USDC, constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.ErrorIs(t, err, entities.ErrSameAddress, "cannot be used for the same token")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 1, nil)
	assert.ErrorIs(t, err, ErrInvalidSqrtRatioX96, "price must be within tick price bounds")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeMedium, new(big.Int).Add(utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(1)), big.NewInt(0), -1, nil)
	assert.ErrorIs(t, err, ErrInvalidSqrtRatioX96, "price must be within tick price bounds")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool medium fee")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool low fee")

	_, err = NewPool(USDC, entities.WETH9[1], constants.FeeHigh, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.NoError(t, err, "works with valid arguments for empty pool high fee")
}

func TestGetAddress(t *testing.T) {
	addr, _ := GetAddress(USDC, DAI, constants.FeeLow, "")
	assert.Equal(t, addr, common.HexToAddress("0x6c6Bc977E13Df9b0de53b251522280BB72383700"), "matches an example")
}

func TestToken0(t *testing.T) {
	pool, _ := NewPool(USDC, DAI, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool.Token0, DAI, "always is the token that sorts before")

	pool, _ = NewPool(DAI, USDC, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool.Token0, DAI, "always is the token that sorts before")
}

func TestToken1(t *testing.T) {
	pool, _ := NewPool(USDC, DAI, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool.Token1, USDC, "always is the token that sorts after")

	pool, _ = NewPool(DAI, USDC, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool.Token1, USDC, "always is the token that sorts after")
}

func TestToken0Price(t *testing.T) {
	a1 := new(big.Int).Mul(big.NewInt(101), big.NewInt(1e6))
	a2 := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	r, _ := utils.GetTickAtSqrtRatio(utils.EncodeSqrtRatioX96(a1, a2))
	pool0, _ := NewPool(USDC, DAI, constants.FeeLow, utils.EncodeSqrtRatioX96(a1, a2), big.NewInt(0), r, nil)
	assert.Equal(t, pool0.Token0Price().ToSignificant(5), "1.01", "returns price of token0 in terms of token1")

	pool1, _ := NewPool(DAI, USDC, constants.FeeLow, utils.EncodeSqrtRatioX96(a1, a2), big.NewInt(0), r, nil)
	assert.Equal(t, pool1.Token0Price().ToSignificant(5), "1.01", "returns price of token0 in terms of token1")
}

func TestToken1Price(t *testing.T) {
	a1 := new(big.Int).Mul(big.NewInt(101), big.NewInt(1e6))
	a2 := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	r, _ := utils.GetTickAtSqrtRatio(utils.EncodeSqrtRatioX96(a1, a2))
	pool0, _ := NewPool(USDC, DAI, constants.FeeLow, utils.EncodeSqrtRatioX96(a1, a2), big.NewInt(0), r, nil)
	assert.Equal(t, pool0.Token1Price().ToSignificant(5), "0.9901", "returns price of token1 in terms of token0")

	pool1, _ := NewPool(DAI, USDC, constants.FeeLow, utils.EncodeSqrtRatioX96(a1, a2), big.NewInt(0), r, nil)
	assert.Equal(t, pool1.Token1Price().ToSignificant(5), "0.9901", "returns price of token1 in terms of token0")
}

func TestPriceOf(t *testing.T) {
	pool, _ := NewPool(USDC, DAI, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	price0, _ := pool.PriceOf(DAI)
	assert.Equal(t, price0, pool.Token0Price(), "returns price of token in terms of other token")
	price1, _ := pool.PriceOf(USDC)
	assert.Equal(t, price1, pool.Token1Price(), "returns price of token in terms of other token")

	_, err := pool.PriceOf(entities.WETH9[1])
	assert.Error(t, err, "invalid token")
}

func TestChainID(t *testing.T) {
	pool0, _ := NewPool(USDC, DAI, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool0.ChainID(), uint(1), "returns the token0 chainId")

	pool1, _ := NewPool(DAI, USDC, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.Equal(t, pool1.ChainID(), uint(1), "returns the token0 chainId")
}

func TestInvolvesToken(t *testing.T) {
	pool, _ := NewPool(USDC, DAI, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	assert.True(t, pool.InvolvesToken(USDC), "involves USDC")
	assert.True(t, pool.InvolvesToken(DAI), "involves DAI")
	assert.False(t, pool.InvolvesToken(entities.WETH9[1]), "does not involve WETH9")
}

func newTestPool() *Pool {
	ticks := []Tick{
		{
			Index:          NearestUsableTick(utils.MinTick, constants.TickSpacings[constants.FeeLow]),
			LiquidityNet:   OneEther,
			LiquidityGross: OneEther,
		},
		{
			Index:          NearestUsableTick(utils.MaxTick, constants.TickSpacings[constants.FeeLow]),
			LiquidityNet:   new(big.Int).Mul(OneEther, constants.NegativeOne),
			LiquidityGross: OneEther,
		},
	}

	p, err := NewTickListDataProvider(ticks, constants.TickSpacings[constants.FeeLow])
	if err != nil {
		panic(err)
	}

	pool, err := NewPool(USDC, DAI, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), OneEther, 0, p)
	if err != nil {
		panic(err)
	}
	return pool
}
func TestGetOutputAmount(t *testing.T) {
	pool := newTestPool()

	// USDC -> DAI
	inputAmount := entities.FromRawAmount(USDC, big.NewInt(100))
	outputAmount, err := pool.GetOutputAmount(inputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, outputAmount.ReturnedAmount.Currency.Equal(DAI))
	assert.Equal(t, outputAmount.ReturnedAmount.Quotient(), big.NewInt(98))

	// DAI -> USDC
	inputAmount = entities.FromRawAmount(DAI, big.NewInt(100))
	outputAmount, err = pool.GetOutputAmount(inputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, outputAmount.ReturnedAmount.Currency.Equal(USDC))
	assert.Equal(t, outputAmount.ReturnedAmount.Quotient(), big.NewInt(98))
}

func TestGetInputAmount(t *testing.T) {
	pool := newTestPool()

	// USDC -> DAI
	outputAmount := entities.FromRawAmount(DAI, big.NewInt(98))
	getInputAmountResult, err := pool.GetInputAmount(outputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, getInputAmountResult.ReturnedAmount.Currency.Equal(USDC))
	assert.Equal(t, getInputAmountResult.ReturnedAmount.Quotient(), big.NewInt(100))

	// DAI -> USDC
	outputAmount = entities.FromRawAmount(USDC, big.NewInt(98))
	getInputAmountResult, err = pool.GetInputAmount(outputAmount, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, getInputAmountResult.ReturnedAmount.Currency.Equal(DAI))
	assert.Equal(t, getInputAmountResult.ReturnedAmount.Quotient(), big.NewInt(100))
}

func newTestBlueprintV3Pool() *Pool {
	liquidityGross, _ := new(big.Int).SetString("5158244147190", 10)
	liquidityNet, _ := new(big.Int).SetString("-5158244147190", 10)

	ticksData := []Tick{
		{
			Index:          -197760,
			LiquidityNet:   liquidityGross,
			LiquidityGross: liquidityGross,
		},
		{
			Index:          -196320,
			LiquidityNet:   liquidityNet,
			LiquidityGross: liquidityGross,
		},
	}

	ticks, err := NewTickListDataProvider(ticksData, constants.TickSpacings[constants.FeeMedium])
	if err != nil {
		panic(err)
	}

	sqrtPriceX96, _ := new(big.Int).SetString("4736796043499064985766337", 10)

	pool, err := NewPool(
		WETH,
		USDT,
		constants.FeeMedium,
		sqrtPriceX96,
		constants.Zero,
		-194505,
		ticks,
	)
	if err != nil {
		panic(err)
	}

	return pool
}

func TestGetInputAmount2(t *testing.T) {
	blueprintV3Pool := newTestBlueprintV3Pool()

	sqrtPriceLimitX96, _ := new(big.Int).SetString("4025227652667870579138333", 10)
	// WETH -> USDT
	outputAmount := entities.FromRawAmount(USDT, big.NewInt(500000000))
	getInputAmountResult, err := blueprintV3Pool.GetInputAmount(outputAmount, sqrtPriceLimitX96)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, getInputAmountResult.ReturnedAmount.Currency.Equal(WETH))
	assert.Equal(t, getInputAmountResult.ReturnedAmount.Quotient(), big.NewInt(7074025631378098))
	assert.Equal(t, getInputAmountResult.RemainingAmountOut.Quotient(), big.NewInt(-480436293))
}
