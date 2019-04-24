# Flogo extension for Hyperledger Fabric chaincode

This Flogo extension is designed to allow delopers to design and implement Hyperledger Fabric chaincode in the zero-code visual programming environment of the TIBCO Flogo® Enterprise.  This extension supports the following release versions:
- [TIBCO Flogo® Enterprise 2.4](https://docs.tibco.com/products/tibco-flogo-enterprise-2-4-0)
- [Hyperledger Fabric 1.4](https://www.hyperledger.org/projects/fabric)

The [Transaction Trigger](trigger/transaction) allows you to configure the chaincode input and output schema, including persistent and/or transient input parameters.

It supports the following activities for storing and querying data on the distributed ledger and/or on private collections.
- [Put](activity/put): Insert or update data on the distributed ledger or a private data collection, and optionally insert its associated compsite keys if they are specified.
- [Put All](activity/putall): Insert a list of records on the distributed ledger or a private collection, and optionally insert composite keys of each record.
- [Get](activity/get): Retrieve a record of a specified key from the distributed ledger or a private collection.
- [Get by Range](activity/getrange): Retrieve all records in a specified range of keys from the distributed ledger or a private collection.  It supports resultset pagination for records from the distributed ledger.
- [Get by Composite Key](activity/getbycompositekey): Retrieve all records by a composite-key filter from the distributed ledger or a private collection.  It supports resultset pagination for records from the distributed ledger.
- [Get History](activity/gethistory): Retrieve the history of a specified key for data on the distributed ledger.
- [Query](activity/query): Retrieve all records by a Couchdb query statement from the distributed ledger or a private collection.  It supports resultset pagination for records from the distributed ledger.
- [Delete](activity/delete): Delete the state of a specified key from the distributed ledger or a private collection, as well as its composite keys.  Optionally, it can delete only the state, or only a composite key.
- [Set Event](activity/setevent): Set a specified event and payload for a blockchain transaction.
- [Set Endorsement Policy](activity/endorsement): Set state-based endorsement policy by adding or deleting an endorsement organization, or by specifying a new endorsement policy.
- [Invoke Chaincode](activity/endorsement): Invoke a local chaincode, and returns response data from the called transaction.

With these extensions, Hyperledger Fabric chaincode can be designed and implemented with zero code. Refer to samples [`marble-app`](../marble-app) and [`marble-private`](../marble-private) for more details of the chaincode models implemented by using the `TIBCO Flogo® Enterprise`.

This extension is packaged as [`fabricExtension.zip`](../fabricExtension.zip), which can be re-created from source by using the script [`zip-fabric.sh`](../zip-fabric.sh).  You can upload the zip-file to the `TIBCO Flogo® Enterprise` as an extension, and start using the trigger and activities.
