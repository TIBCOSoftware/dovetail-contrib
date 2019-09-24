# Flogo extension for Hyperledger Fabric client

This Flogo extension is designed to allow developers to use the Flogo visual programming environment to design and implement apps or services that interact with a Hyperledger Fabric network.  This extension supports the following release versions:
- [TIBCO Flogo® Enterprise 2.6](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1)
- [Hyperledger Fabric 1.4](https://www.hyperledger.org/projects/fabric)
- [Hyperledger Fabric Go SDK v1.0.0-alpha5](https://github.com/hyperledger/fabric-sdk-go)

The [Fabric Connector](connector/fabconnector) allows you to configure connections of the target Hyperledger Fabric network.

The [Fabric Event Listener Trigger](trigger/eventlistener) allows you to implement blockchain apps that listens to Hyperledger Fabric events, including block events, filtered block events, and chaincode events.

The [Fabric Request Activity](activity/fabrequest) allows you to implement client apps that interacts with Hyperledger Fabric network by submitting `query` or `invoke` requests.  A Fabric `invoke` request can execute a specified chaincode to create or update states in distributed ledger or private collections.  A Fabric `query` request can execute a specified chaincode to query current states or history in distributed ledger or private collections without changing any state.

With these extensions, Hyperledger Fabric client apps can be designed and implemented with zero code. Refer to the sample [`equipment`](../samples/equipment) for more details about implementing REST or GraphQL service that interacts with a Hyperledger Fabric network.

To use this extension in Flogo model, you can create [`fabclientExtension.zip`](../fabclientExtension.zip) by using the script [`zip-fabclient.sh`](../zip-fabclient.sh), and then upload the zip-file to the `TIBCO Flogo® Enterprise 2.6` or `Dovetail v0.2.0` as an extension, and so they are available for modeling client apps and services for Hyperledger Fabric.
