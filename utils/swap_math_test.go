package utils

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/KyberNetwork/int256"
	"github.com/KyberNetwork/uniswapv3-sdk-uint256/constants"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeSwapStep(t *testing.T) {

	p1 := EncodeSqrtRatioX96(big.NewInt(1), big.NewInt(1))
	p2 := EncodeSqrtRatioX96(big.NewInt(101), big.NewInt(100))
	p3 := EncodeSqrtRatioX96(big.NewInt(1000), big.NewInt(100))
	p4 := EncodeSqrtRatioX96(big.NewInt(10000), big.NewInt(100))

	tests := []struct {
		price       string
		priceTarget string
		liquidity   string
		amount      string
		fee         constants.FeeAmount

		expAmountIn  string
		expAmountOut string
		expFee       string

		expNextPrice string
	}{
		{p1.String(), p2.String(), "2000000000000000000", "1000000000000000000", 600,
			"9975124224178055", "9925619580021728", "5988667735148", "="},
		{p1.String(), p2.String(), "2000000000000000000", "-1000000000000000000", 600,
			"9975124224178055", "9925619580021728", "5988667735148", "="},

		{p1.String(), p3.String(), "2000000000000000000", "1000000000000000000", 600,
			"999400000000000000", "666399946655997866", "600000000000000", "<"},
		{p1.String(), p4.String(), "2000000000000000000", "-1000000000000000000", 600,
			"2000000000000000000", "1000000000000000000", "1200720432259356", "<"},

		{"417332158212080721273783715441582", "1452870262520218020823638996", "159344665391607089467575320103", "-1", 1,
			"1", "1", "1", "417332158212080721273783715441581"},

		{"2", "1", "1", "3915081100057732413702495386755767", 1,
			"39614081257132168796771975168", "0", "39614120871253040049813", "1"},

		{"2413", "79887613182836312", "1985041575832132834610021537970", "10", 1872,
			"0", "0", "10", "2413"},

		{"20282409603651670423947251286016", "22310650564016837466341976414617", "1024", "-4", 3000,
			"26215", "0", "79", "="},

		{"20282409603651670423947251286016", "18254168643286503381552526157414", "1024", "-263000", 3000,
			"1", "26214", "1", "="},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			sqrtRatioNextX96, amountIn, amountOut, feeAmount, err := ComputeSwapStep(
				uint256.MustFromDecimal(tt.price),
				uint256.MustFromDecimal(tt.priceTarget),
				uint256.MustFromDecimal(tt.liquidity),
				int256.MustFromDec(tt.amount),
				tt.fee,
			)
			require.Nil(t, err)
			if tt.expNextPrice == "=" {
				assert.Equal(t, tt.priceTarget, sqrtRatioNextX96.Dec())
			} else if tt.expNextPrice == "<" {
				assert.Greater(t, tt.priceTarget, sqrtRatioNextX96.Dec())
			} else {
				assert.Equal(t, tt.expNextPrice, sqrtRatioNextX96.Dec())
			}
			assert.Equal(t, tt.expAmountIn, amountIn.Dec())
			assert.Equal(t, tt.expAmountOut, amountOut.Dec())
			assert.Equal(t, tt.expFee, feeAmount.Dec())
		})
	}
}
