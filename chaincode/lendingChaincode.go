package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Loan struct represents a loan request
type Loan struct {
	LoanID           string   `json:"loanId"`
	BorrowerID       string   `json:"borrowerId"`
	LenderID         string   `json:"lenderId"`
	Amount           float64  `json:"amount"`
	InterestRate     float64  `json:"interestRate"`
	Duration         int      `json:"duration"`
	Status           string   `json:"status"` // Pending, Approved, Active, Repaid, Defaulted
	DisbursementDate string   `json:"disbursementDate"`
	RepaymentDue     float64  `json:"repaymentDue"`
	RemainingBalance float64  `json:"remainingBalance"`
	Collateral       string   `json:"collateral"`
	Defaulted        bool     `json:"defaulted"`
	AuditHistory     []string `json:"auditHistory"` // List of loan status changes
}

// SmartContract provides functions for managing loans
type SmartContract struct {
	contractapi.Contract
}

// TODO: Implement function to request a loan
func (s *SmartContract) RequestLoan(ctx contractapi.TransactionContextInterface, loanID, borrowerID string, amount float64, interestRate float64, duration int) error {
	// TODO: Add logic to create a new loan request and store it on the blockchain

	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return fmt.Errorf("loanID cannot be empty")
	}

	if borrowerID == "" {
		fmt.Println("borrowerID cannot be empty")
		return fmt.Errorf("borrowerID cannot be empty")
	}

	if amount <= 0.0 {
		fmt.Println("amount cannot be less then 0")
		return fmt.Errorf("amount cannot be less then 0")
	}

	if interestRate <= 0.0 {
		fmt.Println("interestRate cannot be less then 0")
		return fmt.Errorf("interestRate cannot be less then 0")
	}

	if duration <= 0 {
		fmt.Println("duration cannot be less then 0")
		return fmt.Errorf("duration cannot be less then 0")
	}

	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if existingLoanData != nil && err == nil {
		fmt.Println("loanID already exists")
		return fmt.Errorf("loanID already exists")
	} else if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return fmt.Errorf("error in getting world state %s", err.Error())
	}

	var auditHistory = make([]string, 0)

	var loanData Loan
	loanData.LoanID = loanID
	loanData.BorrowerID = borrowerID
	loanData.Amount = amount
	loanData.InterestRate = interestRate
	loanData.Duration = duration
	loanData.Status = "Pending"
	auditHistory = append(auditHistory, "Pending")
	loanData.AuditHistory = auditHistory

	loanDataInByte, err := json.Marshal(loanData)
	if err != nil {
		fmt.Printf("error in marshal loan data %s\n", err.Error())
		return fmt.Errorf("error in marshal loan data %s", err.Error())
	}

	err = ctx.GetStub().PutState(loanID, loanDataInByte)
	if err != nil {
		fmt.Printf("error in writing to world state %s\n", err.Error())
		return fmt.Errorf("error in writing to world state %s", err.Error())
	}

	return nil
}

// TODO: Implement function to approve a loan
func (s *SmartContract) ApproveLoan(ctx contractapi.TransactionContextInterface, loanID, lenderID string) error {
	// TODO: Add logic to approve a loan and update its status
	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return fmt.Errorf("loanID cannot be empty")
	}

	if lenderID == "" {
		fmt.Println("lenderID cannot be empty")
		return fmt.Errorf("lenderID cannot be empty")
	}

	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return fmt.Errorf("error in getting world state %s", err.Error())
	}

	var loanData Loan

	err = json.Unmarshal(existingLoanData, &loanData)
	if err != nil {
		fmt.Printf("error in unmarshal loan data %s\n", err.Error())
		return fmt.Errorf("error in unmarshal loan data %s", err.Error())
	}

	if loanData.Status == "Approved" {
		fmt.Println("loan already approved")
		return fmt.Errorf("loan already approved")
	}

	if loanData.BorrowerID == lenderID {
		fmt.Println("lenderID is same as borrowerID")
		return fmt.Errorf("lenderID is same as borrowerID")

	}

	loanData.Status = "Approved"
	loanData.AuditHistory = append(loanData.AuditHistory, "Approved")
	loanData.LenderID = lenderID

	loanDataInByte, err := json.Marshal(loanData)
	if err != nil {
		fmt.Printf("error in marshal loan data %s\n", err.Error())
		return fmt.Errorf("error in marshal loan data %s", err.Error())
	}

	err = ctx.GetStub().PutState(loanID, loanDataInByte)
	if err != nil {
		fmt.Printf("error in writing to world state %s\n", err.Error())
		return fmt.Errorf("error in writing to world state %s", err.Error())
	}

	return nil
}

