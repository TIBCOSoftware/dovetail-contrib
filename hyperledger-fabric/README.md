# Flogo extensions for Hyperledger Fabric

This package includes [Flogo®](https://www.flogo.io/) extensions for developing blockchain apps and smart contracts on [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric). It also provides sample blockchain apps using the Flogo visual programming environment with zero-code. These Dovetail extensions can be used to model Hyperledger Fabric chaincode and client apps in [TIBCO Flogo® Enterprise v2.10.0](https://docs.tibco.com/products/tibco-flogo-enterprise-2-10-0).

The following Flogo® extensions are currently available:

- [Fabric Chaincode Extension](fabric), which includes connector, trigger and activities for designing and implementing Hyperledger Fabric chaincode.
- [Fabric Client Extension](fabclient), which includes connector, trigger, and activities for designing and implementing Hyperledger Fabric client apps, such as a REST or GraphQL service that interacts with a Hyperledger Fabric network.

## Getting Started

To setup the local environment on Mac or Linux for Hyperledger Fabric development, follow the instructions [here](development.md).

Then, start by looking at the following end-to-end samples:

- [`marble`](samples/marble), which is a zero-code version of the `marbles02` chaincode in [`fabric samples`](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode) and a REST service for client to submit chaincode transactions. It demonstrates basic features of the Hyperledger Fabric, including the creeation and update of states and composite-keys, and various types of queries for state and history with pagination.
- [`marble-private`](samples/marble-private), which is a zero-code version of the `marbles02_private` chaincode in [`fabric samples`](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode) and a REST service. It demonstrates the use of private collections.
- [`equipment`](samples/equipment), which implements the chaincode and REST and GraphQL services for tracking equipment purchasing and installation coordinated by multiple clients. It demonstrates the use of Hyperledger Fabric events and event listners.
- [`audit`](samples/audit) is an audit-trace app used by the TIBCO AuditSafe cloud service. It supports multi-tenant multi-domain audit log and reporting requirements.
- [`iou`](samples/iou) is an advanced sample that implements a cross-border payment network similar to a simplified [Ripple network](https://www.ripple.com/files/ripple_product_overview.pdf). It implements both a required chaincode and a client app with [GraphQL](https://graphql.org/) service interface. The chaincode uses some more advanced Hyperledger Fabric features, including ABAC and private collections. This sample illustrates how a real-worlld Hyperleddger Fabric app can be implemented with zero-code.

By comparing other implementations of chaincode and client apps, you can see that hundreds of lines of boilerplate code are replaced by a single JSON model file exported from the Dovetail or TIBCO Flogo® Enterprise modeling UI. Besides, by using the Flogo visual programming environment, you do not have to learn much of the blockchain APIs nor special programming language for smart contracts. You can implement chaincode and client apps for Hyperledger Fabric by simple drag-drop-mapping in Flogo.

## Modeling with TIBCO Cloud Integration (TCI)

If you are already a subscriber of [TIBCO Cloud Integration (TCI)](https://cloud.tibco.com/), or you plan to sign-up for a TCI trial, you can easily start the development of Hyperledger Fabric apps by using a Chrome browser. Refer to [Modeling with TCI](tci) for more detailed instructions.

## Deploy to Kubernetes in Cloud

Dovetail apps and chaincodes can be deployed to Kubernetes in any of the supported cloud services, including AWS, Azure, and GCP. Refer to [operation](./operation) for detailed instructions for creating Kubernetes clusters and managing Hyperledger Fabric network and chaincodes.
