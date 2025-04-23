export PATH=${PWD}/../fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=$PWD/../fabric-samples/config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="RBI"
export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/peers/peer0.RBI.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$PWD/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/users/Admin@RBI.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

echo "Create Loan Request"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/../fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n lendingChaincodeFinal --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/peers/peer0.RBI.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/HDFC.example.com/peers/peer0.HDFC.example.com/tls/ca.crt" -c '{"function":"RequestLoan","Args":["loan117", "bororrer117", "1000", "5", "365"]}'

sleep 2
echo "Approve Loan Request"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/../fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n lendingChaincodeFinal --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/peers/peer0.RBI.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/HDFC.example.com/peers/peer0.HDFC.example.com/tls/ca.crt" -c '{"function":"ApproveLoan","Args":["loan117", "lender117"]}'

sleep 2
echo "DisburseLoan Loan"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/../fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n lendingChaincodeFinal --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/peers/peer0.RBI.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/HDFC.example.com/peers/peer0.HDFC.example.com/tls/ca.crt" -c '{"function":"DisburseLoan","Args":["loan117", "31/03/2025"]}'

sleep 2
echo "Repay Loan"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/../fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n lendingChaincodeFinal --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/peers/peer0.RBI.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/HDFC.example.com/peers/peer0.HDFC.example.com/tls/ca.crt" -c '{"function":"RepayLoan","Args":["loan117", "300"]}'

sleep 2
echo "Mark As defaulted Loan"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/../fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n lendingChaincodeFinal --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com/peers/peer0.RBI.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/../fabric-samples/test-network/organizations/peerOrganizations/HDFC.example.com/peers/peer0.HDFC.example.com/tls/ca.crt" -c '{"function":"MarkAsDefaulted","Args":["loan117"]}'