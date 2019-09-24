/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package trigger

import (
	"github.com/TIBCOSoftware/dovetail-contrib/libraries/fabric-go/runtime/services"
	"github.com/TIBCOSoftware/dovetail-contrib/libraries/fabric-go/runtime/transaction"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

type SmartContractTrigger interface {
	trigger.Trigger
	trigger.Initializable
	Invoke(stub services.ContainerService, txn transaction.TransactionService) (status bool, data interface{}, err error)
}
