---
title: flow receiver
weight: 4603
---

# Flow Receiver
This trigger defines a responding or observer flow initiated by an initiating flow defined by Flow Initiator trigger. If a transaction is sent by com.tibco.dovetail.container.cordapp.flows.ObserverFlowInitiator, there is no need to define a receiver flow, com.tibco.dovetail.container.cordapp.flows.ObserverFlowReceiver will be used.

## Settings
| Setting               | Required | Description |
|:------------          |:---------|:------------| 
| flowType              | True     | 
| useAnonymousIdentity  | True     | Select true if annonymous identity must be used for this transaction |
| hasObservers          | True     | Select true if the transaction has observers |
| initiatorFlow         | True     | The initiating flow's class name, including package |


