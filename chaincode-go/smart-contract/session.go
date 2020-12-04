/*
SPDX-License-Identifier: Apache-2.0
*/

package session

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// Session data
type Session struct {
	Admin       	string             `json:"admin"`
	Orgs         	[]string           `json:"organizations"`
	PrivateBids  	map[string]BidHash `json:"privateBids"`
	FinalizedBids 	map[string]FullBid `json:"finalizedBids"`
	Status       	string             `json:"status"`
}

// FullBid is the structure of a revealed bid
type FullBid struct {
	BidType     BidType `json:"bidType"`
	Volume    	int    	`json:"volume"`
	Org      	string 	`json:"org"`
	Bidder   	string 	`json:"bidder"`
	Status      string  `json:"status"`
}

// BidHash is the structure of a private bid
type BidHash struct {
	Org  string `json:"org"`
	Hash string `json:"hash"`
}

type BidType string
const(
	Sell = "sell"
	Buy = "buy"
)

const bidKeyType = "bid"

// CreateSession creates on session on the public channel. The identity that
// submits the transacion becomes the admin of the session
func (s *SmartContract) CreateSession(ctx contractapi.TransactionContextInterface, sessionID string) error {

	// get ID of submitting client
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

	// get org of submitting client
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

	// Create session
	bidders := make(map[string]BidHash)
	finalizedBids := make(map[string]FullBid)

	session := Session{
		Admin:     		clientID,
		Orgs:       	[]string{clientOrgID},
		PrivateBids:  	bidders,
		FinalizedBids: 	finalizedBids,
		Status:     	"Open",
	}

	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// put session into state
	err = ctx.GetStub().PutState(sessionID, sessionBytes)
	if err != nil {
		return fmt.Errorf("failed to put session in public data: %v", err)
	}

	// set the admin of the session as an endorser
	err = setAssetStateBasedEndorsement(ctx, sessionID, clientOrgID)
	if err != nil {
		return fmt.Errorf("failed setting state based endorsement for new organization: %v", err)
	}

	return nil
}

// Bid is used to add a user's bid to the session. The bid is stored in the private
// data collection on the peer of the bidder's organization. The function returns
// the session ID so that users can identify and query their bid
func (s *SmartContract) Bid(ctx contractapi.TransactionContextInterface, sessionID string) (string, error) {

	// get bid from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "", fmt.Errorf("error getting transient: %v", err)
	}

	BidJSON, ok := transientMap["bid"]
	if !ok {
		return "", fmt.Errorf("bid key not found in the transient map")
	}

	// get the implicit collection name using the bidder's organization ID
	collection, err := getCollectionName(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get implicit collection name: %v", err)
	}

	// the bidder has to target their peer to store the bid
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return "", fmt.Errorf("Cannot store bid on this peer, not a member of this org: Error %v", err)
	}

	// the session ID is used as a unique index for the bid
	txID := ctx.GetStub().GetTxID()

	// create a composite key using the session ID
	bidKey, err := ctx.GetStub().CreateCompositeKey(bidKeyType, []string{sessionID, txID})
	if err != nil {
		return "", fmt.Errorf("failed to create composite key: %v", err)
	}

	// put the bid into the organization's implicit data collection
	err = ctx.GetStub().PutPrivateData(collection, bidKey, BidJSON)
	if err != nil {
		return "", fmt.Errorf("failed to input volume into collection: %v", err)
	}

	// return the trannsaction ID so that the uset can identify their bid
	return txID, nil
}

