package vmcontext

import (
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/consensus"
	"github.com/celo-org/celo-blockchain/contracts"
	"github.com/celo-org/celo-blockchain/contracts/reserve"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/core/vm"
	"github.com/celo-org/celo-blockchain/log"
	"github.com/celo-org/celo-blockchain/params"
)

// chainContext defines methods required to build a context
// a copy of this exist on core.ChainContext (needed to break dependency)
type chainContext interface {
	// Engine retrieves the blockchain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to the given hash and number.
	GetHeader(common.Hash, uint64) *types.Header

	// Config returns the blockchain's chain configuration
	Config() *params.ChainConfig
}

// New creates a new context for use in the EVM.
func NewBlockContext(header *types.Header, chain chainContext, author *common.Address) vm.BlockContext {
	var (
		beneficiary common.Address
		baseFee     *big.Int
	)

	// If we don't have an explicit author (i.e. not mining), extract from the header
	if author == nil {
		beneficiary, _ = chain.Engine().Author(header) // Ignore error, we're past header validation
	} else {
		beneficiary = *author
	}
	if header.BaseFee != nil {
		baseFee = new(big.Int).Set(header.BaseFee)
	}
	return vm.BlockContext{
		CanTransfer: CanTransfer,
		Transfer:    TobinTransfer,
		GetHash:     GetHashFn(header, chain),
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        new(big.Int).SetUint64(header.Time),
		Difficulty:  new(big.Int).Set(header.Difficulty),
		BaseFee:     baseFee,
		GasLimit:    header.GasLimit,
	}
}

func GetRegisteredAddress(evm *vm.EVM, registryId common.Hash) (common.Address, error) {
	caller := &SharedEVMRunner{evm}
	return contracts.GetRegisteredAddress(caller, registryId)
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(ref *types.Header, chain chainContext) func(uint64) common.Hash {
	// Cache will initially contain [refHash.parent],
	// Then fill up with [refHash.p, refHash.pp, refHash.ppp, ...]
	var cache []common.Hash

	return func(n uint64) common.Hash {
		// If there's no hash cache yet, make one
		if len(cache) == 0 {
			cache = append(cache, ref.ParentHash)
		}
		if idx := ref.Number.Uint64() - n - 1; idx < uint64(len(cache)) {
			return cache[idx]
		}
		// No luck in the cache, but we can start iterating from the last element we already know
		lastKnownHash := cache[len(cache)-1]
		lastKnownNumber := ref.Number.Uint64() - uint64(len(cache))

		for {
			header := chain.GetHeader(lastKnownHash, lastKnownNumber)
			if header == nil {
				break
			}
			cache = append(cache, header.ParentHash)
			lastKnownHash = header.ParentHash
			lastKnownNumber = header.Number.Uint64() - 1
			if n == lastKnownNumber {
				return lastKnownHash
			}
		}
		return common.Hash{}
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas into account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}

// TobinTransfer performs a transfer that may take a tax from the sent amount and give it to the reserve.
// If the calculation or transfer of the tax amount fails for any reason, the regular transfer goes ahead.
// NB: Gas is not charged or accounted for this calculation.
func TobinTransfer(evm *vm.EVM, sender, recipient common.Address, amount *big.Int) {
	// Run only primary evm.Call() with tracer
	if evm.GetDebug() {
		evm.SetDebug(false)
		defer func() { evm.SetDebug(true) }()
	}

	if amount.Cmp(big.NewInt(0)) != 0 {
		caller := &SharedEVMRunner{evm}
		tax, taxRecipient, err := reserve.ComputeTobinTax(caller, sender, amount)
		if err == nil {
			Transfer(evm.StateDB, sender, recipient, new(big.Int).Sub(amount, tax))
			Transfer(evm.StateDB, sender, taxRecipient, tax)
			return
		} else {
			log.Error("Failed to get tobin tax", "error", err)
		}
	}

	// Complete a normal transfer if the amount is 0 or the tobin tax value is unable to be fetched and parsed.
	// We transfer even when the amount is 0 because state trie clearing [EIP161] is necessary at the end of a transaction
	Transfer(evm.StateDB, sender, recipient, amount)
}
