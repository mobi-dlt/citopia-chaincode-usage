#!/bin/bash

CC_VERSION=${1}

install() {
    NAME=$1

    echo
    echo "============ Installing chaincode ${NAME} ============"
    echo
    docker exec -e "CORE_PEER_TLS_ENABLED=true" \
    -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/home/managedblockchain-tls-chain.pem" \
    -e "CORE_PEER_LOCALMSPID=$MSP" \
    -e "CORE_PEER_MSPCONFIGPATH=$MSP_PATH" \
    -e "CORE_PEER_ADDRESS=$PEER" \
    cli peer chaincode install \
    -n ${NAME} -v ${CC_VERSION} -p github.com/chaincode/src/${NAME}
}

install "trip-contract"
