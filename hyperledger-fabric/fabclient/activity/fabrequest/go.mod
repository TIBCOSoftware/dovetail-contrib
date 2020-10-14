module github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabclient/activity/fabrequest

go 1.14

replace github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabclient/common => /Users/yxu/work/dovetail/dovetail-contrib/hyperledger-fabric/fabclient/common

replace github.com/project-flogo/flow => github.com/yxuco/flow v1.1.1

replace github.com/project-flogo/core => github.com/yxuco/core v1.1.1

require (
	github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabclient/common v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/project-flogo/core v1.1.0
	github.com/stretchr/testify v1.6.1
)
