/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package analytics

// Imports
import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// Constants
const (
	ivLogLevel = "logLevel"
	ivMessage  = "message"
	ivErrcode  = "errorCode"
	ivStub     = "containerServiceStub"
	ivContract = "FlowCC"
)

// describes the metadata of the activity as found in the activity.json file
type AnalyticsActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AnalyticsActivity{metadata: metadata}
}

func (a *AnalyticsActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *AnalyticsActivity) Eval(context activity.Context) (done bool, err error) {
	// Add implementation here
	return true, nil
}
