package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// AssetsChaincode creates assets
type AssetsChaincode struct {
}

func main() {
	err := shim.Start(new(AssetsChaincode))
	if err != nil {
		fmt.Printf("Error starting AssetsChaincode chaincode: %s", err)
	}
}

// Init is where initialization and resets should happen
func (t *AssetsChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Constructs an asset name. Concantenates the asset name and the transaction id
func (t *AssetsChaincode) newAssetName(assetName string, stub shim.ChaincodeStubInterface) string {
	return fmt.Sprintf("%s-%s", assetName, stub.GetTxID())
}

// createAsset creates an asset on the ledger
func (t *AssetsChaincode) createAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("invalid number of arguments. Expect asset name and asset data")
	}

	err := stub.PutState(t.newAssetName(args[0], stub), []byte(args[1]))
	if err != nil {
		return nil, errors.New("Failed to create asset -> " + args[0])
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *AssetsChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("invoke is running " + function)

	switch function {
	case "init":
		return t.Init(stub, "init", args)

	case "create":
		return t.createAsset(stub, args)
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Get an asset by it's name
func (t *AssetsChaincode) getAsset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("invalid number of arguments. Expect asset name")
	}
	return stub.GetState(args[0])
}

// Query is our entry point for queries
func (t *AssetsChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	switch function {

	case "asset":
		assetDataBytes, err := t.getAsset(stub, args)
		if err != nil && len(assetDataBytes) == 0 {
			return nil, errors.New("asset not found")
		}
		return assetDataBytes, err

	default:
		fmt.Println("query did not find func: " + function)
	}

	return nil, errors.New("Received unknown function query: " + function)
}
