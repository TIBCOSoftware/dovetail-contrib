module github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/activity/getbycompositekey

go 1.14

replace github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common => /Users/yxu/work/dovetail/dovetail-contrib/hyperledger-fabric/fabric/common

replace github.com/project-flogo/flow => /Users/yxu/go/src/github.com/project-flogo/flow

replace github.com/project-flogo/core => github.com/yxuco/core v1.1.1-0.20201008191127-cf81f70151d2

require (
	github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common v0.0.0-00010101000000-000000000000
	github.com/hyperledger/fabric v1.4.9
	github.com/pkg/errors v0.9.1
	github.com/project-flogo/core v1.1.0
	github.com/stretchr/testify v1.6.1
)
