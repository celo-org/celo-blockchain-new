package testutil

import (
	"github.com/celo-org/celo-blockchain/contracts/abis"
)

type BlockchainParametersMock struct {
	ContractMock
}

func NewBlockchainParametersMock() *BlockchainParametersMock {
	mock := &BlockchainParametersMock{
	}

	contract := NewContractMock(abis.BlockchainParameters, mock)
	mock.ContractMock = contract
	return mock
}
