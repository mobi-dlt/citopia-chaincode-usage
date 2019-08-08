# Citopia Chaincode Usage

This repository contains Citopia chaincodes, 
start up instructions and chaincode API.

### Setup AWS account

1. Create an AWS account
2. Send you AWS account ID to Citopia administrator
3. Accept AMB (Amazon Managed Blockchain) join proposal
4. Go through [AMB setup steps](https://docs.aws.amazon.com/managed-blockchain/latest/managementguide/get-started-create-endpoint.html) from 2 to 5
5. Copy and send to Citopia administrators the following files
    ```
    /admin-msp/admincerts
    /admin-msp/cacerts
    ```
6. Clone current repository on your Amazon EC2 instance:

   ```
   git clone https://github.com/mobi-dlt/citopia-chaincodes.git 
   ```

### Install chaincode on EC2 instance

Run the following command to install chaincode on a peer node:

```bash
docker exec -e "CORE_PEER_TLS_ENABLED=true" \
-e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/home/managedblockchain-tls-chain.pem" \
-e "CORE_PEER_LOCALMSPID=$MSP" \
-e "CORE_PEER_MSPCONFIGPATH=$MSP_PATH" \
-e "CORE_PEER_ADDRESS=$PEER" \
cli peer chaincode install \
-n trip-contract -v 1.0.0 -p github.com/chaincode/trip-contract
```

### Query the Chaincode

You may need to wait a brief moment for the chaincode instantiation to complete before you run
 the following command to query a value:
 
```bash
docker exec -e "CORE_PEER_TLS_ENABLED=true" \
-e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/home/managedblockchain-tls-chain.pem" \
-e "CORE_PEER_ADDRESS=$PEER" \
-e "CORE_PEER_LOCALMSPID=$MSP" \
-e "CORE_PEER_MSPCONFIGPATH=$MSP_PATH" \
cli peer chaincode query -C citopia-channel \
-n trip-contract -c '{"function":"findTrips","Args":[]}'
```

### Chaincode API

#### Trip contract

 * `findTrip(args[])` - find trip by given id
 
 Parameters: 
 
 `args[0]` - trip id

 * `findTrips(args[])` - find trips by filter
 
 Parameters: 

`args[0]` - stringified filter

`filter`:
 * `providerId: string` - current provider id (**required**)
 * `userId: string` - trip user id
 * `serviceId: string` - user service id
 * `serviceTypes: string[]` - possible values: "type-a", "type-b", "type-c"
 * `completed: number` - `0` for `false` or `1` for `true`
 * `paid: number` - `0` for `false` or `1` for `true`

    
### API usage example
 
See node.js usage example [here](./nodejs-server-example).
