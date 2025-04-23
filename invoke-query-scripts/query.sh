export PATH=${PWD}/../fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=$PWD/../fabric-samples/config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="RBI"
export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/peers/peer0.RBI.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$PWD/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/users/Admin@RBI.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

echo "Query Loan ID"
peer chaincode query -C mychannel -n lendingChaincodeFinal -c '{"Args":["QueryLoan", "loan117"]}'