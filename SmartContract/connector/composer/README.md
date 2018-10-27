---
title: Composer
weight: 4603
---

# Composer
This connector parses Hyperledger Composer .bna file to generate common assets, concepts, events and transactions json schemas to be used in activities

## Settings
| Setting   | Required | Description |
|:----------|:---------|:------------|
| name      | True     | Common data model name |
| mode      | True     | Either use an existing .bna file which is recommended, or author asset models inline |
| modelFile | True     | Select an existing .bna file from file system  |


