package ledger

// Imports
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	dthlf "github.com/TIBCOSoftware/dovetail-contrib/blockchain/hyperledger-fabric"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"
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
		jsonMetadataBytes, err := ioutil.ReadFile("../../../SmartContract/activity/ledger/activity.json")
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
	fmt.Println("--------put--------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "linearId"
	operation := "PUT"
	isArray := false
	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)
	hlfStub.MockTransactionStart("issueIOU")
	shim.SetLoggingLevel(shim.LogDebug)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("operation", operation)
	tc.SetInput("isArray", isArray)
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
	key, _ := hlfStub.CreateCompositeKey(assetName, []string{"ae846b5b-2189-46b4-bda7-8511ff33ba6b_null"})
	value, err := hlfStub.GetState(key)
	assert.Nil(t, err)
	fmt.Printf("state=%v\n", string(value))
	hlfStub.MockTransactionStart("issueIOU")
}
func TestEvalArray(t *testing.T) {
	fmt.Println("--------put array--------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "linearId"
	operation := "PUT"
	isArray := true
	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)
	hlfStub.MockTransactionStart("issueIOU")
	shim.SetLoggingLevel(shim.LogDebug)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("operation", operation)
	tc.SetInput("isArray", isArray)
	tc.SetInput("containerServiceStub", containerServiceStub)

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	inputmp["Value"] = `[{"lender":"net.i2p.crypto.eddsa.EdDSAPublicKey@c3cefb34","borrower":"net.i2p.crypto.eddsa.EdDSAPublicKey@a7348296","amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_1"}, {"lender":"net.i2p.crypto.eddsa.EdDSAPublicKey@c3cefb34","borrower":"net.i2p.crypto.eddsa.EdDSAPublicKey@a7348296","amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_2"}]`
	v, err := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	err = json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err = act.Eval(tc)
	// We assume there will be no errors
	assert.Nil(t, err)
	key, _ := hlfStub.CreateCompositeKey(assetName, []string{"ae846b5b-2189-46b4-bda7-8511ff33ba6b_1"})
	value, err := hlfStub.GetState(key)
	assert.Nil(t, err)
	fmt.Printf("state=%v\n", string(value))

	key, _ = hlfStub.CreateCompositeKey(assetName, []string{"ae846b5b-2189-46b4-bda7-8511ff33ba6b_2"})
	value, err = hlfStub.GetState(key)
	assert.Nil(t, err)
	fmt.Printf("state=%v\n", string(value))
	hlfStub.MockTransactionStart("issueIOU")
}

func put(containerServiceStub *dthlf.HyperledgerFabricContainerService) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "linearId"
	operation := "PUT"
	isArray := true
	shim.SetLoggingLevel(shim.LogDebug)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("operation", operation)
	tc.SetInput("isArray", isArray)
	tc.SetInput("containerServiceStub", containerServiceStub)

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	inputmp["Value"] = `[{"lender":"net.i2p.crypto.eddsa.EdDSAPublicKey@c3cefb34","borrower":"net.i2p.crypto.eddsa.EdDSAPublicKey@a7348296","amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_1"}, {"lender":"net.i2p.crypto.eddsa.EdDSAPublicKey@c3cefb34","borrower":"net.i2p.crypto.eddsa.EdDSAPublicKey@a7348296","amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_2"}]`
	v, _ := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("err %v", err)
	}
}

func putComposite(containerServiceStub *dthlf.HyperledgerFabricContainerService) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "lender,borrower"
	operation := "PUT"
	isArray := true
	shim.SetLoggingLevel(shim.LogDebug)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("operation", operation)
	tc.SetInput("isArray", isArray)
	tc.SetInput("containerServiceStub", containerServiceStub)
	tc.SetInput(ivCompositeKeys, "borrower,lender")

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	inputmp["Value"] = `[{"lender":"Alice","borrower":2,"amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_1"}, {"lender":"Charlie","borrower":"Frank","amt":{"quantity":10000,"currency":"USD"},"paid":{"quantity":0,"currency":"USD"},"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_2"}]`
	v, _ := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("err %v\n", err)
	}

	result := tc.GetOutput(ovOutput).(*data.ComplexObject).Value.([]interface{})

	for _, m := range result {
		fmt.Printf("Putcomposite output=%v\n", m)
	}

}

func TestEvalGet(t *testing.T) {
	fmt.Println("--------GET--------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "linearId"
	operation := "GET"
	isArray := true
	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)
	hlfStub.MockTransactionStart("issueIOU")
	shim.SetLoggingLevel(shim.LogDebug)

	put(containerServiceStub)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("operation", operation)
	tc.SetInput("isArray", isArray)
	tc.SetInput("containerServiceStub", containerServiceStub)

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	inputmp["Value"] = `[{"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_1"}, {"linearId":"ae846b5b-2189-46b4-bda7-8511ff33ba6b_2"}]`
	v, _ := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("err %v", err)
	}
	// We assume there will be no errors

	complexvalue := tc.GetOutput(ovOutput).(*data.ComplexObject)
	values := complexvalue.Value.([]interface{})
	for _, v := range values {
		if v != nil {
			fmt.Printf("state=%v\n", v)
		}
	}
	hlfStub.MockTransactionStart("issueIOU")
}

