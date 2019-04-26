# flogo-enterprise-app
This package includes [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-4-0) extensions for developing blockchain apps and smart contracts on [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric).

It provides the following Flogo® extensions:
- [Fabric Chaincode Extension](fabric), which includes triggers and activities for designing and implementing Hyperledger Fabric chaincode with zero-code.
- [Fabric Client Extension](fabclient), which includes connector and activities for designing and implementing Hyperledger Fabric client apps, such as a REST service that interacts with a Hyperledger Fabric network, by visual programming with zero-code.

Samples for Hyperledger Fabric chaincode:
- [`marble-app`](marble-app)
- [`marble-private`](marble-private)

These are the zero-code version of the `marbles02` and `marbles02_private` chaincode in [`fabric samples`](https://github.com/hyperledger/fabric-samples/tree/release-1.4/chaincode).  The hundreds of lines of boilerplate code are replaced by a JSON model file exported from the TIBCO Flogo® Enterprise, where the chaincode is modeled by drag-and-drop.

Sample REST services as a Hyperledger Fabric client:
- [`marble-client`](marble-client)
- [`marble-private-client`](marble-private-client)

These sample REST services are clients of the `marble-app` and `marble-private` networks.  They implement REST APIs that interact with the distributed ledger in Hyperledger Fabric network.  These Flogo apps do not require any code, each app contains only a JSON model file exported from the the TIBCO Flogo® Enterprise, where the REST APIs are modeled by drag-and-drop.
