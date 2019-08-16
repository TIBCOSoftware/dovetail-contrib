---
title: flow initiator
weight: 4603
---

# Flow Initiator
This trigger starts an initiating flow

## Settings
| Setting               | Required | Description |
|:------------          |:---------|:------------|
| useAnonymousIdentity  | True     | Select true if annonymous identity must be used for this transaction |
| hasObservers          | True     | Select true if the transaction has observers |
| observerManual        | True     | If there are observers, specify if this transaction should be sent in this flow or a separate flow. when a separate flow is selected, com.tibco.dovetail.container.cordapp.flows.ObserverFlowInitiator will be used. |


