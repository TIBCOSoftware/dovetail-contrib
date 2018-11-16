---
title: Project Dovetail™ Smart Contract Action Model Component
weight: 4603
---
# Project Dovetail™ Smart Contract Action Model Component

Collection of Smart Contract activities, triggers and connectors 

## Contributions

### Activities
* [aggregate](activity/aggregate): Aggregate of numeric values
* [collection](activity/collection): Support operations on a collection of primitive and object data types 
* [history](activity/history): Get historical transactions, Hyperledger Fabric only
* [ledger](activity/ledger): Read and write access to ledger 
* [logger](activity/logger): Simple Logger for blockchain
* [mapper](activity/mapper): Simple mapper
* [query](activity/query): Custom CouchDB Query, Hyperledger Fabric only
* [txreply](activity/txreply): Send transaction reply

### Triggers
* [transaction](trigger/transaction): Smart contract transaction trigger. 

### Connectors
Connectors can be used to bring in business data types and structs defined in the tools of your choice and convert them into json schema used by Project Dovetail™ Studio. 

* [composer](connector/composer): Hyperledger Composer connector

## Installation

* From SmartContract directory, run ```zip  -u -r contrib-smartcontract.zip SmartContract/*```
* Upload DovetailSmartContractExtension.zip file through Project Dovetail™ Studio extension tab.


