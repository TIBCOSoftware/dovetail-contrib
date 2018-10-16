package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Trigger implements a simple chaincode to invoke a trigger
type Trigger struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *Trigger) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *Trigger) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	_, args := stub.GetFunctionAndParameters()

	result, err := transaction.Invoke(stub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	// Return the result as success payload
	return shim.Success(result)
}

func main() {
	if err := shim.Start(new(Trigger)); err != nil {
		fmt.Printf("Error starting Trigger transaction: %s", err)
	}
}
