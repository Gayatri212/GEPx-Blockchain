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

async function endTransaction(ccp,wallet,user,transactionID) {
    try {

        const gateway = new Gateway();
      //connect using Discovery enabled

      await gateway.connect(ccp,
          { wallet: wallet, identity: user, discovery: { enabled: true, asLocalhost: true } });

        const network = await gateway.getNetwork(myChannel);
        const contract = network.getContract(myChaincodeName);

        // Query the transaction to get the list of endorsing orgs.
        //console.log('\n--> Evaluate Transaction: query the transaction you want to end');
        let transactionString = await contract.evaluateTransaction('QueryTransaction',transactionID);
        //console.log('*** Result:  Bid: ' + prettyJSONString(transactionString.toString()));
        var transactionJSON = JSON.parse(transactionString);

        let statefulTxn = contract.createTransaction('EndTransaction');

        if (transactionJSON.organizations.length == 2) {
            statefulTxn.setEndorsingOrganizations(transactionJSON.organizations[0],transactionJSON.organizations[1]);
        } else {
            statefulTxn.setEndorsingOrganizations(transactionJSON.organizations[0]);
            }

        console.log('\n--> Submit the transaction to end the transaction');
        await statefulTxn.submit(transactionID);
        console.log('*** Result: committed');

        console.log('\n--> Evaluate Transaction: query the updated transaction');
        let result = await contract.evaluateTransaction('QueryTransaction',transactionID);
        console.log('*** Result: Transaction: ' + prettyJSONString(result.toString()));

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
            console.log("Usage: node endTransaction.js org userID transactionID");
            process.exit(1);
        }

        const org = process.argv[2]
        const user = process.argv[3];
        const transactionID = process.argv[4];

        if (org == 'Org1' || org == 'org1') {

            const orgMSP = 'Org1MSP';
            const ccp = buildCCPOrg1();
            const walletPath = path.join(__dirname, 'wallet/org1');
            const wallet = await buildWallet(Wallets, walletPath);
            await endTransaction(ccp,wallet,user,transactionID);
        }
        else if (org == 'Org2' || org == 'org2') {

            const orgMSP = 'Org2MSP';
            const ccp = buildCCPOrg2();
            const walletPath = path.join(__dirname, 'wallet/org2');
            const wallet = await buildWallet(Wallets, walletPath);
            await endTransaction(ccp,wallet,user,transactionID);
        }  else {
            console.log("Usage: node endTransaction.js org userID transactionID");
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
