package abis

// This is taken from celo-monorepo/packages/protocol/build/<env>/contracts/Registry.json
const RegistryStr = `[
	{
		"constant": true,
		"inputs": [
			{
				"name": "identifier",
				"type": "bytes32"
			}
		],
		"name": "getAddressFor",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`

const BlockchainParametersStr = `[
	{
		"constant": true,
		"inputs": [],
		"name": "getMinimumClientVersion",
		"outputs": [
			{
			"name": "major",
			"type": "uint256"
			},
			{
			"name": "minor",
			"type": "uint256"
			},
			{
			"name": "patch",
			"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`
