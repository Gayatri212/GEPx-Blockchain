/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-samples/GEPx-Blockchain/chaincode-go/smart-contract"
)

func main() {
	transactionSmartContract, err := contractapi.NewChaincode(&transaction.SmartContract{})
	if err != nil {
		log.Panicf("Error creating transaction chaincode: %v", err)
	}

	if err := transactionSmartContract.Start(); err != nil {
		log.Panicf("Error starting transaction chaincode: %v", err)
	}
}
