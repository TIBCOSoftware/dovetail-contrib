---
title: Composer
weight: 4603
---

# Composer
This connector imports Business Network Archive (.bna) file and converts assets, transactions and events defined in [Hyperledger Composer modeling language](https://hyperledger.github.io/composer/v0.19/reference/cto_language.html) into json schemas, not all features are supported, see below for support and limitations
  * **Features Supported**
    - Resource definitions, including assets, concepts, enums, transactions, events and participants.
    - Relationships
    - Imports
    - Decorators, Project Dovetail™ has a predefine set of decorators
       - Asset Decorators
            * @CordaClass("corda class name") : map concept or asset to Corda type, e.g. @CordaClass("net.corda.core.contracts.LinearState") defined in dovetail.system.cto file
            * @CordaParticipants("$tx.path.to.party1", ""$tx.path.to.party2", "...") : specify participants of corda ContractState, if not present, default to top leve Party reference
            * @PrimaryCompositeKey("list of comma delimited attributes"): provide composite key support, identifier specified in the model will be ignored
            * @SecondaryCompositeKeys("list of comma delimited attributes", "another composite key"): provide multiple composite keys for partial key lookup
       - Transaction Decorators
            * @InitiatedBy("$tx.path.to.authorizedparty.attribute", "attributename=value"): specify which participants are authorized to invoke a transaction
            * @Query() : indicate a transaction is query, will not be included in Corda Contract

  * **Features Not Supported Yet**
    - Field validators

  * **Limited Support**
     - In order to be blockchain agnostic, Project Dovetail™ maps participants to network participant's identity, e.g. MSP identifier for Hyperledger Fabric, and party's legal name for R3 Corda, participants are not stored on the ledger. Project Dovetail™ defines a Party participant in its system namespace, it should be used by all resource definitions where participants are required, and Party should be alreays be by reference since they are exsiting entities.

## Settings
| Setting   | Required | Description |
|:----------|:---------|:------------|
| name      | True     | Common data model name |
| mode      | True     | Either use an existing .bna file which is recommended, or author asset models inline |
| modelFile | True     | Select an existing .bna file from file system  |


