# Dovetail extensions for Hyperledger Fabric
This package includes [Flogo®](https://www.flogo.io/) extensions for developing blockchain apps and smart contracts on [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric).  It also provides sample blockchain apps using the Flogo visual programming environment with zero-code.  These Dovetail extensions can be used to model Hyperledger Fabric chaincode and client apps in either [TIBCO Flogo® Enterprise v2.6.1](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) or [Dovetail v0.2.0](https://github.com/TIBCOSoftware/dovetail).

The following Flogo® extensions are currently available:
- [Fabric Chaincode Extension](fabric), which includes connector, trigger and activities for designing and implementing Hyperledger Fabric chaincode.
- [Fabric Client Extension](fabclient), which includes connector, trigger, and activities for designing and implementing Hyperledger Fabric client apps, such as a REST or GraphQL service that interacts with a Hyperledger Fabric network.

## Getting Started
Start by looking at the following end-to-end samples:
- [`marble`](samples/marble), which is a zero-code version of the `marbles02` chaincode in [`fabric samples`](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode) and a REST service for client to submit chaincode transactions.  It demonstrates basic features of the Hyperledger Fabric, including the creeation and update of states and composite-keys, and various types of queries for state and history with pagination.
- [`marble-private`](samples/marble-private), which is a zero-code version of the `marbles02_private` chaincode in [`fabric samples`](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode) and a REST service.  It demonstrates the use of private collections.
- [ `abac`](samples/abac), which is a sample to demonstrate the use of `Attribute Based Access Control`(ABAC). This feature is important for use-cases, such as IOU trading or exchange, that requires user-level access control.
- [ `equipment`](samples/equipment), which implements the chaincode and REST and GraphQL services for tracking equipment purchasing and installation coordinated by multiple clients. It demonstrates the use of Hyperledger Fabric events.

By comparing other implementations of chaincode and client apps, you can see that hundreds of lines of boilerplate code are replaced by a single JSON model file exported from the TIBCO Flogo® Enterprise or Dovetail.  Besides, by using the Flogo visual programming environment, you do not have to learn much of the blockchain APIs nor special programming language for the smart contract.  You can implement chaincode and client apps for Hyperledger Fabric by simple drag-drop-mapping in Flogo.

![Output sample](import-app.gif)

<video src="import-app.mp4" width="320" height="200" import app></video>