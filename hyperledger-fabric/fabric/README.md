# Flogo extension for Hyperledger Fabric chaincode

This Flogo extension is designed to allow delopers to design and implement Hyperledger Fabric chaincode in the zero-code visual programming environment of the TIBCO Flogo® Enterprise.  This extension supports the following release versions:
- [TIBCO Flogo® Enterprise 2.4](https://docs.tibco.com/products/tibco-flogo-enterprise-2-4-0)
- [Hyperledger Fabric 1.4](https://www.hyperledger.org/projects/fabric)

The [Transaction Trigger](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/trigger/transaction) allows you to configure the chaincode input and output schema, including persistent and/or transient input parameters.

It supports the following activities for storing and querying data on the distributed ledger and/or on private collections.
- [Put](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/put): Insert or update data on the distributed ledger or a private data collection, and optionally insert its associated compsite keys if they are specified.
- [Put All](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/putall): Insert a list of records on the distributed ledger or a private collection, and optionally insert composite keys of each record.
- [Get](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/get): Retrieve a record of a specified key from the distributed ledger or a private collection.
- [Get by Range](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/getrange): Retrieve all records in a specified range of keys from the distributed ledger or a private collection.  It supports resultset pagination for records from the distributed ledger.
- [Get by Composite Key](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/getbycompositekey): Retrieve all records by a composite-key filter from the distributed ledger or a private collection.  It supports resultset pagination for records from the distributed ledger.
- [Get History](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/gethistory): Retrieve the history of a specified key for data on the distributed ledger.
- [Query](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/query): Retrieve all records by a Couchdb query statement from the distributed ledger or a private collection.  It supports resultset pagination for records from the distributed ledger.
- [Delete](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/delete): Delete the state of a specified key from the distributed ledger or a private collection, as well as its composite keys.  Optionally, it can delete only the state, or only a composite key.
- [Set Event](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/setevent): Set a specified event and payload for a blockchain transaction.
- [Set Endorsement Policy](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/endorsement): Set state-based endorsement policy by adding or deleting an endorsement organization, or by specifying a new endorsement policy.
- [Invoke Chaincode](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/fabric/activity/endorsement): Invoke a local chaincode, and returns response data from the called transaction.

With these extensions, Hyperledger Fabric chaincode can be designed and implemented with zero code. Refer to samples [`marble-app`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/marble-app) and [`marble-private`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/tree/master/marble-private) for more details of the chaincode models implemented by using the `TIBCO Flogo® Enterprise`.

This extension is packaged as [`fabricExtension.zip`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/fabricExtension.zip), which can be re-created from source by using the script [`zip-fabric.sh`](https://github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/blob/master/zip-fabric.sh).  You can upload the zip-file to the `TIBCO Flogo® Enterprise` as an extension, and start using the trigger and activities.
