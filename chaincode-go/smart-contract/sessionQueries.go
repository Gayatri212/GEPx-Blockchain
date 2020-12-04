/*
SPDX-License-Identifier: Apache-2.0
*/

package session

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// QuerySession allows all members of the channel to read a public session
func (s *SmartContract) QuerySession(ctx contractapi.TransactionContextInterface, sessionID string) (*Session, error) {

	sessionJSON, err := ctx.GetStub().GetState(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session object %v: %v", sessionID, err)
	}
	if sessionJSON == nil {
		return nil, fmt.Errorf("session does not exist")
	}

	var session *Session
	err = json.Unmarshal(sessionJSON, &session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// QueryBid allows the submitter of the bid to read their bid from public state
func (s *SmartContract) QueryBid(ctx contractapi.TransactionContextInterface, sessionID string, txID string) (*FullBid, error) {

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

	bidKey, err := ctx.GetStub().CreateCompositeKey(bidKeyType, []string{sessionID, txID})
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
func UpdateStatus(ctx contractapi.TransactionContextInterface,sessionID string,sessionJSON Session, bid FullBid, status string, bidKey string) error {
	NewBid := FullBid{
		BidType:  bid.BidType,
		Volume:   bid.Volume,
		Org:      bid.Org,
		Bidder:   bid.Bidder,
		Status:	  status,
	}
	
	revealedBids := make(map[string]FullBid)
	revealedBids = sessionJSON.FinalizedBids
	revealedBids[bidKey] = NewBid
	sessionJSON.FinalizedBids = revealedBids

	newSessionBytes, _ := json.Marshal(sessionJSON)

	// put session with bid added back into state
	err := ctx.GetStub().PutState(sessionID, newSessionBytes)
	if err != nil {
		return fmt.Errorf("failed to update session: %v", err)
	}
	return nil
}