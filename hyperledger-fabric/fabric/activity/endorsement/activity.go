package endorsement

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/statebased"
	"github.com/pkg/errors"
	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
)

const (
	ivKey           = "key"
	ivOperation     = "operation"
	ivRole          = "role"
	ivOrganizations = "organizations"
	ivPolicy        = "policy"
	ivIsPrivate     = "isPrivate"
	ivCollection    = "collection"
	ovCode          = "code"
	ovMessage       = "message"
	ovKey           = "key"
	ovResult        = "result"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-endorsement")

func init() {
	common.SetChaincodeLogLevel(log)
}

// FabricEndorsementActivity is a stub for executing Hyperledger Fabric get operations
type FabricEndorsementActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new FabricEndorsementActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FabricEndorsementActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *FabricEndorsementActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *FabricEndorsementActivity) Eval(ctx activity.Context) (done bool, err error) {
	// check input args
	key, ok := ctx.GetInput(ivKey).(string)
	if !ok || key == "" {
		log.Error("state key is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "state key is not specified")
		return false, errors.New("state key is not specified")
	}
	log.Debugf("state key: %s\n", key)
	ops, ok := ctx.GetInput(ivOperation).(string)
	if !ok || ops == "" {
		log.Error("operation is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "operation is not specified")
		return false, errors.New("operation is not specified")
	}
	log.Debugf("operation: %s\n", ops)

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, err.Error())
		return false, err
	}

	if isPrivate, ok := ctx.GetInput(ivIsPrivate).(bool); ok && isPrivate {
		// set endorsement policy for a key on a private collection
		return setPrivatePolicy(ctx, stub, key, ops)
	}

	// set endorsement policy for the key
	return setPolicy(ctx, stub, key, ops)
}

