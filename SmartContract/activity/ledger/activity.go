/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package ledger

import (
	impl "github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/activity/ledger"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return impl.NewActivity(metadata)
}
