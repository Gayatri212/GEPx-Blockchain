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

async function createTransaction(ccp,wallet,user,transactionID) {
    try {

        const gateway = new Gateway();
      //connect using Discovery enabled

      await gateway.connect(ccp,
          { wallet: wallet, identity: user, discovery: { enabled: true, asLocalhost: true } });

        const network = await gateway.getNetwork(myChannel);
        const contract = network.getContract(myChaincodeName);

        let statefulTxn = contract.createTransaction('CreateTransaction');

        console.log('\n--> Submit Transaction: Propose a new transaction');
        await statefulTxn.submit(transactionID);
        console.log('*** Result: committed');

        console.log('\n--> Evaluate Transaction: query the transaction that was just created');
        let result = await contract.evaluateTransaction('QueryTransaction',transactionID);
        console.log('*** Result: Transaction: ' + prettyJSONString(result.toString()));

        gateway.disconnect();
    } catch (error) {
        console.error(`******** FAILED to submit bid: ${error}`);
	}
}

async function main() {
    try {

        if (process.argv[2] == undefined || process.argv[3] == undefined
            || process.argv[4] == undefined) {
            console.log("Usage: node createTransaction.js org userID transactionID");
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
            await createTransaction(ccp,wallet,user,transactionID);
        }
        else if (org == 'Org2' || org == 'org2') {

            const orgMSP = 'Org2MSP';
            const ccp = buildCCPOrg2();
            const walletPath = path.join(__dirname, 'wallet/org2');
            const wallet = await buildWallet(Wallets, walletPath);
            await createTransaction(ccp,wallet,user,transactionID);
        }  else {
            console.log("Usage: node createTransaction.js org userID transactionID");
            console.log("Org must be Org1 or Org2");
          }
    } catch (error) {
		console.error(`******** FAILED to run the application: ${error}`);
    }
}


main();
