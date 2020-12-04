/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const { buildCCPOrg1, buildCCPOrg2, buildWallet } = require('../../test-application/javascript/AppUtil.js');

const myChannel = 'mychannel';
const myChaincodeName = 'gepx';


function prettyJSONString(inputString) {
    if (inputString) {
        return JSON.stringify(JSON.parse(inputString), null, 2);
    }
    else {
        return inputString;
    }
}

async function submitBid(ccp,wallet,user,sessionID,bidID) {
    try {

        const gateway = new Gateway();
      //connect using Discovery enabled

      await gateway.connect(ccp,
          { wallet: wallet, identity: user, discovery: { enabled: true, asLocalhost: true } });

        const network = await gateway.getNetwork(myChannel);
        const contract = network.getContract(myChaincodeName);

        console.log('\n--> Evaluate Session: query the session you want to join');
        let sessionString = await contract.evaluateTransaction('QuerySession',sessionID);
        var sessionJSON = JSON.parse(sessionString);

        let statefulTxn = contract.createTransaction('SubmitBid');

        if (sessionJSON.organizations.length == 2) {
            statefulTxn.setEndorsingOrganizations(sessionJSON.organizations[0],sessionJSON.organizations[1]);
            } else {
            statefulTxn.setEndorsingOrganizations(sessionJSON.organizations[0]);
            }

        console.log('\n--> Submit Session: add bid to the session');
        await statefulTxn.submit(sessionID,bidID);

        console.log('\n--> Evaluate Session: query the session to see that our bid was added');
        let result = await contract.evaluateTransaction('QuerySession',sessionID);
        console.log('*** Result: session: ' + prettyJSONString(result.toString()));

        gateway.disconnect();
    } catch (error) {
        console.error(`******** FAILED to submit bid: ${error}`);
		process.exit(1);
	}
}

async function main() {
    try {

        if (process.argv[2] == undefined || process.argv[3] == undefined
            || process.argv[4] == undefined || process.argv[5] == undefined) {
            console.log("Usage: node submitBid.js org userID sessionID bidID");
            process.exit(1);
        }

        const org = process.argv[2]
        const user = process.argv[3];
        const sessionID = process.argv[4];
        const bidID = process.argv[5];

        if (org == 'Org1' || org == 'org1') {

            const orgMSP = 'Org1MSP';
            const ccp = buildCCPOrg1();
            const walletPath = path.join(__dirname, 'wallet/org1');
            const wallet = await buildWallet(Wallets, walletPath);
            await submitBid(ccp,wallet,user,sessionID,bidID);
        }
        else if (org == 'Org2' || org == 'org2') {

            const orgMSP = 'Org2MSP';
            const ccp = buildCCPOrg2();
            const walletPath = path.join(__dirname, 'wallet/org2');
            const wallet = await buildWallet(Wallets, walletPath);
            await submitBid(ccp,wallet,user,sessionID,bidID);
        }
        else {
            console.log("Usage: node submitBid.js org userID sessionID bidID");
            console.log("Org must be Org1 or Org2");
          }
    } catch (error) {
		console.error(`******** FAILED to run the application: ${error}`);
    if (error.stack) {
        console.error(error.stack);
    }
    process.exit(1);
    }
}


main();
