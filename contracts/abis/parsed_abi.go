package abis

import (
	"fmt"
	"strings"

	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/celo-org/celo-blockchain/common"
)

var (
	Registry *abi.ABI = mustParseAbi("Registry", RegistryStr)
)

func mustParseAbi(name, abiStr string) *abi.ABI {
	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		panic(fmt.Sprintf("Error reading ABI %s err=%s", name, err))
	}
	return &parsedAbi
}

var byRegistryId = map[common.Hash]*abi.ABI{}

func AbiFor(registryId common.Hash) *abi.ABI {
	found, ok := byRegistryId[registryId]
	if !ok {
		return nil
	}
	return found
}