---
title: transaction filter
weight: 4603
---

# Transaction Filter
Filter out input/output/reference states, commands, notary or time window to verify a transaction before signing, used in [responder flow](../../trigger/flowreceiver/README.md)

## Settings
| Setting       | Required | Description                                                                       |
|:--------------|:---------|:----------------------------------------------------------------------------------|
| filter        | True     | Select what data to filter from the ledger transaction                            |
| artifact      | False    | Select the artifact type to filter on                                             |
