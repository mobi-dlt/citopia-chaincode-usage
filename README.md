# Citopia API

This repository contains Citopia chaincodes, 
start up instructions and chaincode API.

## Getting started

1. Create an AWS account
2. Send you AWS account ID to Citopia administrator
3. Accept AMB (Amazon Managed Blockchain) join proposal
4. Go through [AMB setup steps](https://docs.aws.amazon.com/managed-blockchain/latest/managementguide/get-started-create-endpoint.html) from 2 to 5
5. Copy and send to Citopia administrators the following files
    ```
    /admin-msp/admincerts
    /admin-msp/cacerts
    ```
6. Clone this repository on your Amazon EC2 instance:

   ```
   git clone https://github.com/mobi-dlt/citopia-chaincodes.git 
   ```

7. Install and run chaincode on your Amazon EC2 instance.

    Run the following command to install chaincode on a peer node:
    
    ```sh
    docker exec -e "CORE_PEER_TLS_ENABLED=true" \
    -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/home/managedblockchain-tls-chain.pem" \
    -e "CORE_PEER_LOCALMSPID=$MSP" \
    -e "CORE_PEER_MSPCONFIGPATH=$MSP_PATH" \
    -e "CORE_PEER_ADDRESS=$PEER" \
    cli peer chaincode install \
    -n trip-contract -v 1.0.0 -p github.com/chaincode/trip-contract
    ```
    
    Run the following command to instantiate the chaincode:
    
    ```sh
    docker exec -e "CORE_PEER_TLS_ENABLED=true" \
    -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/home/managedblockchain-tls-chain.pem" \
    -e "CORE_PEER_LOCALMSPID=$MSP" \
    -e "CORE_PEER_MSPCONFIGPATH=$MSP_PATH" \
    -e "CORE_PEER_ADDRESS=$PEER" \
    cli peer chaincode instantiate \
    -o $ORDERER -C citopia-channel -n trip-contract -v 1.0.0 \
    -c '{"Args":[]}' \
    --cafile /opt/home/managedblockchain-tls-chain.pem --tls
    ```
    
    You may have to wait a minute or two for the instantiation to propagate to the peer node. 
    Use the following command to verify instantiation:
    
    ```sh
    docker exec -e "CORE_PEER_TLS_ENABLED=true" \
    -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/home/managedblockchain-tls-chain.pem" \
    -e "CORE_PEER_LOCALMSPID=$MSP" \
    -e  "CORE_PEER_MSPCONFIGPATH=$MSP_PATH" \
    -e "CORE_PEER_ADDRESS=$PEER"  \
    cli peer chaincode list --instantiated \
    -o $ORDERER -C citopia-channel \
    --cafile /opt/home/managedblockchain-tls-chain.pem --tls
    ```
    
    After making changes in the chaincode you must upgrade it using the following command:
    
8. Upgrade the chaincode

    In case of changes to the contract logic by Citopia (for example, adding new properties to the model),
     you will also need to update your version of the contract. To do this, run the following command
    
    ```sh
    docker exec -e "CORE_PEER_TLS_ENABLED=true" \
    -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/home/managedblockchain-tls-chain.pem" \
    -e "CORE_PEER_LOCALMSPID=$MSP" \
    -e "CORE_PEER_MSPCONFIGPATH=$MSP_PATH" \
    -e "CORE_PEER_ADDRESS=$PEER" \
    cli peer chaincode upgrade \
    -o $ORDERER -C citopia-channel -n trip-contract -v 1.0.0 \
    -c '{"Args":[]}' \
    --cafile /opt/home/managedblockchain-tls-chain.pem --tls
    ```
    
9. Build your own server (see `/nodejs-server-example`)