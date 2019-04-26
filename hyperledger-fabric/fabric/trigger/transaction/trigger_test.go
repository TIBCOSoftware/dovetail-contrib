
package transaction


import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	return nil, nil
}

func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	return nil, nil
}

const testConfig string = `{
                "id": "mytrigger",
                    "settings": {
                    "setting": "somevalue"
                },
                "handlers": [
                    {
                        "actionId": "test_action",
                        "settings": {
                            "handler_setting": "somevalue"
                        }
                    }
                ]
            }`

func TestInit(t *testing.T) {
	md := &MyTriggerFactory{trigger.NewMetadata(jsonMetadata)}
	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	f := md.New(config)

	_, isNew := f.(trigger.Initializable)

	if !isNew {
		runner := &TestRunner{}
		tgr, isOld := f.(trigger.InitOld)
		if isOld {
			tgr.Init(runner)

		}
	}
}
