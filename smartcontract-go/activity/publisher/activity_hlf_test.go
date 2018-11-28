/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package publisher

// Imports
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"

	dthlf "github.com/TIBCOSoftware/dovetail-contrib/container/hyperledgerfabric"
)

type TestChainCode struct{}

func (cc *TestChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *TestChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke TestChainCode")
	return shim.Success(nil)
}

// activityMetadata is the metadata of the activity as described in activity.json
// We'll store it as a variable to reuse it across multiple testcases
var activityMetadata *activity.Metadata

// getActivityMetadata reads the activity.json file and sets the activityMetadata variable
// if the variable already contains metadata it simply returns the current value rather than reading the file again
func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("../../../SmartContract/activity/publisher/activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}
		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}
	return activityMetadata
}

// TestActivityRegistration checks whether the activity can be registered, and is registered in the engine
func TestActivityRegistration(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

// TestEval tests the Eval function and sends a message to IFTTT
// Make sure that you have updated the values below
func TestEval(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	event := "com.tibco.cp.IOUIssued"
	eventMetadata := "meta data"
	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)
	hlfStub.MockTransactionStart("issueIOU")
	shim.SetLoggingLevel(shim.LogDebug)

	tc.SetInput(ivEvent, event)
	tc.SetInput(ivMetadata, eventMetadata)
	tc.SetInput("containerServiceStub", containerServiceStub)

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	inputmp["Value"] = `{"lender":"net.i2p.crypto.eddsa.EdDSAPublicKey@c3cefb34","borrower":"net.i2p.crypto.eddsa.EdDSAPublicKey@a7348296","amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_null"}`
	v, err := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	err = json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err = act.Eval(tc)
	// We assume there will be no errors
	assert.Nil(t, err)

	hlfStub.MockTransactionStart("issueIOU")
}
