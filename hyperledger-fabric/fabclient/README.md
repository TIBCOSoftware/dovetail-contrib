# Flogo extension for Hyperledger Fabric client

This Flogo extension is designed to allow delopers to use the zero-code visual programming environment of the TIBCO Flogo® Enterprise to design and implement apps or services that interact with Hyperledger Fabric networks.  This extension supports the following release versions:
- [TIBCO Flogo® Enterprise 2.6](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1)
- [Hyperledger Fabric 1.4](https://www.hyperledger.org/projects/fabric)
- [Hyperledger Fabric Go SDK v1.0.0-alpha5](https://github.com/hyperledger/fabric-sdk-go)

The [Fabric Connector](connector/fabconnector) allows you to configure the target fabric network used by the app.

It supports the following activity for submitting an `invoke` or `query` request to a Hyperledger Fabric network.
- [Fabric Request](activity/fabrequest): It sends a fabric `invoke` request to execute a specified chaincode that inserts or updates records in the distributed ledger or private collections.  It sends a fabric `query` request to execute a specified chaincode that queries the records in the distributed ledger or private collections without changing any state.

More activities can be added to integrate with blockchain events.

With these extensions, Hyperledger Fabric client apps can be designed and implemented with zero code. Refer to the sample [`marble-client`](../marble-client) for more details of a REST service implemented by using `TIBCO Flogo® Enterprise` that updates and retrieves data on a distributed ledger of Hyperledger Fabric.

This extension is packaged as [`fabclientExtension.zip`](../fabclientExtension.zip), which can be re-created from source by using the script [`zip-fabclient.sh`](../zip-fabclient.sh).  You can upload the zip-file to the `TIBCO Flogo® Enterprise` as an extension, and start using the connector and activities.
