# Flogo extension for Hyperledger Fabric chaincode

This Flogo extension is designed to allow developers to design and implement Hyperledger Fabric chaincode in the Flogo visual programming environment.  This extension supports the following release versions:
- [TIBCO Flogo® Enterprise 2.8](https://docs.tibco.com/products/tibco-flogo-enterprise-2-8-0)
- [Hyperledger Fabric 1.4.4](https://www.hyperledger.org/projects/fabric)

The [Transaction Trigger](trigger/transaction) allows you to configure the chaincode input and output schema, including normal and/or transient input parameters.

It supports the following activities for storing and querying data on the distributed ledger and/or on private collections.
- [Put](activity/put): Insert or update data on the distributed ledger or a private data collection, and optionally insert its associated compsite keys if they are specified.
- [Put All](activity/putall): Insert a list of records on the distributed ledger or a private collection, and optionally insert composite keys of each record.
- [Get](activity/get): Retrieve a state by a specified key from the distributed ledger or a private collection.
- [Get by Range](activity/getrange): Retrieve all states in a specified range of keys from the distributed ledger or a private collection.  It supports resultset pagination for states from the distributed ledger.
- [Get by Composite Key](activity/getbycompositekey): Retrieve all states by a composite-key filter from the distributed ledger or a private collection.  It supports resultset pagination for states from the distributed ledger.
- [Get History](activity/gethistory): Retrieve the history of a specified key for data on the distributed ledger.
- [Query](activity/query): Retrieve all states by a Couchdb query statement from the distributed ledger or a private collection.  It supports resultset pagination for states from the distributed ledger.
- [Delete](activity/delete): Mark the state as deleted for a specified key from the distributed ledger or a private collection, and deletes its composite keys.  Optionally, it can delete only the state, or only a composite key.
- [Set Event](activity/setevent): Set a specified event and payload for a blockchain transaction.
- [Set Endorsement Policy](activity/endorsement): Set state-based endorsement policy by adding or deleting an endorsement organization, or by specifying a new endorsement policy.
- [Invoke Chaincode](activity/invokechaincode): Invoke a local chaincode, and returns response data from the called transaction.
- [Cid](activity/cid): It inspects the client identification and returns the client's name, MSP, and attributes that can be used for ABAC(Attribute Based Access Control).

With these extensions, Hyperledger Fabric chaincode can be designed and implemented with zero code. Refer to samples [`marble-app`](../samples/marble-app) and [`marble-private`](../samples/marble-private) for more details about implementing chaincode for Hyperledger Fabric.

To use this extension in Flogo model, you can create [`fabricExtension.zip`](../fabricExtension.zip) by using the script [`zip-fabric.sh`](../zip-fabric.sh), and then upload the zip-file to the `TIBCO Flogo® Enterprise 2.8` as an extension, and so they are available for modeling chaincode.
