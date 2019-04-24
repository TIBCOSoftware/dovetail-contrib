# fabric-illustration
This app illustrates all [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-4-0) triggers and activities that can be used to implement chaincode for [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric).  It includes all chaincode APIs currently supported by Hyperledger Fabric release 1.4.

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.4](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html)
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric Samples](https://github.com/hyperledger/fabric-samples)
- Clone [This Flogo extension](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric)

## Upload and view the illustration
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.4.0/doc/pdf/TIB_flogo_2.4_users_guide.pdf?id=1)
- Upload [`fabticExtension.zip`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can recreate this `zip` by using the script [`zip-fabric.sh`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/zip-fabric.sh)
- Create new Flogo App of name `all_fabric` and choose `Import app` to import the model [`all_fabric.json`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/fabric-illustration/all_fabric.json)
- You can then view or edit the trigger and activities in UI modeler of the TIBCO Flogo® Enterprise.

## Build and deploy chaincode to Hyperledger Fabric

This model is for illustration only.  For building and deploying real chaincode, refer the sample [marble-app](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/marble-app).