func setPrivatePolicy(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key, operation string) (bool, error) {
	// set endorsement policy on a private collection
	collection, ok := ctx.GetInput(ivCollection).(string)
	if !ok || collection == "" {
		log.Error("private collection is not specified\n")
		ctx.SetOutput(ovCode, 400)
		ctx.SetOutput(ovMessage, "private collection is not specified")
		return false, errors.New("private collection is not specified")
	}
	ep, err := ccshim.GetPrivateDataValidationParameter(collection, key)
	if err != nil {
		log.Errorf("failed to retrieve policy for private collection %s: %+v\n", collection, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve policy for private collection %s: %+v", collection, err))
		return false, errors.Wrapf(err, "failed to retrieve policy for private collection %s", collection)
	}

	stateEP, err := getUpdatedPolicy(ctx, operation, ep)
	if err != nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to create policy: %+v", err))
		return false, errors.Wrapf(err, "failed to create policy")
	}

	if operation != "LIST" {
		epBytes, err := stateEP.Policy()
		if err != nil {
			log.Errorf("failed to marshal policy: %+v\n", err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to marshal policy: %+v", err))
			return false, errors.Wrapf(err, "failed to marshal policy")
		}

		// update endorsement policy for key
		if err := ccshim.SetPrivateDataValidationParameter(collection, key, epBytes); err != nil {
			log.Errorf("failed to set policy on private collecton %s: %+v\n", collection, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to set policy on private collecton %s: %+v", collection, err))
			return false, errors.Wrapf(err, "failed to to set policy on private collecton %s", collection)
		}
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("updated policy for key %s on private collection %s", key, collection))
	ctx.SetOutput(ovKey, key)
	orgs := stateEP.ListOrgs()
	if len(orgs) > 0 {
		ctx.SetOutput(ovResult, strings.Join(orgs, ","))
	}
	return true, nil
}

func getUpdatedPolicy(ctx activity.Context, operation string, ep []byte) (statebased.KeyEndorsementPolicy, error) {
	switch operation {
	case "ADD":
		return addOrgsToPolicy(ctx, ep)
	case "DELETE":
		return deleteOrgsFromPolicy(ctx, ep)
	case "LIST":
		return statebased.NewStateEP(ep)
	case "SET":
		return createNewPolicy(ctx)
	default:
		log.Errorf("operation %s is not supported", operation)
		return nil, errors.Errorf("operation %s is not supported", operation)
	}
}

func createNewPolicy(ctx activity.Context) (statebased.KeyEndorsementPolicy, error) {
	// create new policy from policy string
	newPolicy, ok := ctx.GetInput(ivPolicy).(string)
	if !ok {
		log.Errorf("policy is not specified for SET operation\n")
		return nil, errors.New("policy is not specified for SET operation")
	}
	envelope, err := cauthdsl.FromString(newPolicy)
	if err != nil {
		log.Errorf("failed to parse policy string %s: %+v\n", newPolicy, err)
		return nil, errors.Wrapf(err, "failed to parse policy string %s", newPolicy)
	}
	epBytes, err := proto.Marshal(envelope)
	if err != nil {
		log.Errorf("failed to marshal signature policy: %+v\n", err)
		return nil, errors.Wrapf(err, "failed to marshal signature policy")
	}
	return statebased.NewStateEP(epBytes)
}

func deleteOrgsFromPolicy(ctx activity.Context, ep []byte) (statebased.KeyEndorsementPolicy, error) {
	stateEP, err := statebased.NewStateEP(ep)
	if err != nil {
		log.Errorf("failed to construct policy from channel default: %+v\n", err)
		return nil, err
	}
	orgs, err := getOrganizations(ctx)
	if err != nil {
		return nil, err
	}
	stateEP.DelOrgs(orgs...)
	return stateEP, nil
}

func addOrgsToPolicy(ctx activity.Context, ep []byte) (statebased.KeyEndorsementPolicy, error) {
	stateEP, err := statebased.NewStateEP(ep)
	if err != nil {
		log.Errorf("failed to construct policy from channel default: %+v\n", err)
		return nil, err
	}
	orgs, err := getOrganizations(ctx)
	if err != nil {
		return nil, err
	}
	role, ok := ctx.GetInput(ivRole).(string)
	if !ok {
		log.Errorf("role is not specified for Add operation\n")
		return nil, errors.New("role is not specified for Add operation")
	}
	err = stateEP.AddOrgs(statebased.RoleType(role), orgs...)
	return stateEP, err
}

func getOrganizations(ctx activity.Context) ([]string, error) {
	orgs, ok := ctx.GetInput(ivOrganizations).(string)
	if !ok {
		log.Errorf("organization is not specified for Add operation\n")
		return nil, errors.New("organization is not specified for Add operation")
	}
	orgArray := strings.Split(orgs, ",")
	for i := range orgArray {
		orgArray[i] = strings.TrimSpace(orgArray[i])
	}
	return orgArray, nil
}

func setPolicy(ctx activity.Context, ccshim shim.ChaincodeStubInterface, key, operation string) (bool, error) {
	// set endorsement policy for a key
	ep, err := ccshim.GetStateValidationParameter(key)
	if err != nil {
		log.Errorf("failed to retrieve policy for key %s: %+v\n", key, err)
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to retrieve policy for key %s: %+v", key, err))
		return false, errors.Wrapf(err, "failed to retrieve policy for key %s", key)
	}

	stateEP, err := getUpdatedPolicy(ctx, operation, ep)
	if err != nil {
		ctx.SetOutput(ovCode, 500)
		ctx.SetOutput(ovMessage, fmt.Sprintf("failed to create policy: %+v", err))
		return false, errors.Wrapf(err, "failed to create policy")
	}

	if operation != "LIST" {
		epBytes, err := stateEP.Policy()
		if err != nil {
			log.Errorf("failed to marshal policy: %+v\n", err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to marshal policy: %+v", err))
			return false, errors.Wrapf(err, "failed to marshal policy")
		}

		// update endorsement policy for key
		if err := ccshim.SetStateValidationParameter(key, epBytes); err != nil {
			log.Errorf("failed to set policy for key %s: %+v\n", key, err)
			ctx.SetOutput(ovCode, 500)
			ctx.SetOutput(ovMessage, fmt.Sprintf("failed to set policy for key %s: %+v", key, err))
			return false, errors.Wrapf(err, "failed to to set policy for key %s", key)
		}
	}

	ctx.SetOutput(ovCode, 200)
	ctx.SetOutput(ovMessage, fmt.Sprintf("updated policy for key %s", key))
	ctx.SetOutput(ovKey, key)
	orgs := stateEP.ListOrgs()
	if len(orgs) > 0 {
		ctx.SetOutput(ovResult, strings.Join(orgs, ","))
	}
	return true, nil
}
