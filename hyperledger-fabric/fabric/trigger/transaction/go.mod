module github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/trigger/transaction

go 1.14

replace github.com/project-flogo/flow => github.com/yxuco/flow v1.1.1

replace github.com/project-flogo/core => github.com/yxuco/core v1.1.1

require (
	github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric/common v1.0.0
	github.com/hyperledger/fabric v1.4.9
	github.com/project-flogo/core v1.1.0
	github.com/stretchr/testify v1.6.1
	github.com/xeipuuv/gojsonschema v1.2.0
)
