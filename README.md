# GEPx-Blockchain
Hyperledger fabric POC

### Use Case - GE Power Exchange 

- This is simple power exchange model which demonstrates volume based power exchange settlement.

### Steps to run application

#### Pre-requisites 
- Go (version 1.15.x)
- NodeJS (v14.15.1)
- Docker
- Docker-compose

#### Installing Hyperledger Fabric

```
cd $GOPATH/src/github.com
sudo curl -sSL https://bit.ly/2ysbOFE | bash -s
```
This will install hyperledger fabric including [fabric-samples](https://github.com/hyperledger/fabric-samples) folder which we will use for creating test-network with CA.

#### Cloning this repository

```
cd fabric-samples
git clone https://github.com/Gayatri212/GEPx-Blockchain.git
```
Make sure you clone this repository inside fabric-samples folder

#### Create channel with CA

```
cd fabric-samples/test-network
./network.sh up createChannel -ca
```
This will create a channel to which we will add our orgs for POC

#### Deploy Smart-Contract on channel

```
./network.sh deployCC -ccn gepx -ccp ../GEPx-Blockchain/chaincode-go/ -ccep "OR('Org1MSP.peer','Org2MSP.peer')"
```
This will add org1 and org2 to the channel and deploy chaincode named gepx on them. By default it will chreate one peer per org.
```
2020-12-02 05:09:23.306 UTC [chaincodeCmd] ClientWait -> INFO 001 txid [8e1e1414b2f0e2d2891ea5f8258b9d814a7ae414437f666632fa1f299a505f39] committed with status (VALID) at localhost:7051
2020-12-02 05:09:23.314 UTC [chaincodeCmd] ClientWait -> INFO 002 txid [8e1e1414b2f0e2d2891ea5f8258b9d814a7ae414437f666632fa1f299a505f39] committed with status (VALID) at localhost:9051
Chaincode definition committed on channel 'mychannel'
Using organization 1
Querying chaincode definition on peer0.org1 on channel 'mychannel'...
Attempting to Query committed status on peer0.org1, Retry after 3 seconds.
+ peer lifecycle chaincode querycommitted --channelID mychannel --name gepx
+ res=0
Committed chaincode definition for chaincode 'gepx' on channel 'mychannel':
Version: 1.0, Sequence: 1, Endorsement Plugin: escc, Validation Plugin: vscc, Approvals: [Org1MSP: true, Org2MSP: true]
Query chaincode definition successful on peer0.org1 on channel 'mychannel'
Using organization 2
Querying chaincode definition on peer0.org2 on channel 'mychannel'...
Attempting to Query committed status on peer0.org2, Retry after 3 seconds.
+ peer lifecycle chaincode querycommitted --channelID mychannel --name gepx
+ res=0
Committed chaincode definition for chaincode 'gepx' on channel 'mychannel':
Version: 1.0, Sequence: 1, Endorsement Plugin: escc, Validation Plugin: vscc, Approvals: [Org1MSP: true, Org2MSP: true]
Query chaincode definition successful on peer0.org2 on channel 'mychannel'
Chaincode initialization is not required
```
Ending of logs will look like this

#### Running Application

```
cd fabric-samples/GEPx-Blockchain/application-javascript
npm install
```
This will install all the dependencies of application and you will get following in logs
```
found 0 vulnerabilities



   ╭────────────────────────────────────────────────────────────────╮
   │                                                                │
   │      New patch version of npm available! 6.14.8 → 6.14.9       │
   │   Changelog: https://github.com/npm/cli/releases/tag/v6.14.9   │
   │               Run npm install -g npm to update!                │
   │                                                                │
   ╰────────────────────────────────────────────────────────────────╯
```

1. Enroll organization
```
node enrollAdmin.js org1
node enrollAdmin.js org2
```
Logs
```
--> Enrolling the Org1 CA admin
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a CA Client named ca-org1
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1
Successfully enrolled admin user and imported it into the wallet
```

2. Register admin user and create transaction
```
node registerEnrollUser.js org1 adminuser
node createTransaction.js org1 adminuser tx1
```
Logs
```
$ node registerEnrollUser.js org1 adminuser

--> Register and enrolling new user
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a CA Client named ca-org1
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1
Successfully registered and enrolled user adminuser and imported it into the wallet

$ node createTransaction.js org1 adminuser tx1
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Submit Transaction: Propose a new transaction
*** Result: committed

--> Evaluate Transaction: query the transaction that was just created
*** Result: Transaction: {
  "admin": "eDUwOTo6Q049YWRtaW51c2VyLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "organizations": [
    "Org1MSP"
  ],
  "privateBids": {},
  "revealedBids": {},
  "status": "Open"
}
```

3. Create seller and buyers for bidding
```
$ node registerEnrollUser.js org1 seller1

--> Register and enrolling new user
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a CA Client named ca-org1
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1
Successfully registered and enrolled user seller1 and imported it into the wallet

$ node registerEnrollUser.js org1 buyer1

--> Register and enrolling new user
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a CA Client named ca-org1
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1
Successfully registered and enrolled user buyer1 and imported it into the wallet
```

4. Create and submit bids
```
$ node bid.js org1 seller1 tx1 100 sell
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Evaluate Transaction: get your client ID
*** Result:  Bidder ID is eDUwOTo6Q049c2VsbGVyMSxPVT1jbGllbnQrT1U9b3JnMStPVT1kZXBhcnRtZW50MTo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==

--> Submit Transaction: Create the bid that is stored in your organization's private data collection
*** Result: committed
*** Result ***SAVE THIS VALUE*** BidID: 2f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3

--> Evaluate Transaction: read the bid that was just created
*** Result:  Bid: {
  "bidType": "sell",
  "volume": 100,
  "org": "Org1MSP",
  "bidder": "eDUwOTo6Q049c2VsbGVyMSxPVT1jbGllbnQrT1U9b3JnMStPVT1kZXBhcnRtZW50MTo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==",
  "status": "Placed"
}

$ node bid.js org1 buyer1 tx1 90 buy
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Evaluate Transaction: get your client ID
*** Result:  Bidder ID is eDUwOTo6Q049YnV5ZXIxLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT

--> Submit Transaction: Create the bid that is stored in your organization's private data collection
*** Result: committed
*** Result ***SAVE THIS VALUE*** BidID: 07d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099

--> Evaluate Transaction: read the bid that was just created
*** Result:  Bid: {
  "bidType": "buy",
  "volume": 90,
  "org": "Org1MSP",
  "bidder": "eDUwOTo6Q049YnV5ZXIxLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "status": "Placed"
}

```

For submit bid use BidID generated by bid.js
```
$ node submitBid.js org1 seller1 tx1 2f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Evaluate Transaction: query the transaction you want to join

--> Submit Transaction: add bid to the transaction

--> Evaluate Transaction: query the transaction to see that our bid was added
*** Result: transaction: {
  "admin": "eDUwOTo6Q049YWRtaW51c2VyLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "organizations": [
    "Org1MSP"
  ],
  "privateBids": {
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "org": "Org1MSP",
      "hash": "fb3e86ae86ab36145c627e002aa8470555ce41cd95b6a80a6ab828d7c433fad0"
    }
  },
  "revealedBids": {},
  "status": "Open"
}

$ node submitBid.js org1 buyer1 tx1 07d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Evaluate Transaction: query the transaction you want to join

--> Submit Transaction: add bid to the transaction

--> Evaluate Transaction: query the transaction to see that our bid was added
*** Result: transaction: {
  "admin": "eDUwOTo6Q049YWRtaW51c2VyLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "organizations": [
    "Org1MSP"
  ],
  "privateBids": {
    "\u0000bid\u0000tx1\u000007d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099\u0000": {
      "org": "Org1MSP",
      "hash": "ef635f39a865cc71730a341b90d16277b7646fe5b1e34e02ca3af6db9e7e67ae"
    },
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "org": "Org1MSP",
      "hash": "fb3e86ae86ab36145c627e002aa8470555ce41cd95b6a80a6ab828d7c433fad0"
    }
  },
  "revealedBids": {},
  "status": "Open"
}

```

5. Stop Bidding -
Once all bids are placed admin user can stop the bidding process then users can finalize their bids.
```
$ node closeTransaction.js org1 adminuser tx1
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Submit Transaction: close transaction
*** Result: committed

--> Evaluate Transaction: query the updated transaction
*** Result: Transaction: {
  "admin": "eDUwOTo6Q049YWRtaW51c2VyLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "organizations": [
    "Org1MSP"
  ],
  "privateBids": {
    "\u0000bid\u0000tx1\u000007d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099\u0000": {
      "org": "Org1MSP",
      "hash": "ef635f39a865cc71730a341b90d16277b7646fe5b1e34e02ca3af6db9e7e67ae"
    },
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "org": "Org1MSP",
      "hash": "fb3e86ae86ab36145c627e002aa8470555ce41cd95b6a80a6ab828d7c433fad0"
    }
  },
  "revealedBids": {},
  "status": "Close"
}
```

6. Finalize bids
```
$ node revealBid.js org1 seller1 tx1 2f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Evaluate Transaction: read your bid
*** Result:  Bid: {
  "bidType": "sell",
  "volume": 100,
  "org": "Org1MSP",
  "bidder": "eDUwOTo6Q049c2VsbGVyMSxPVT1jbGllbnQrT1U9b3JnMStPVT1kZXBhcnRtZW50MTo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==",
  "status": "Placed"
}

--> Evaluate Transaction: query the transaction to see that our bid was added
*** Result: Transaction: {
  "admin": "eDUwOTo6Q049YWRtaW51c2VyLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "organizations": [
    "Org1MSP"
  ],
  "privateBids": {
    "\u0000bid\u0000tx1\u000007d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099\u0000": {
      "org": "Org1MSP",
      "hash": "ef635f39a865cc71730a341b90d16277b7646fe5b1e34e02ca3af6db9e7e67ae"
    },
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "org": "Org1MSP",
      "hash": "fb3e86ae86ab36145c627e002aa8470555ce41cd95b6a80a6ab828d7c433fad0"
    }
  },
  "revealedBids": {
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "bidType": "sell",
      "volume": 100,
      "org": "Org1MSP",
      "bidder": "eDUwOTo6Q049c2VsbGVyMSxPVT1jbGllbnQrT1U9b3JnMStPVT1kZXBhcnRtZW50MTo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==",
      "status": "Finalized"
    }
  },
  "status": "Close"
}


$ node revealBid.js org1 buyer1 tx1 07d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Evaluate Transaction: read your bid
*** Result:  Bid: {
  "bidType": "buy",
  "volume": 90,
  "org": "Org1MSP",
  "bidder": "eDUwOTo6Q049YnV5ZXIxLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "status": "Placed"
}

--> Evaluate Transaction: query the transaction to see that our bid was added
*** Result: Transaction: {
  "admin": "eDUwOTo6Q049YWRtaW51c2VyLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "organizations": [
    "Org1MSP"
  ],
  "privateBids": {
    "\u0000bid\u0000tx1\u000007d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099\u0000": {
      "org": "Org1MSP",
      "hash": "ef635f39a865cc71730a341b90d16277b7646fe5b1e34e02ca3af6db9e7e67ae"
    },
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "org": "Org1MSP",
      "hash": "fb3e86ae86ab36145c627e002aa8470555ce41cd95b6a80a6ab828d7c433fad0"
    }
  },
  "revealedBids": {
    "\u0000bid\u0000tx1\u000007d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099\u0000": {
      "bidType": "buy",
      "volume": 90,
      "org": "Org1MSP",
      "bidder": "eDUwOTo6Q049YnV5ZXIxLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
      "status": "Finalized"
    },
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "bidType": "sell",
      "volume": 100,
      "org": "Org1MSP",
      "bidder": "eDUwOTo6Q049c2VsbGVyMSxPVT1jbGllbnQrT1U9b3JnMStPVT1kZXBhcnRtZW50MTo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==",
      "status": "Finalized"
    }
  },
  "status": "Close"
}
```

6. End transaction - 
End transaction will update the status of bids approved/partially approved/declined based on smart contract. Only Finalized bids will be considered for approval.
```
$ node endTransaction.js org1 adminuser tx1
Loaded the network configuration located at /opt/go/src/github.com/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json
Built a file system wallet at /opt/go/src/github.com/fabric-samples/GEPx-Blockchain/application-javascript/wallet/org1

--> Submit the transaction to end the transaction
*** Result: committed

--> Evaluate Transaction: query the updated transaction
*** Result: Transaction: {
  "admin": "eDUwOTo6Q049YWRtaW51c2VyLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
  "organizations": [
    "Org1MSP"
  ],
  "privateBids": {
    "\u0000bid\u0000tx1\u000007d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099\u0000": {
      "org": "Org1MSP",
      "hash": "ef635f39a865cc71730a341b90d16277b7646fe5b1e34e02ca3af6db9e7e67ae"
    },
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "org": "Org1MSP",
      "hash": "fb3e86ae86ab36145c627e002aa8470555ce41cd95b6a80a6ab828d7c433fad0"
    }
  },
  "revealedBids": {
    "\u0000bid\u0000tx1\u000007d17392eeeaf3fbdd930aa1753765dca157a0b463bef88f581ed4b94e17f099\u0000": {
      "bidType": "buy",
      "volume": 90,
      "org": "Org1MSP",
      "bidder": "eDUwOTo6Q049YnV5ZXIxLE9VPWNsaWVudCtPVT1vcmcxK09VPWRlcGFydG1lbnQxOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPUR1cmhhbSxTVD1Ob3J0aCBDYXJvbGluYSxDPVVT",
      "status": "Aprroved"
    },
    "\u0000bid\u0000tx1\u00002f21a9738ef1bd651154a3015af8f566986b138ef453830aebad325e0c8e18b3\u0000": {
      "bidType": "sell",
      "volume": 100,
      "org": "Org1MSP",
      "bidder": "eDUwOTo6Q049c2VsbGVyMSxPVT1jbGllbnQrT1U9b3JnMStPVT1kZXBhcnRtZW50MTo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==",
      "status": "Partially Aprroved"
    }
  },
  "status": "ended"
}
```

#### Deleting Database
```
rm -rf wallet
```

#### Shut down network
```
cd fabric-samples/test-network
./network.sh down
```
