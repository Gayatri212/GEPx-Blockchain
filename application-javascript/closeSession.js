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

async function closeSession(ccp,wallet,user,sessionID) {
    try {

        const gateway = new Gateway();
      //connect using Discovery enabled

      await gateway.connect(ccp,
          { wallet: wallet, identity: user, discovery: { enabled: true, asLocalhost: true } });

        const network = await gateway.getNetwork(myChannel);
        const contract = network.getContract(myChaincodeName);

        // Query the session to get the list of endorsing orgs.
        //console.log('\n--> Evaluate Session: query the session you want to close');
        let sessionString = await contract.evaluateTransaction('QuerySession',sessionID);
        //console.log('*** Result:  Bid: ' + prettyJSONString(sessionString.toString()));
        var sessionJSON = JSON.parse(sessionString);

        let statefulTxn = contract.createTransaction('CloseSession');

        if (sessionJSON.organizations.length == 2) {
            statefulTxn.setEndorsingOrganizations(sessionJSON.organizations[0],sessionJSON.organizations[1]);
        } else {
            statefulTxn.setEndorsingOrganizations(sessionJSON.organizations[0]);
            }

        console.log('\n--> Submit Session: close session');
        await statefulTxn.submit(sessionID);
        console.log('*** Result: committed');

        console.log('\n--> Evaluate Session: query the updated session');
        let result = await contract.evaluateTransaction('QuerySession',sessionID);
        console.log('*** Result: Session: ' + prettyJSONString(result.toString()));

        gateway.disconnect();
    } catch (error) {
        console.error(`******** FAILED to submit bid: ${error}`);
        process.exit(1);
	}
}

async function main() {
    try {

        if (process.argv[2] == undefined || process.argv[3] == undefined
            || process.argv[4] == undefined) {
            console.log("Usage: node closeSession.js org userID sessionID");
            process.exit(1);
        }

        const org = process.argv[2]
        const user = process.argv[3];
        const sessionID = process.argv[4];

        if (org == 'Org1' || org == 'org1') {

            const orgMSP = 'Org1MSP';
            const ccp = buildCCPOrg1();
            const walletPath = path.join(__dirname, 'wallet/org1');
            const wallet = await buildWallet(Wallets, walletPath);
            await closeSession(ccp,wallet,user,sessionID);
        }
        else if (org == 'Org2' || org == 'org2') {

            const orgMSP = 'Org2MSP';
            const ccp = buildCCPOrg2();
            const walletPath = path.join(__dirname, 'wallet/org2');
            const wallet = await buildWallet(Wallets, walletPath);
            await closeSession(ccp,wallet,user,sessionID);
        }  else {
            console.log("Usage: node closeSession.js org userID sessionID ");
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
