package main

import (
	SmartContract "MyContract/MySmartContract"
	_ "MyContract/MySmartContract"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	chaincode, err := contractapi.NewChaincode(
		new(SmartContract.SimpleContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}

}
