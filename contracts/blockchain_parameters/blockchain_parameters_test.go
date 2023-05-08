package blockchain_parameters

import (
	"math/big"
	"testing"

	"github.com/celo-org/celo-blockchain/contracts"
	"github.com/celo-org/celo-blockchain/contracts/testutil"
	"github.com/celo-org/celo-blockchain/params"
	. "github.com/onsi/gomega"
)

func TestGetMinimumVersion(t *testing.T) {
	testutil.TestFailOnFailingRunner(t, getMinimumVersion)
	testutil.TestFailsWhenContractNotDeployed(t, contracts.ErrSmartContractNotDeployed, getMinimumVersion)

	t.Run("should return minimum version", func(t *testing.T) {
		g := NewGomegaWithT(t)

		runner := testutil.NewSingleMethodRunner(
			params.BlockchainParametersRegistryId,
			"getMinimumClientVersion",
			func() (*big.Int, *big.Int, *big.Int) {
				return big.NewInt(5), big.NewInt(4), big.NewInt(3)
			},
		)

		version, err := getMinimumVersion(runner)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(version).To(Equal(&params.VersionInfo{Major: 5, Minor: 4, Patch: 3}))
	})
}