func TestEvalGetComposite(t *testing.T) {
	fmt.Println("--------GETComposite--------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "lender,borrower"
	operation := "GET"
	//	isArray := true
	isArray := false
	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)
	hlfStub.MockTransactionStart("issueIOU")
	shim.SetLoggingLevel(shim.LogDebug)

	putComposite(containerServiceStub)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("operation", operation)
	tc.SetInput("isArray", isArray)
	tc.SetInput("containerServiceStub", containerServiceStub)

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	//inputmp["Value"] = `[{"lender":"Alice", "borrower":"Bob"}, {"lender":"Charlie", "borrower":"Frank"}]`
	inputmp["Value"] = `{"lender":"Alice2notthere", "borrower":2}`

	v, _ := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("err %v\n", err)
	}

	// We assume there will be no errors

	complexvalue := tc.GetOutput(ovOutput).(*data.ComplexObject)
	switch values := complexvalue.Value.(type) {
	case []interface{}:
		for _, v := range values {
			if v != nil {
				fmt.Printf("state=%v\n", v)
			}
		}
		break
	case interface{}:
		fmt.Printf("state=%v\n", v)
	}
	hlfStub.MockTransactionStart("issueIOU")
}

func TestEvalLookupComposite(t *testing.T) {
	fmt.Println("--------GETCompositeLOOKUP--------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "lender,borrower"
	operation := "LOOKUP"
	//	isArray := true
	isArray := true
	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)
	hlfStub.MockTransactionStart("issueIOU")
	shim.SetLoggingLevel(shim.LogDebug)

	putComposite(containerServiceStub)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("operation", operation)
	tc.SetInput("isArray", isArray)
	tc.SetInput(ivCompositeKeys, "borrower,lender")
	tc.SetInput(ivAssetLookupKey, "lender,borrower")
	tc.SetInput("containerServiceStub", containerServiceStub)

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	//inputmp["Value"] = `{"lender":"Alice"}`
	inputmp["Value"] = `[{"lender":"Charlie"}, {"lender":"Alice"}]`
	v, _ := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("err %v\n", err)
	}
	// We assume there will be no errors
	assert.Nil(t, err)
	complexvalue := tc.GetOutput(ovOutput).(*data.ComplexObject)
	switch values := complexvalue.Value.(type) {
	case []interface{}:
		for _, v := range values {
			if v != nil {
				fmt.Printf("state=%v\n", v)
			}
		}
		break
	case interface{}:
		fmt.Printf("state=%v\n", v)
	}

	hlfStub.MockTransactionStart("issueIOU")
}

func TestEvalDeleteComposite(t *testing.T) {
	fmt.Println("--------GETCompositeDelete--------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	assetName := "com.tibco.cp.iou"
	assetKey := "lender,borrower"
	operation := "DELETE"
	//	isArray := true
	isArray := false
	hlfStub := shim.NewMockStub("IOU", &TestChainCode{})
	containerServiceStub := dthlf.NewHyperledgerFabricContainerService(hlfStub)
	hlfStub.MockTransactionStart("issueIOU")
	shim.SetLoggingLevel(shim.LogDebug)

	putComposite(containerServiceStub)

	tc.SetInput("assetName", assetName)
	tc.SetInput("identifier", assetKey)
	tc.SetInput("operation", operation)
	tc.SetInput("isArray", isArray)
	tc.SetInput(ivCompositeKeys, "borrower,lender")
	tc.SetInput("containerServiceStub", containerServiceStub)

	//fmt.Printf("stub = %#v\n", tc.GetInput("containerServiceStub"))

	inputmp := make(map[string]interface{})
	//inputmp["Value"] = `{"lender":"Alice"}`
	inputmp["Value"] = `{"lender":"Charlie", "borrower":"Frank"}`
	v, _ := json.Marshal(inputmp)
	complexObject := &data.ComplexObject{}
	json.Unmarshal([]byte(v), complexObject)

	tc.SetInput("input", complexObject)
	//fmt.Printf("data = %#v\n", tc.GetInput("data"))

	// Execute the activity
	_, err := act.Eval(tc)
	if err != nil {
		fmt.Printf("err %v\n", err)
	}
	// We assume there will be no errors
	assert.Nil(t, err)
	complexvalue := tc.GetOutput(ovOutput).(*data.ComplexObject)
	switch values := complexvalue.Value.(type) {
	case []interface{}:
		for _, v := range values {
			if v != nil {
				fmt.Printf("state1=%v\n", v)
			}
		}
		break
	case interface{}:
		fmt.Printf("state2=%v\n", values)
	}

	hlfStub.MockTransactionStart("issueIOU")
}
