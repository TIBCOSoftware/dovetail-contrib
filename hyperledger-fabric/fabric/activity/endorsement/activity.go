package endorsement

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/statebased"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-endorsement")

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	common.SetChaincodeLogLevel(log)
	_ = activity.Register(&Activity{}, New)
}

// Activity is a stub for executing Hyperledger Fabric get operations
type Activity struct {
}

// New creates a new Activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	return &Activity{}, nil
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	input := &Input{}
	if err = ctx.GetInputObject(input); err != nil {
		return false, err
	}

	if input.StateKey == "" {
		log.Error("state key is not specified\n")
		output := &Output{Code: 400, Message: "state key is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("state key: %s\n", input.StateKey)
	if input.Operation == "" {
		log.Error("operation is not specified\n")
		output := &Output{Code: 400, Message: "operation is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	log.Debugf("operation: %s\n", input.Operation)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	if input.IsPrivate {
		// set endorsement policy for a key on a private collection
		return setPrivatePolicy(ctx, stub, input)
	}

	// set endorsement policy for the key
	return setPolicy(ctx, stub, input)
}

func setPrivatePolicy(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	// set endorsement policy on a private collection
	if input.PrivateCollection == "" {
		log.Error("private collection is not specified\n")
		output := &Output{Code: 400, Message: "private collection is not specified"}
		ctx.SetOutputObject(output)
		return false, errors.New(output.Message)
	}
	ep, err := ccshim.GetPrivateDataValidationParameter(input.PrivateCollection, input.StateKey)
	if err != nil {
		log.Errorf("failed to retrieve policy for private collection %s: %+v\n", input.PrivateCollection, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve policy for private collection %s", input.PrivateCollection)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	stateEP, err := getUpdatedPolicy(ctx, ep, input)
	if err != nil {
		output := &Output{Code: 500, Message: "failed to create policy"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if input.Operation != "LIST" {
		epBytes, err := stateEP.Policy()
		if err != nil {
			log.Errorf("failed to marshal policy: %+v\n", err)
			output := &Output{Code: 500, Message: "failed to marshal policy"}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}

		// update endorsement policy for key
		if err := ccshim.SetPrivateDataValidationParameter(input.PrivateCollection, input.StateKey, epBytes); err != nil {
			log.Errorf("failed to set policy on private collecton %s: %+v\n", input.PrivateCollection, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to set policy on private collecton %s", input.PrivateCollection)}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}
	}

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("updated policy for key %s on private collection %s", input.StateKey, input.PrivateCollection),
		StateKey: input.StateKey,
		Result:   "",
	}
	orgs := stateEP.ListOrgs()
	if len(orgs) > 0 {
		output.Result = strings.Join(orgs, ",")
	}
	ctx.SetOutputObject(output)
	return true, nil
}

func getUpdatedPolicy(ctx activity.Context, ep []byte, input *Input) (statebased.KeyEndorsementPolicy, error) {
	switch input.Operation {
	case "ADD":
		return addOrgsToPolicy(ctx, ep, input)
	case "DELETE":
		return deleteOrgsFromPolicy(ctx, ep, input)
	case "LIST":
		return statebased.NewStateEP(ep)
	case "SET":
		return createNewPolicy(input.Policy)
	default:
		log.Errorf("operation %s is not supported", input.Operation)
		return nil, errors.Errorf("operation %s is not supported", input.Operation)
	}
}

func createNewPolicy(policy string) (statebased.KeyEndorsementPolicy, error) {
	// create new policy from policy string
	if policy == "" {
		log.Errorf("policy is not specified for SET operation\n")
		return nil, errors.New("policy is not specified for SET operation")
	}
	envelope, err := cauthdsl.FromString(policy)
	if err != nil {
		log.Errorf("failed to parse policy string %s: %+v\n", policy, err)
		return nil, errors.Wrapf(err, "failed to parse policy string %s", policy)
	}
	epBytes, err := proto.Marshal(envelope)
	if err != nil {
		log.Errorf("failed to marshal signature policy: %+v\n", err)
		return nil, errors.Wrapf(err, "failed to marshal signature policy")
	}
	return statebased.NewStateEP(epBytes)
}

func deleteOrgsFromPolicy(ctx activity.Context, ep []byte, input *Input) (statebased.KeyEndorsementPolicy, error) {
	stateEP, err := statebased.NewStateEP(ep)
	if err != nil {
		log.Errorf("failed to construct policy from channel default: %+v\n", err)
		return nil, err
	}
	orgs, err := getOrganizations(input.Organizations)
	if err != nil {
		return nil, err
	}
	stateEP.DelOrgs(orgs...)
	return stateEP, nil
}

func addOrgsToPolicy(ctx activity.Context, ep []byte, input *Input) (statebased.KeyEndorsementPolicy, error) {
	stateEP, err := statebased.NewStateEP(ep)
	if err != nil {
		log.Errorf("failed to construct policy from channel default: %+v\n", err)
		return nil, err
	}
	orgs, err := getOrganizations(input.Organizations)
	if err != nil {
		return nil, err
	}
	if input.Role == "" {
		log.Errorf("role is not specified for Add operation\n")
		return nil, errors.New("role is not specified for Add operation")
	}
	err = stateEP.AddOrgs(statebased.RoleType(input.Role), orgs...)
	return stateEP, err
}

func getOrganizations(orgs string) ([]string, error) {
	if orgs == "" {
		log.Errorf("organization is not specified\n")
		return nil, errors.New("organization is not specified")
	}
	orgArray := strings.Split(orgs, ",")
	for i := range orgArray {
		orgArray[i] = strings.TrimSpace(orgArray[i])
	}
	return orgArray, nil
}

func setPolicy(ctx activity.Context, ccshim shim.ChaincodeStubInterface, input *Input) (bool, error) {
	// set endorsement policy for a key
	ep, err := ccshim.GetStateValidationParameter(input.StateKey)
	if err != nil {
		log.Errorf("failed to retrieve policy for key %s: %+v\n", input.StateKey, err)
		output := &Output{Code: 500, Message: fmt.Sprintf("failed to retrieve policy for key %s", input.StateKey)}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	stateEP, err := getUpdatedPolicy(ctx, ep, input)
	if err != nil {
		output := &Output{Code: 500, Message: "failed to create policy"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	if input.Operation != "LIST" {
		epBytes, err := stateEP.Policy()
		if err != nil {
			log.Errorf("failed to marshal policy: %+v\n", err)
			output := &Output{Code: 500, Message: "failed to marshal policy"}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}

		// update endorsement policy for key
		if err := ccshim.SetStateValidationParameter(input.StateKey, epBytes); err != nil {
			log.Errorf("failed to set policy for key %s: %+v\n", input.StateKey, err)
			output := &Output{Code: 500, Message: fmt.Sprintf("failed to to set policy for key %s", input.StateKey)}
			ctx.SetOutputObject(output)
			return false, errors.Wrapf(err, output.Message)
		}
	}

	output := &Output{
		Code:     200,
		Message:  fmt.Sprintf("updated policy for key %s", input.StateKey),
		StateKey: input.StateKey,
		Result:   "",
	}
	orgs := stateEP.ListOrgs()
	if len(orgs) > 0 {
		output.Result = strings.Join(orgs, ",")
	}
	ctx.SetOutputObject(output)
	return true, nil
}
