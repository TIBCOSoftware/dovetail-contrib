package history

// Imports
import (
	"fmt"

	"github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/utils"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Constants
const (
	ivAssetName = "assetName"
	ivData      = "input"
	ivAssetKey  = "identifier"
	ivStub      = "containerServiceStub"
	ovOutput    = "output"
)

// describes the metadata of the activity as found in the activity.json file
type HistoryActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &HistoryActivity{metadata: metadata}
}

func (a *HistoryActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval executes the activity
func (a *HistoryActivity) Eval(context activity.Context) (done bool, err error) {
	stub, err := utils.GetContainerStub(context)
	if err != nil {
		return false, err
	}

	logger := stub.GetLogService()
	logger.Debug("Enter activity history ...")

	assetValue, ok := context.GetInput(ivData).(*data.ComplexObject)
	if !ok {
		return false, fmt.Errorf("input is not initialized")
	}

	assetName, ok := context.GetInput(ivAssetName).(string)
	if !ok {
		return false, fmt.Errorf("asset name is not initialized")
	}

	assetKey, ok := context.GetInput(ivAssetKey).(string)
	if !ok {
		return false, fmt.Errorf("asset key is not initialized")
	}

	input, err := utils.GetInputData(assetValue, false)
	if err != nil {
		return false, err
	}
	fmt.Printf("data = %#v\n", input)
	result, err := stub.GetDataService().GetHistory(assetName, assetKey, input[0].(map[string]interface{}))
	if err != nil {
		return false, err
	}

	output, err := utils.ParseRecords(result)
	if err != nil {
		return false, err
	}
	complexOutput := &data.ComplexObject{}
	complexOutput.Value = output
	context.SetOutput(ovOutput, complexOutput)
	logger.Debug("Exit history activity")
	return true, nil
}
