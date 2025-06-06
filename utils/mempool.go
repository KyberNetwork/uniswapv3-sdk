package utils

import (
	"math/big"
	"sync"
)

var intPool = &sync.Pool{
	New: func() interface{} {
		return new(big.Int)
	},
}

func GetBigInt() *big.Int {
	return intPool.Get().(*big.Int)
}

func PutBigInt(i *big.Int) {
	i.SetInt64(0) // or you might want to reset to some default value
	intPool.Put(i)
}
