package testutil

import (
	"github.com/celo-org/celo-blockchain/params"
)

type CeloMock struct {
	Runner   *MockEVMRunner
	Registry *RegistryMock
}

func NewCeloMock() CeloMock {
	celo := CeloMock{
		Runner:   NewMockEVMRunner(),
		Registry: NewRegistryMock(),
	}

	celo.Runner.RegisterContract(params.RegistrySmartContractAddress, celo.Registry)

	return celo
}
