/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}

	if function == "write" {
		return t.write(stub, args)
	}

	if function == "create-table" {
		err := t.createUserTable(stub, args)
		if err != nil {
			fmt.Println("Failed to create table -> ", err.Error())
			return nil, errors.New("Failed to create table")
		}

		fmt.Println("Successfully created user table ")
		return []byte("successfully created"), nil
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) createUserTable(stub shim.ChaincodeStubInterface, args []string) error {

	// Define table column
	var userTableColumnDefs []*shim.ColumnDefinition
	nameColumnDef := shim.ColumnDefinition{Name: "name", Type: shim.ColumnDefinition_STRING, Key: true}
	ageColumnDef := shim.ColumnDefinition{Name: "age", Type: shim.ColumnDefinition_INT32, Key: false}
	genderColumnDef := shim.ColumnDefinition{Name: "gender", Type: shim.ColumnDefinition_INT32, Key: false}
	userTableColumnDefs = append(userTableColumnDefs, &nameColumnDef)
	userTableColumnDefs = append(userTableColumnDefs, &ageColumnDef)
	userTableColumnDefs = append(userTableColumnDefs, &genderColumnDef)

	// create table
	err := stub.CreateTable("user", userTableColumnDefs)

	return err
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string
	var err error
	fmt.Println("running write()")

	// tbl, _ := stub.GetTable("")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of variable and value to set")
	}

	name = args[0]
	value = args[0]
	err = stub.PutState(name, []byte(value))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {
		return t.read(stub, args)
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsBytes, err := stub.GetState(name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsBytes, nil
}
