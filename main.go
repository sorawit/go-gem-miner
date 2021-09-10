package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

type MiningInfo struct {
	chainId    string
	entropy    []byte
	gem        common.Address
	miner      common.Address
	kind       string
	nonce      string
	difficulty *big.Int
}

type Result struct {
	// TODO
}

/// Converts the given string text to a big integer. Panic on failure.
func textToBigInt(text string) *big.Int {
	v := new(big.Int)
	_, ok := v.SetString(text, 10)
	if !ok {
		panic(fmt.Sprintf("textoBigInt failed: %s is not an integer", text))
	}
	return v
}

func run(iter chan int64, batch int64, salt *big.Int, info *MiningInfo) {
	one := big.NewInt(1)
	result := new(big.Int)
	maxUint256 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	target := new(big.Int).Div(maxUint256, info.difficulty)
	for {
		select {
		default:
			for i := int64(0); i < batch; i++ {
				hash := solsha3.SoliditySHA3(
					[]string{
						"uint256", "bytes32", "address",
						"address", "uint256", "uint256", "uint256",
					},
					[]interface{}{
						info.chainId, info.entropy, info.gem,
						info.miner, info.kind, info.nonce, salt,
					},
				)
				result.SetBytes(hash)
				if result.Cmp(target) < 0 {
					// TODO
				}
				salt = salt.Add(salt, one)
			}
			iter <- batch
		}
	}
}

func main() {
	n := 16
	batch := int64(5000)
	iter := make(chan int64)
	for i := 0; i < n; i++ {
		go run(iter, batch, big.NewInt(0), &MiningInfo{
			difficulty: big.NewInt(1),
			// chainId: "1",
			// entropy: "",
		})
	}
	total := int64(0)
	start := time.Now().UnixNano()
	for {
		count := <-iter
		total += count
		now := time.Now().UnixNano()
		fmt.Printf("\rtotal hashes %d, hashes per second : %d", total, total/((now-start)/1e9+1))
	}
}
