# Citopia starter kit
This repository contains Citopia chaincodes, start up instructions and server example


## Getting started

1. Create an AWS account
2. Send you AWS account ID to Citopia administrator
3. Accept AWS managed blockchain join proposal
4. Go through AWS Managed Blockchain set up steps from 2 to 5 (https://docs.aws.amazon.com/managed-blockchain/latest/managementguide/get-started-create-endpoint.html)
5. Copy and send to Citopia administrators the following files
    ```
    /admin-msp/admincerts
    /admin-msp/cacerts
    ```
6. Install chaincodes
    ```
    cd chaincode
    sh installAllChaincodes.sh VERSION (e.g. 1.0.0)
    ```
    wait for chaincodes getting installed and then run instantiate script
    ```
    sh instantiateAllChaincodes.sh VERSION (e.g. 1.0.0)
    ```
    after chaincode been updated, increase chaincode version and run update scrips
    ```
    sh upgradeAllChaincodes.sh VERSION (e.g. 1.0.1)
    ```
    
    If you want to install a new chaincode, you can add a new line in the end of these scripts
     (e.g. `install "my-contract"`) or use a separate scrips:
     
     ```
     sh installChaincode.sh NAME (e.g mycc) VERSION (e.g. 1.0.0)
     sh instantiateChaincode.sh NAME (e.g mycc) VERSION (e.g. 1.0.0)
     sh upgradeChaincode.sh NAME (e.g mycc) VERSION (e.g. 1.0.1)
     ```
    
7. Build your own server (see `/server` for NodeJS server example)