func (s *SmartContract) DisburseLoan(ctx contractapi.TransactionContextInterface, loanID string, date string) error {
	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return fmt.Errorf("loanID cannot be empty")
	}

	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return fmt.Errorf("error in getting world state %s", err.Error())
	}

	var loanData Loan

	err = json.Unmarshal(existingLoanData, &loanData)
	if err != nil {
		fmt.Printf("error in unmarshal loan data %s\n", err.Error())
		return fmt.Errorf("error in unmarshal loan data %s", err.Error())
	}

	if loanData.Status != "Approved" {
		fmt.Println("loan is not approved")
		return fmt.Errorf("loan is not approved")
	}

	loanData.Status = "Active"
	loanData.RepaymentDue = loanData.Amount
	loanData.DisbursementDate = date
	loanData.AuditHistory = append(loanData.AuditHistory, "Active")

	loanDataInByte, err := json.Marshal(loanData)
	if err != nil {
		fmt.Printf("error in marshal loan data %s\n", err.Error())
		return fmt.Errorf("error in marshal loan data %s", err.Error())
	}

	err = ctx.GetStub().PutState(loanID, loanDataInByte)
	if err != nil {
		fmt.Printf("error in writing to world state %s\n", err.Error())
		return fmt.Errorf("error in writing to world state %s", err.Error())
	}

	return nil

}

// TODO: Implement function to repay a loan
func (s *SmartContract) RepayLoan(ctx contractapi.TransactionContextInterface, loanID string, amount float64) error {
	// TODO: Add logic to process a loan repayment and update the remaining balance

	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return fmt.Errorf("loanID cannot be empty")
	}

	if amount <= 0.0 {
		fmt.Println("amount cannot be less then 0")
		return fmt.Errorf("amount cannot be less then 0")
	}

	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return fmt.Errorf("error in getting world state %s", err.Error())
	}

	var loanData Loan

	err = json.Unmarshal(existingLoanData, &loanData)
	if err != nil {
		fmt.Printf("error in unmarshal loan data %s\n", err.Error())
		return fmt.Errorf("error in unmarshal loan data %s", err.Error())
	}

	if loanData.Status != "Active" {
		fmt.Println("loan is not active")
		return fmt.Errorf("loan is not active")
	}
	remainingAmount := loanData.Amount - amount
	loanData.RepaymentDue = remainingAmount
	loanData.RemainingBalance = remainingAmount

	if remainingAmount <= 0.0 {
		loanData.Status = "Repaid"
		loanData.AuditHistory = append(loanData.AuditHistory, "Repaid")
	}

	loanDataInByte, err := json.Marshal(loanData)
	if err != nil {
		fmt.Printf("error in marshal loan data %s\n", err.Error())
		return fmt.Errorf("error in marshal loan data %s", err.Error())
	}

	err = ctx.GetStub().PutState(loanID, loanDataInByte)
	if err != nil {
		fmt.Printf("error in writing to world state %s\n", err.Error())
		return fmt.Errorf("error in writing to world state %s", err.Error())
	}

	return nil

}

func (s *SmartContract) CheckLoanStatus(ctx contractapi.TransactionContextInterface, loanID string) (string, error) {
	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return "", fmt.Errorf("loanID cannot be empty")
	}
	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return "", fmt.Errorf("error in getting world state %s", err.Error())
	}

	var loanData Loan

	err = json.Unmarshal(existingLoanData, &loanData)
	if err != nil {
		fmt.Printf("error in unmarshal loan data %s\n", err.Error())
		return "", fmt.Errorf("error in unmarshal loan data %s", err.Error())
	}

	return loanData.Status, nil

}

