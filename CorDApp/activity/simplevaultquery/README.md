---
title: simple vault query
weight: 4603
---

# Simple Vault Query
Query vault for linear state assets

The query output is an array of objects with structure of data and ref. data object contains the value of the state, ref is a pointer to the state. "ref" is used to map to BuildTransactionProposal activity input attributes that are type of asset by reference

## Settings
| Setting       | Required | Description                                                                       |
|:--------------|:---------|:----------------------------------------------------------------------------------|
| model         | True     | Select contract model                                                             |
| assetName     | True     | Select asset to query                                                             |
| status        | True     | Asset status, default to UNCONSUMED                                               |
| assetType     | True     | default to LinearState                                                            | 
