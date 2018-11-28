/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package collection

// Imports
import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"

	impl "github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/activity/collection"
)

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return impl.NewActivity(metadata)
}
