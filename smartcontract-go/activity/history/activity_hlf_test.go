package history

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

	dthlf "github.com/TIBCOSoftware/dovetail-contrib/blockchain/hyperledger-fabric"
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
		jsonMetadataBytes, err := ioutil.ReadFile("../../../SmartContract/activity/history/activity.json")
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
func put(stub *shim.MockStub) {
	shim.SetLoggingLevel(shim.LogDebug)
	assetName := "com.tibco.cp.iou"
	key, _ := stub.CreateCompositeKey(assetName, []string{"rec1"})
	value := `{"lender":"Alice","Charlie","amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"rec1"}`
	stub.PutState(key, []byte(value))
}

func putComposite(stub *shim.MockStub) {

	assetName := "com.tibco.cp.iou"
	value := `{"lender":"Alice","Charlie","amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"rec1"}`
	key, _ := stub.CreateCompositeKey(assetName, []string{"Alice", "Charlie"})
	stub.PutState(key, []byte(value))
}

func TestEvalHistory(t *testing.T) {
	fmt.Println("--------GET--------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "linearId"

	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)
	hlfStub.MockTransactionStart("issueIOU")
	shim.SetLoggingLevel(shim.LogDebug)

	put(hlfStub)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("containerServiceStub", containerServiceStub)

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	inputmp["Value"] = `{"linearId":"rec1"}`
	v, _ := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("err %v\n", err)
		return
	}

	// We assume there will be no errors

	complexvalue := tc.GetOutput("output").(*data.ComplexObject)
	values := complexvalue.Value
	switch t := values.(type) {
	case string:
		fmt.Printf("string value %s \n", t)
		break
	case []interface{}:
		for _, v := range t {
			if v != nil {
				fmt.Printf("state=%v\n", v)
			}
		}
	}

	hlfStub.MockTransactionStart("issueIOU")
}
