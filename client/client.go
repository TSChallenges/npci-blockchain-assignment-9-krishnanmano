/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	// "github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/status"
)

const (
	mspID        = "RBI"
	cryptoPath   = "../fabric-samples/test-network/organizations/peerOrganizations/RBI.example.com"
	certPath     = cryptoPath + "/users/User1@RBI.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@RBI.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.RBI.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.RBI.example.com"
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	chaincodeName := "lendingChaincodeFinal"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	createLoanRequest(contract)
	readLoanByID(contract)
	
	approveLoanRequest(contract)
	readLoanByID(contract)
	
	disburseLoan(contract)
	readLoanByID(contract)
	
	repayLoan(contract)
	readLoanByID(contract)
	
	defaultedLoan(contract)
	readLoanByID(contract)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

func createLoanRequest(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: createLoanRequest, creates new loanRequest\n")

	_, err := contract.SubmitTransaction("RequestLoan", "loan0001", "borrower0001", "1000", "5", "365")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func approveLoanRequest(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: approveLoanRequest, approves loanRequest\n")

	_, err := contract.SubmitTransaction("ApproveLoan", "loan0001", "lender0001")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func repayLoan(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: repayLoan, repay loan\n")

	_, err := contract.SubmitTransaction("RepayLoan", "loan0001", "300")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func disburseLoan(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: disburseLoan\n")

	_, err := contract.SubmitTransaction("DisburseLoan", "loan0001", "31/03/2025")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func defaultedLoan(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: defaultedLoan\n")

	_, err := contract.SubmitTransaction("MarkAsDefaulted", "loan0001")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func readLoanByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: readLoanByID, function returns loan attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("QueryLoan", "loan0001")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := string(evaluateResult)

	fmt.Printf("*** Query Loan by Id result: %s\n", result)
}