package logger

// Imports
import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"

	dthlf "github.com/TIBCOSoftware/dovetail-contrib/blockchain/hyperledger-fabric"
)

// activityMetadata is the metadata of the activity as described in activity.json
// We'll store it as a variable to reuse it across multiple testcases
var activityMetadata *activity.Metadata

type TestChainCode struct{}

func (cc *TestChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (cc *TestChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke TestChainCode")
	return shim.Success(nil)
}

// getActivityMetadata reads the activity.json file and sets the activityMetadata variable
// if the variable already contains metadata it simply returns the current value rather than reading the file again
func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("../../../SmartContract/activity/logger/activity.json")
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

	shim.SetLoggingLevel(shim.LogDebug)
	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)

	tc.SetInput(ivLogLevel, "INFO")
	tc.SetInput(ivMessage, "this is testing")
	tc.SetInput("containerServiceStub", containerServiceStub)
	// Execute the activity
	_, err := act.Eval(tc)

	// We assume there will be no errors
	assert.Nil(t, err)
}
