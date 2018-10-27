package collection

// Imports
import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"

	impl "github.com/TIBCOSoftware/dovetail-contrib/smartcontract-go/activity/collection"
)

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return impl.NewActivity(metadata)
}
