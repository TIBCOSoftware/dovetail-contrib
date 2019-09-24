/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package empty

// Imports
import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// describes the metadata of the activity as found in the activity.json file
type NullActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &NullActivity{metadata: metadata}
}

func (a *NullActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *NullActivity) Eval(context activity.Context) (done bool, err error) {
	return true, nil
}
