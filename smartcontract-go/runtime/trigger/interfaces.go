package trigger

import (
	"github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/runtime/services"
	"github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/runtime/transaction"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

type SmartContractTrigger interface {
	trigger.Trigger
	trigger.Initializable
	Invoke(stub services.ContainerService, txn transaction.TransactionService) (status bool, data interface{}, err error)
}
