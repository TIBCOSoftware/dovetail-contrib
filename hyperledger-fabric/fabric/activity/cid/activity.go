package cid

import (
	"encoding/json"

	"github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	ci "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"
	"github.com/project-flogo/core/activity"
)

// Create a new logger
var log = shim.NewLogger("activity-fabric-cid")

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

	// get chaincode stub
	stub, err := common.GetChaincodeStub(ctx)
	if err != nil || stub == nil {
		log.Errorf("failed to retrieve fabric stub: %+v\n", err)
		output := &Output{Code: 500, Message: err.Error()}
		ctx.SetOutputObject(output)
		return false, err
	}

	// retrieve data for the key
	return retrieveCid(ctx, stub)
}

func retrieveCid(ctx activity.Context, ccshim shim.ChaincodeStubInterface) (bool, error) {
	// get client identity
	c, err := ci.New(ccshim)
	if err != nil {
		log.Errorf("failed to extract client identity from stub: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to extract client identity from stub"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)

	}

	// retrieve data from client identity
	id, err := c.GetID()
	if err != nil {
		log.Errorf("failed to retrieve client ID: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to retrieve client ID"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	mspid, err := c.GetMSPID()
	if err != nil {
		log.Errorf("failed to retrieve client MSPID: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to retrieve client MSPID"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	log.Debug("client MSPID:", mspid)

	cert, err := c.GetX509Certificate()
	if err != nil {
		log.Errorf("failed to retrieve client certificate: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to retrieve client certificate"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	name := cert.Subject.CommonName
	log.Debug("client subject cn:", name)

	schema, err := common.GetActivityInputSchema(ctx, "attrs")
	if err != nil {
		log.Error("schema not defined for CID attributes\n")
		output := &Output{Code: 400, Message: "schema not defined for CID attributes"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}

	attrs, err := getCidAttributeSpec(schema)
	for k := range attrs {
		v, ok, err := c.GetAttributeValue(k)
		if err != nil {
			log.Errorf("failed to retrieve attribute %s: %+v\n", k, err)
		} else if !ok {
			log.Infof("attribute %s is not found", k)
		} else {
			log.Infof("found attribute %s = %s", k, v)
			attrs[k] = v
		}
	}

	output := &Output{Code: 200,
		Cid:   id,
		Mspid: mspid,
		Name:  name,
		Attrs: attrs,
	}
	msgBytes, err := json.Marshal(output)
	if err != nil {
		log.Errorf("failed to serialize JSON output: %+v\n", err)
		output := &Output{Code: 500, Message: "failed to serialize JSON output"}
		ctx.SetOutputObject(output)
		return false, errors.Wrapf(err, output.Message)
	}
	output.Message = string(msgBytes)
	log.Info("CID content:", output.Message)

	ctx.SetOutputObject(output)
	return true, nil
}

func getCidAttributeSpec(metadata string) (map[string]string, error) {
	// extract object field name and type from JSON schema
	var objectProps struct {
		Props map[string]struct {
			FieldType string `json:"type"`
		} `json:"properties"`
	}
	if err := json.Unmarshal([]byte(metadata), &objectProps); err != nil {
		log.Errorf("failed to extract properties from metadata: %+v", err)
		return nil, err
	}
	if objectProps.Props == nil {
		log.Debug("no attribute specified in metadata %s\n", metadata)
		return nil, nil
	}

	// collect object property name and types
	attrs := make(map[string]string)
	for k, v := range objectProps.Props {
		log.Debugf("CID attribute %s type %s\n", k, v.FieldType)
		attrs[k] = v.FieldType
	}
	return attrs, nil
}
