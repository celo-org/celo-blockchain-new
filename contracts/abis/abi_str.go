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

