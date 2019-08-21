---
title: Transaction
weight: 4603
---

# Transaction
At design time, it works with Hyperledger Composer Connector to display predefined user types for modeling, at runtime, it receives smart contract transactions from distributed ledger platform, resolves transaction input, dispatches transactions to appropriate flow handler, and sends transaction reponses, if any, back to caller.

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| model       | True     | Common data model name |
| createAll   | True     | Create flows for all transactions defined in the model, or select a specific transaction |
| transaction | True     | Select the transaction to implement |


