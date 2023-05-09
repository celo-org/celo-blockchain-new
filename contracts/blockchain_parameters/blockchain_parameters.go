// Copyright 2017 The Celo Authors
// This file is part of the celo library.
//
// The celo library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The celo library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the celo library. If not, see <http://www.gnu.org/licenses/>.

package blockchain_parameters

import (
	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/celo-org/celo-blockchain/contracts"
	"github.com/celo-org/celo-blockchain/log"
	"github.com/celo-org/celo-blockchain/params"
)

var (
)

func logError(method string, err error) {
	if err == contracts.ErrRegistryContractNotDeployed {
		log.Debug("Error calling "+method, "err", err, "contract", hexutil.Encode(params.BlockchainParametersRegistryId[:]))
	} else {
		log.Warn("Error calling "+method, "err", err, "contract", hexutil.Encode(params.BlockchainParametersRegistryId[:]))
	}
}