func (s *SmartContract) MarkAsDefaulted(ctx contractapi.TransactionContextInterface, loanID string) error {
	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return fmt.Errorf("loanID cannot be empty")
	}
	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return fmt.Errorf("error in getting world state %s", err.Error())
	}

	var loanData Loan

	err = json.Unmarshal(existingLoanData, &loanData)
	if err != nil {
		fmt.Printf("error in unmarshal loan data %s\n", err.Error())
		return fmt.Errorf("error in unmarshal loan data %s", err.Error())
	}

	if loanData.Status != "Active" {
		fmt.Println("loan is not active")
		return fmt.Errorf("loan is not active")
	}

	loanData.Status = "Defaulted"
	loanData.Defaulted = true
	loanData.AuditHistory = append(loanData.AuditHistory, "Defaulted")

	loanDataInByte, err := json.Marshal(loanData)
	if err != nil {
		fmt.Printf("error in marshal loan data %s\n", err.Error())
		return fmt.Errorf("error in marshal loan data %s", err.Error())
	}

	err = ctx.GetStub().PutState(loanID, loanDataInByte)
	if err != nil {
		fmt.Printf("error in writing to world state %s\n", err.Error())
		return fmt.Errorf("error in writing to world state %s", err.Error())
	}

	return nil
}

func (s *SmartContract) AddCollateral(ctx contractapi.TransactionContextInterface, loanID string, collateral string) error {
	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return fmt.Errorf("loanID cannot be empty")
	}
	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return fmt.Errorf("error in getting world state %s", err.Error())
	}

	var loanData Loan

	err = json.Unmarshal(existingLoanData, &loanData)
	if err != nil {
		fmt.Printf("error in unmarshal loan data %s\n", err.Error())
		return fmt.Errorf("error in unmarshal loan data %s", err.Error())
	}

	if loanData.Status != "Active" {
		fmt.Println("loan is not active")
		return fmt.Errorf("loan is not active")
	}

	loanData.Collateral = collateral

	loanDataInByte, err := json.Marshal(loanData)
	if err != nil {
		fmt.Printf("error in marshal loan data %s\n", err.Error())
		return fmt.Errorf("error in marshal loan data %s", err.Error())
	}

	err = ctx.GetStub().PutState(loanID, loanDataInByte)
	if err != nil {
		fmt.Printf("error in writing to world state %s\n", err.Error())
		return fmt.Errorf("error in writing to world state %s", err.Error())
	}

	return nil
}

func (s *SmartContract) GetLoanHistory(ctx contractapi.TransactionContextInterface, loanID string) ([]string, error) {
	// TODO: Add logic to fetch loan details from the blockchain
	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return []string{}, fmt.Errorf("loanID cannot be empty")
	}

	mspId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		fmt.Printf("error in getting msp id %s\n", err.Error())
		return []string{}, fmt.Errorf("error in getting msp id %s", err.Error())
	}

	if mspId != "RBI" {
		fmt.Println("only rbi can invoke this fucntion")
		return []string{}, fmt.Errorf("only rbi invoke this fucntion")
	}

	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return []string{}, fmt.Errorf("error in getting world state %s", err.Error())
	}

	var loanData Loan

	err = json.Unmarshal(existingLoanData, &loanData)
	if err != nil {
		fmt.Printf("error in unmarshal loan data %s\n", err.Error())
		return []string{}, fmt.Errorf("error in unmarshal loan data %s", err.Error())
	}
	return loanData.AuditHistory, nil
}

// TODO: Implement function to query a loan by ID
func (s *SmartContract) QueryLoan(ctx contractapi.TransactionContextInterface, loanID string) (*Loan, error) {
	// TODO: Add logic to fetch loan details from the blockchain
	if loanID == "" {
		fmt.Println("loanID cannot be empty")
		return nil, fmt.Errorf("loanID cannot be empty")
	}
	existingLoanData, err := ctx.GetStub().GetState(loanID)
	if err != nil {
		fmt.Printf("error in getting world state %s\n", err.Error())
		return nil, fmt.Errorf("error in getting world state %s", err.Error())
	}

	var loanData Loan

	err = json.Unmarshal(existingLoanData, &loanData)
	if err != nil {
		fmt.Printf("error in unmarshal loan data %s\n", err.Error())
		return nil, fmt.Errorf("error in unmarshal loan data %s", err.Error())
	}
	return &loanData, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %v", err)
	}
}