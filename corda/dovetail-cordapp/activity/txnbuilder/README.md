---
title: transaction builder
weight: 4603
---

# Transaction Builder
This is the main activity of initiator flows to build transaction proposal to be signed and recorded to ledger. It invokes the smart contract transaction to create output states, commands, and signature requirements

Smart contract transaction input attributes that are type of asset by reference (e.g. -->IOU iou) should be mapped to "ref" field of vault query activitoies, e.g. SimpleVaultQuery and CashWallet.

## Settings
| Setting       | Required | Description                                                                       |
|:--------------|:---------|:----------------------------------------------------------------------------------|
| contract      | True     | Select the smart contract model                                                   |
| transaction   | true     | Select the transaction                                                            |


