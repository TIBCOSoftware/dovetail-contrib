/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package txnreply

// Imports
import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// Constants
const (
	ivMessage = "message"
)

// describes the metadata of the activity as found in the activity.json file
type TxnResponseActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &TxnResponseActivity{metadata: metadata}
}

func (a *TxnResponseActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *TxnResponseActivity) Eval(context activity.Context) (done bool, err error) {

	err = fmt.Errorf(context.GetInput(ivMessage).(string))

	return false, err
}
