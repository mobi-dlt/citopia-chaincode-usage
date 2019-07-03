import {ChannelEventHub} from "fabric-client";
import * as config from "../config/connection.json"

const FabricClient = require("fabric-client");

export async function invokeChaincode(
  channelName: string,
  username: string,
  chaincodeName: string,
  fcn: string,
  args: string[]
) {
  let errorMessage: string | null = null;
  let txIdAsString: string = "";
  try {
    // Connecting to the channel
    const client = FabricClient.loadFromConfig(config);
    await client.initCredentialStores();
    await client.getUserContext(username, true);
    const channel = client.getChannel(channelName);
    const txId = client.newTransactionID(true);
    txIdAsString = txId.getTransactionID();

    // Send transaction proposal
    let results = await channel.sendTransactionProposal({
      targets: channel.getPeers(),
      chaincodeId: chaincodeName,
      fcn: fcn,
      args: args,
      txId: txId
    });

    // the returned object has both the endorsement results
    // and the actual proposal, the proposal will be needed
    // later when we send a transaction to the ordering service
    const proposalResponses = results[0];
    const proposal = results[1];

    // lets have a look at the responses to see if they are
    // all good, if good they will also include signatures
    // required to be committed
    let successfulResponses = true;
    for (const i in proposalResponses) {
      let oneSuccessfulResponse = false;
      if (
        proposalResponses &&
        proposalResponses[i].response &&
        proposalResponses[i].response.status === 200
      ) {
        oneSuccessfulResponse = true;
      } else {
        console.error(
          "##### invokeChaincode - received unsuccessful proposal response: ",
          proposalResponses
        );
      }
      successfulResponses = successfulResponses && oneSuccessfulResponse;
    }

    if (successfulResponses) {
      // wait for the channel-based event hub to tell us
      // that the commit was good or bad on each peer in our organization
      const promises = [];
      let eventHubs = channel.getChannelEventHubsForOrg();
      eventHubs.forEach((eventHub: ChannelEventHub) => {
        let invokeEventPromise = new Promise((resolve, reject) => {
          let eventTimeout = setTimeout(() => {
            let message = "REQUEST_TIMEOUT:" + eventHub.getPeerAddr();
            console.error(message);
            eventHub.disconnect();
          }, 3000);
          eventHub.registerTxEvent(
            txIdAsString,
            (tx: any, code: any, block_num: any) => {
              console.info(
                "##### invokeChaincode - The invoke chaincode transaction has been committed on peer %s",
                eventHub.getPeerAddr()
              );
              console.info(
                "##### invokeChaincode - Transaction %s has status of %s in block %s",
                tx,
                code,
                block_num
              );
              clearTimeout(eventTimeout);

              if (code !== "VALID") {
                let message = `##### invokeChaincode - The invoke chaincode transaction was invalid, code: ${code}`;
                console.error(message);
                reject(new Error(message));
              } else {
                let message =
                  "##### invokeChaincode - The invoke chaincode transaction was valid.";
                console.info(message);
                resolve(message);
              }
            },
            err => {
              clearTimeout(eventTimeout);
              console.error(err);
              reject(err);
            },
            // the default for 'unregister' is true for transaction listeners
            // so no real need to set here, however for 'disconnect'
            // the default is false as most event hubs are long running
            // in this use case we are using it only once
            { unregister: true, disconnect: true }
          );
          eventHub.connect();
        });
        promises.push(invokeEventPromise);
      });

      const ordererRequest = {
        txId: txId,
        proposalResponses: proposalResponses,
        proposal: proposal
      };
      const sendPromise = channel.sendTransaction(ordererRequest);
      // put the send to the ordering service last so that the events get registered and
      // are ready for the orderering and committing
      promises.push(sendPromise);
      let results = await Promise.all(promises);
      let response = results.pop(); //  ordering service results are last in the results
      if (response.status === "SUCCESS") {
        console.info(
          "##### invokeChaincode - Successfully sent transaction to the ordering service."
        );
      } else {
        console.info(
          `##### invokeChaincode - Failed to order the transaction. Error code: ${
            response.status
          }`
        );
      }

      // now see what each of the event hubs reported
      for (let i in results) {
        let eventHubResult = results[i];
        let eventHub = eventHubs[i];
        console.info(
          "##### invokeChaincode - Event results for event hub :%s",
          eventHub.getPeerAddr()
        );
        if (typeof eventHubResult === "string") {
          console.info("##### invokeChaincode - " + eventHubResult);
        } else {
          if (!errorMessage) errorMessage = eventHubResult.toString();
          console.info("##### invokeChaincode - " + eventHubResult.toString());
        }
      }
    } else {
      console.info(
        "##### invokeChaincode - Failed to send Proposal and receive all good ProposalResponse"
      );
    }
  } catch (error) {
    console.error(
      "##### invokeChaincode - Failed to invoke due to error: " + error.stack
        ? error.stack
        : error
    );
    errorMessage = error.toString();
  }

  if (!errorMessage) {
    console.info("##### invokeChaincode - Successfully invoked chaincode");
    return { transactionId: txIdAsString };
  } else {
    const message = `##### invokeChaincode - Failed to invoke chaincode. cause: ${errorMessage}`;
    console.error(message);
    throw new Error(message);
  }
}

export async function queryChaincode(
  channelName: string,
  username: string,
  chaincodeName: string,
  fcn: string,
  args: string[]
) {
  const client = FabricClient.loadFromConfig(config);
  await client.initCredentialStores();
  await client.getUserContext(username, true);
  const channel = client.getChannel(channelName);

  const request = {
    targets: channel.getPeers(),
    chaincodeId: chaincodeName,
    fcn: fcn,
    args: args
  };
  const responses = await channel.queryByChaincode(request);
  const response = responses[0].toString();
  if (response.indexOf("Error") !== -1) {
    throw new Error(`error in query result: ${response}`);
  }
  return response;
}
