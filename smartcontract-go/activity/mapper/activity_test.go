package transform

// Imports
import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/split"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"

	_ "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/split"
)

// activityMetadata is the metadata of the activity as described in activity.json
// We'll store it as a variable to reuse it across multiple testcases
var activityMetadata *activity.Metadata

// getActivityMetadata reads the activity.json file and sets the activityMetadata variable
// if the variable already contains metadata it simply returns the current value rather than reading the file again
func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("../../../SmartContract/activity/mapper/activity.json")
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

// Make sure that you have updated the values below
func TestEval(t *testing.T) {
	fmt.Println("------Testprimitive------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	val := make(map[string]interface{})
	val["field"] = 100
	complexObject := &data.ComplexObject{}
	complexObject.Value = val
	tc.SetInput(ivInput, complexObject)
	tc.SetInput(ivDatatype, "Long")
	tc.SetInput(ivIsArray, false)

	// Execute the activity
	_, err := act.Eval(tc)

	// We assume there will be no errors
	assert.Nil(t, err)

	output := tc.GetOutput("output")
	result, _ := data.CoerceToComplexObject(output)
	assert.EqualValues(t, 100, result.Value.(map[string]interface{})["field"].(int))
}

func TestUerDefinedEval(t *testing.T) {
	fmt.Println("------TestUserDefined------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	complexObject := &data.ComplexObject{}
	complexObject.Value = `{"name":"abc"}`
	tc.SetInput(ivUserInput, complexObject)
	tc.SetInput(ivDatatype, "User Defined...")
	tc.SetInput(ivIsArray, false)
	// Execute the activity
	_, err := act.Eval(tc)

	// We assume there will be no errors
	assert.Nil(t, err)

	output := tc.GetOutput("output")
	result, _ := data.CoerceToComplexObject(output)
	fmt.Printf("output = %v\n", result)
}

func TestAssetEval(t *testing.T) {
	fmt.Println("------Testasset------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	complexObject := &data.ComplexObject{}
	complexObject.Value = `[{"name":"abc"}]`
	tc.SetInput(ivInput, complexObject)
	tc.SetInput(ivDatatype, "com.tibco.test")
	tc.SetInput(ivIsArray, true)
	// Execute the activity
	_, err := act.Eval(tc)

	// We assume there will be no errors
	assert.Nil(t, err)

	output := tc.GetOutput("output")
	result, _ := data.CoerceToComplexObject(output)
	fmt.Printf("output = %v\n", result)
}

func TestPrimitiveObjArrayEval(t *testing.T) {
	fmt.Println("------TestprimitiveObjArray------")
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	fmt.Println("obj array to obj array")
	input := make([]map[string]interface{}, 0)
	val1 := make(map[string]interface{})
	val1["field"] = 100
	input = append(input, val1)

	val2 := make(map[string]interface{})
	val2["field"] = 200
	input = append(input, val2)

	complexObject := &data.ComplexObject{}
	complexObject.Value = input
	tc.SetInput(ivInput, complexObject)
	tc.SetInput(ivDatatype, "Long")
	tc.SetInput(ivIsArray, true)
	tc.SetInput(ivInputArrayType, "Object Array")
	tc.SetInput(ivOutputArrayType, "Object Array")
	// Execute the activity
	_, err := act.Eval(tc)
	// We assume there will be no errors
	assert.Nil(t, err)
	output := tc.GetOutput("output")
	result, _ := data.CoerceToComplexObject(output)
	fmt.Printf("obj array to obj array, output=%v\n", result.Value.([]map[string]interface{}))

	fmt.Println("obj to primitive array, using previous")
	tc.SetInput(ivOutputArrayType, "Primitive Array")
	// Execute the activity
	_, err = act.Eval(tc)
	// We assume there will be no errors
	assert.Nil(t, err)
	output = tc.GetOutput("output")
	result, _ = data.CoerceToComplexObject(output)
	fmt.Printf("obj array to primitive array, output=%v\n", result.Value.([]interface{}))

	fmt.Println("primitive to obj array using previous output")
	tc.SetInput(ivInputArrayType, "Primitive Array")
	tc.SetInput(ivOutputArrayType, "Object Array")
	tc.SetInput(ivInput, result)
	// Execute the activity
	_, err = act.Eval(tc)
	// We assume there will be no errors
	assert.Nil(t, err)
	output = tc.GetOutput("output")
	result2, _ := data.CoerceToComplexObject(output)
	fmt.Printf("primitive array to obj array, output=%v\n", result2.Value.([]interface{}))

	fmt.Println("primitive to primitive array")
	tc.SetInput(ivInputArrayType, "Primitive Array")
	tc.SetInput(ivOutputArrayType, "Primitive Array")
	tc.SetInput(ivInput, result)
	// Execute the activity
	_, err = act.Eval(tc)

	// We assume there will be no errors
	assert.Nil(t, err)

	output = tc.GetOutput("output")
	result2, _ = data.CoerceToComplexObject(output)
	fmt.Printf("primitive array to primitive array, output=%v\n", result2.Value.([]interface{}))
}

func TestSplitEval(t *testing.T) {
	fmt.Println("------TestSplit------")
	sp := &split.Split{}
	out := sp.Eval("a,b", ",")
	fmt.Printf("split('a,b')=%v\n", out)

	out = sp.Eval("a", ",")
	fmt.Printf("split('a')=%v\n", out)
	/*	act := NewActivity(getActivityMetadata())
		tc := test.NewTestActivityContext(act.Metadata())

		complexObject := &data.ComplexObject{}
		complexObject.Value = `abc`
		tc.SetInput(ivInput, complexObject)
		tc.SetInput(ivDatatype, "String")
		tc.SetInput(ivIsArray, true)
		tc.SetInput(ivInputArrayType, "Primitive Array")
		tc.SetInput(ivOutputArrayType, "Object Array")
		// Execute the activity
		_, err := act.Eval(tc)

		// We assume there will be no errors
		assert.Nil(t, err)

		output := tc.GetOutput("output")
		result, _ := data.CoerceToComplexObject(output)
		fmt.Printf("output = %v\n", result)*/
}
