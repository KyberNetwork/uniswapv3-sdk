package periphery

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	AddressSize         = 20
	FeeSize             = 3
	NextOffset          = AddressSize + FeeSize
	PopOffset           = NextOffset + AddressSize
	MultiplePoolsLength = PopOffset + NextOffset
)

func GetFirstPool(path []byte) []byte {
	if len(path) == 0 || len(path) < PopOffset {
		return nil
	}
	return path[0:PopOffset]
}

func HasMultiplePools(path []byte) bool {
	return len(path) >= MultiplePoolsLength
}

func SkipToken(path []byte) []byte {
	if len(path) < NextOffset {
		return nil
	}

	return path[NextOffset : len(path)-NextOffset]
}

func ToAddress(bytes []byte, start *big.Int) (common.Address, error) {
	offset := start.Int64()

	if len(bytes) < int(offset)+20 {
		return common.Address{}, errors.New("toAddress outOfBounds")
	}

	tempBigInt := new(big.Int).SetBytes(bytes[offset : offset+20])
	tempAddress := common.BigToAddress(tempBigInt)

	return tempAddress, nil
}

func ToFee(bytes []byte, offset int) (*big.Int, error) {
	if len(bytes) < offset+3 {
		return nil, errors.New("toAddress outOfBounds")
	}
	var tempUint int64

	tempUint = int64(bytes[offset]) << 16
	tempUint |= int64(bytes[offset+1]) << 8
	tempUint |= int64(bytes[offset+2])

	return new(big.Int).SetInt64(tempUint), nil
}

func DecodeFirstPool(path []byte) (common.Address, common.Address, *big.Int, error) {
	token0, err := ToAddress(path, new(big.Int).SetInt64(0))
	if err != nil {
		return common.Address{}, common.Address{}, nil, err
	}
	fee, err := ToFee(path, AddressSize)
	if err != nil {
		return common.Address{}, common.Address{}, nil, err
	}
	token1, err := ToAddress(path, new(big.Int).SetInt64(NextOffset))
	if err != nil {
		return common.Address{}, common.Address{}, nil, err
	}
	return token0, token1, fee, nil
}
