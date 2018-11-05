package publisher

// Imports
import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"

	"github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/utils"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// Constants
const (
	ivEvent    = "event"
	ivData     = "input"
	ivStub     = "containerServiceStub"
	ivMetadata = "eventMetadata"
)

// describes the metadata of the activity as found in the activity.json file
type EventPublisherActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &EventPublisherActivity{metadata: metadata}
}

func (a *EventPublisherActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval executes the activity
func (a *EventPublisherActivity) Eval(context activity.Context) (done bool, err error) {
	stub, err := utils.GetContainerStub(context)
	if err != nil {
		return false, err
	}

	event, ok := context.GetInput(ivEvent).(string)
	if !ok {
		return false, fmt.Errorf("event name is not initialized")
	}

	evtMetadata, ok := context.GetInput(ivMetadata).(string)
	if !ok {
		return false, fmt.Errorf("operation is not initialized")
	}
	evtValue, err := data.CoerceToComplexObject(context.GetInput(ivData))
	if err != nil {
		return false, fmt.Errorf("event value is not initialized")
	}

	payload, err := json.Marshal(evtValue.Value)
	if err != nil {
		return false, err
	}

	err = stub.GetEventService().Publish(event, evtMetadata, payload)
	if err != nil {
		return false, err
	}

	return true, nil
}
