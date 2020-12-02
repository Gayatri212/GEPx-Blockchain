/*
SPDX-License-Identifier: Apache-2.0
*/

package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// QueryTransaction allows all members of the channel to read a public transaction
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, transactionID string) (*Transaction, error) {

	transactionJSON, err := ctx.GetStub().GetState(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction object %v: %v", transactionID, err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("transaction does not exist")
	}

	var transaction *Transaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// QueryBid allows the submitter of the bid to read their bid from public state
func (s *SmartContract) QueryBid(ctx contractapi.TransactionContextInterface, transactionID string, txID string) (*FullBid, error) {

	err := verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get implicit collection name: %v", err)
	}

	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client identity %v", err)
	}

	collection, err := getCollectionName(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get implicit collection name: %v", err)
	}

	bidKey, err := ctx.GetStub().CreateCompositeKey(bidKeyType, []string{transactionID, txID})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %v", err)
	}

	bidJSON, err := ctx.GetStub().GetPrivateData(collection, bidKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get bid %v: %v", bidKey, err)
	}
	if bidJSON == nil {
		return nil, fmt.Errorf("bid %v does not exist", bidKey)
	}

	var bid *FullBid
	err = json.Unmarshal(bidJSON, &bid)
	if err != nil {
		return nil, err
	}

	// check that the client querying the bid is the bid owner
	if bid.Bidder != clientID {
		return nil, fmt.Errorf("Permission denied, client id %v is not the owner of the bid", clientID)
	}

	return bid, nil
}

// GetID is an internal helper function to allow users to get their identity
func (s *SmartContract) GetID(ctx contractapi.TransactionContextInterface) (string, error) {

	// Get the MSP ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get verified MSPID: %v", err)
	}

	return clientID, nil
}

// Internal call to update status of bids
func UpdateStatus(ctx contractapi.TransactionContextInterface,transactionID string,transactionJSON Transaction, bid FullBid, status string, bidKey string) error {
	NewBid := FullBid{
		BidType:  bid.BidType,
		Volume:   bid.Volume,
		Org:      bid.Org,
		Bidder:   bid.Bidder,
		Status:	  status,
	}
	
	revealedBids := make(map[string]FullBid)
	revealedBids = transactionJSON.RevealedBids
	revealedBids[bidKey] = NewBid
	transactionJSON.RevealedBids = revealedBids

	newTransactionBytes, _ := json.Marshal(transactionJSON)

	// put transaction with bid added back into state
	err := ctx.GetStub().PutState(transactionID, newTransactionBytes)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}
	return nil
}