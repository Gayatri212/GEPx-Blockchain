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
	sessionSmartContract, err := contractapi.NewChaincode(&session.SmartContract{})
	if err != nil {
		log.Panicf("Error creating session chaincode: %v", err)
	}

	if err := sessionSmartContract.Start(); err != nil {
		log.Panicf("Error starting session chaincode: %v", err)
	}
}