// SubmitBid is used by the bidder to add the hash of that bid stored in private data to the
// session. Note that this function alters the session in private state, and needs
// to meet the session endorsement policy. Session ID is used identify the bid
func (s *SmartContract) SubmitBid(ctx contractapi.TransactionContextInterface, sessionID string, txID string) error {

	// get the MSP ID of the bidder's org
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP ID: %v", err)
	}

	// get the session from state
	sessionBytes, err := ctx.GetStub().GetState(sessionID)
	var sessionJSON Session

	if sessionBytes == nil {
		return fmt.Errorf("Session not found: %v", sessionID)
	}
	err = json.Unmarshal(sessionBytes, &sessionJSON)
	if err != nil {
		return fmt.Errorf("failed to create session object JSON: %v", err)
	}

	// the session needs to be Placed for users to add their bid
	Status := sessionJSON.Status
	if Status != "Open" {
		return fmt.Errorf("cannot change finalized or ended session")
	}

	// get the inplicit collection name of bidder's org
	collection, err := getCollectionName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get implicit collection name: %v", err)
	}

	// use the session ID passed as a parameter to create composite bid key
	bidKey, err := ctx.GetStub().CreateCompositeKey(bidKeyType, []string{sessionID, txID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// get the hash of the bid stored in private data collection
	bidHash, err := ctx.GetStub().GetPrivateDataHash(collection, bidKey)
	if err != nil {
		return fmt.Errorf("failed to read bid bash from collection: %v", err)
	}
	if bidHash == nil {
		return fmt.Errorf("bid hash does not exist: %s", bidKey)
	}

	// store the hash along with the bidder's organization
	NewHash := BidHash{
		Org:  clientOrgID,
		Hash: fmt.Sprintf("%x", bidHash),
	}

	bidders := make(map[string]BidHash)
	bidders = sessionJSON.PrivateBids
	bidders[bidKey] = NewHash
	sessionJSON.PrivateBids = bidders

	// Add the bidding organization to the list of participating organizations if it is not already
	Orgs := sessionJSON.Orgs
	if !(contains(Orgs, clientOrgID)) {
		newOrgs := append(Orgs, clientOrgID)
		sessionJSON.Orgs = newOrgs

		err = addAssetStateBasedEndorsement(ctx, sessionID, clientOrgID)
		if err != nil {
			return fmt.Errorf("failed setting state based endorsement for new organization: %v", err)
		}
	}

	newTransactionBytes, _ := json.Marshal(sessionJSON)

	err = ctx.GetStub().PutState(sessionID, newTransactionBytes)
	if err != nil {
		return fmt.Errorf("failed to update session: %v", err)
	}

	return nil
}

// FinalizeBid is used by a bidder to reveal their bid after the session is Finalized
func (s *SmartContract) FinalizeBid(ctx contractapi.TransactionContextInterface, sessionID string, txID string) error {

	// get bid from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	transientBidJSON, ok := transientMap["bid"]
	if !ok {
		return fmt.Errorf("bid key not found in the transient map")
	}

	// get implicit collection name of organization ID
	collection, err := getCollectionName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get implicit collection name: %v", err)
	}

	// use session ID to create composit bid key
	bidKey, err := ctx.GetStub().CreateCompositeKey(bidKeyType, []string{sessionID, txID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	// get bid hash of bid if private bid on the public ledger
	bidHash, err := ctx.GetStub().GetPrivateDataHash(collection, bidKey)
	if err != nil {
		return fmt.Errorf("failed to read bid bash from collection: %v", err)
	}
	if bidHash == nil {
		return fmt.Errorf("bid hash does not exist: %s", bidKey)
	}

	// get session from public state
	sessionBytes, err := ctx.GetStub().GetState(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session %v: %v", sessionID, err)
	}
	if sessionBytes == nil {
		return fmt.Errorf("Session interest object %v not found", sessionID)
	}

	var sessionJSON Session
	err = json.Unmarshal(sessionBytes, &sessionJSON)
	if err != nil {
		return fmt.Errorf("failed to create session object JSON: %v", err)
	}

	// Complete a series of three checks before we add the bid to the session

	// check 1: check that the session is Finalized. We cannot reveal a
	// bid to an Placed session
	Status := sessionJSON.Status
	if Status != "Close" {
		return fmt.Errorf("cannot reveal bid for Placed or ended session")
	}

	// check 2: check that hash of revealed bid matches hash of private bid
	// on the public ledger. This checks that the bidder is telling the truth
	// about the value of their bid

	hash := sha256.New()
	hash.Write(transientBidJSON)
	calculatedBidJSONHash := hash.Sum(nil)

	// verify that the hash of the passed immutable properties matches the on-chain hash
	if !bytes.Equal(calculatedBidJSONHash, bidHash) {
		return fmt.Errorf("hash %x for bid JSON %s does not match hash in session: %x",
			calculatedBidJSONHash,
			transientBidJSON,
			bidHash,
		)
	}

	// check 3; check hash of relealed bid matches hash of private bid that was
	// added earlier. This ensures that the bid has not changed since it
	// was added to the session

	bidders := sessionJSON.PrivateBids
	privateBidHashString := bidders[bidKey].Hash

	onChainBidHashString := fmt.Sprintf("%x", bidHash)
	if privateBidHashString != onChainBidHashString {
		return fmt.Errorf("hash %s for bid JSON %s does not match hash in session: %s, bidder must have changed bid",
			privateBidHashString,
			transientBidJSON,
			onChainBidHashString,
		)
	}

	// we can add the bid to the session if all checks have passed
	type transientBidInput struct {
		BidType	 BidType `json:"bidType"`
		Volume   int    `json:"volume"`
		Org      string `json:"org"`
		Bidder   string `json:"bidder"`
	}

	// unmarshal bid imput
	var bidInput transientBidInput
	err = json.Unmarshal(transientBidJSON, &bidInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

	// marshal transient parameters and ID and MSPID into bid object
	NewBid := FullBid{
		BidType:  bidInput.BidType,
		Volume:   bidInput.Volume,
		Org:      bidInput.Org,
		Bidder:   bidInput.Bidder,
		Status:	  "Finalized",
	}

	// check 4: make sure that the session is being submitted is the bidder
	if bidInput.Bidder != clientID {
		return fmt.Errorf("Permission denied, client id %v is not the owner of the bid", clientID)
	}

	finalizedBids := make(map[string]FullBid)
	finalizedBids = sessionJSON.FinalizedBids
	finalizedBids[bidKey] = NewBid
	sessionJSON.FinalizedBids = finalizedBids

	newTransactionBytes, _ := json.Marshal(sessionJSON)

	// put session with bid added back into state
	err = ctx.GetStub().PutState(sessionID, newTransactionBytes)
	if err != nil {
		return fmt.Errorf("failed to update session: %v", err)
	}

	return nil
}

// CloseSession can be used by the admin to close the session. This prevents
// bids from being added to the session, and allows users to reveal their bid
func (s *SmartContract) CloseSession(ctx contractapi.TransactionContextInterface, sessionID string) error {

	sessionBytes, err := ctx.GetStub().GetState(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session %v: %v", sessionID, err)
	}

	if sessionBytes == nil {
		return fmt.Errorf("Session interest object %v not found", sessionID)
	}

	var sessionJSON Session
	err = json.Unmarshal(sessionBytes, &sessionJSON)
	if err != nil {
		return fmt.Errorf("failed to create session object JSON: %v", err)
	}

	// the session can only be Finalized by the admin

	// get ID of submitting client
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

	Admin := sessionJSON.Admin
	if Admin != clientID {
		return fmt.Errorf("session can only be Finalized by admin: %v", err)
	}

	Status := sessionJSON.Status
	if Status != "Open" {
		return fmt.Errorf("cannot close session that is not open")
	}

	sessionJSON.Status = string("Close")

	FinalizedTransaction, _ := json.Marshal(sessionJSON)

	err = ctx.GetStub().PutState(sessionID, FinalizedTransaction)
	if err != nil {
		return fmt.Errorf("failed to close session: %v", err)
	}

	return nil
}

// EndSession both changes the session status to Finalized and calculates the winners
// of the session
func (s *SmartContract) EndSession(ctx contractapi.TransactionContextInterface, sessionID string) error {

	sessionBytes, err := ctx.GetStub().GetState(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session %v: %v", sessionID, err)
	}

	if sessionBytes == nil {
		return fmt.Errorf("Session interest object %v not found", sessionID)
	}

	var sessionJSON Session
	err = json.Unmarshal(sessionBytes, &sessionJSON)
	if err != nil {
		return fmt.Errorf("failed to create session object JSON: %v", err)
	}

	// Check that the session is being ended by the admin

	// get ID of submitting client
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

	Admin := sessionJSON.Admin
	if Admin != clientID {
		return fmt.Errorf("session can only be ended by admin: %v", err)
	}

	Status := sessionJSON.Status
	if Status != "Close" {
		return fmt.Errorf("Can only end a closed session")
	}

	// get the list of revealed bids
	revealedBidMap := sessionJSON.FinalizedBids
	if len(sessionJSON.FinalizedBids) == 0 {
		return fmt.Errorf("No bids have been revealed, cannot end session: %v", err)
	}

	var totalSell, totalBuy int
	// approve or decline bids
	for _, bid := range revealedBidMap {
		if bid.BidType == "sell" || bid.BidType == "Sell" {
			totalSell += bid.Volume
		}

		if bid.BidType == "buy" || bid.BidType == "Buy" {
			totalBuy += bid.Volume
		}
	}

	for bidKey, bid := range revealedBidMap {
		if bid.BidType == "sell" || bid.BidType == "Sell" {
			if bid.Volume < totalBuy {
				totalBuy = totalBuy - bid.Volume
				err = UpdateStatus(ctx,sessionID,sessionJSON, bid,"Approved",bidKey)
				if err != nil {
					return fmt.Errorf("failed to update session: %v", err)
				}			
			}else if totalBuy != 0 {
				totalBuy -= totalBuy
				err = UpdateStatus(ctx,sessionID,sessionJSON, bid,"Partially Approved",bidKey)
				if err != nil {
					return fmt.Errorf("failed to update session: %v", err)
				}
			}else {
				err = UpdateStatus(ctx,sessionID,sessionJSON, bid,"Not Approved",bidKey)
				if err != nil {
					return fmt.Errorf("failed to update session: %v", err)
				}
			}
		}
		if bid.BidType == "buy" || bid.BidType == "Buy" {
			if bid.Volume < totalSell {
				totalSell = totalSell - bid.Volume
				err = UpdateStatus(ctx,sessionID,sessionJSON, bid,"Approved",bidKey)
				if err != nil {
					return fmt.Errorf("failed to update session: %v", err)
				}				
			}else if totalSell != 0 {
				totalSell -= totalSell
				err = UpdateStatus(ctx,sessionID,sessionJSON, bid,"Partially Approved",bidKey)
				if err != nil {
					return fmt.Errorf("failed to update session: %v", err)
				}
			}else {
				err = UpdateStatus(ctx,sessionID,sessionJSON, bid,"Not Approved",bidKey)
				if err != nil {
					return fmt.Errorf("failed to update session: %v", err)
				}
			}
		}
	}

	sessionJSON.Status = string("ended")

	FinalizedTransaction, _ := json.Marshal(sessionJSON)

	err = ctx.GetStub().PutState(sessionID, FinalizedTransaction)
	if err != nil {
		return fmt.Errorf("failed to end session: %v", err)
	}
	return nil
}
