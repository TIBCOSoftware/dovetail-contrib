---
title: Publisher
weight: 4603
---

# Publisher
This activity publish events, for blockchain that does not support eventing, this activity has no effect.

## Schema
Inputs and Outputs:

```json
{
  "inputs": [
            {
                "name": "model",
                "type": "string"
            },
           {
                "name": "event",
                "type": "string"
           },
           {
            "name": "input",
            "type": "complex_object",
            "required": true
           },
           {
                "name": "eventMetadata",
                "type": "string",
                "required":false
           },
           {
               "name": "containerServiceStub",
               "type": "any",
               "required":true
           }
    ],
  
    "outputs": [
        
    ]
}
```

## Settings
| Setting               | Required | Description |
|:----------------------|:---------|:------------|
| model                 | True     | Select the common data model, must be the same as the one selected in Trigger |
| event                 | True     | Select the event name defined in common data model to publsih |
| eventMetadata         | False    | Event metadata  |
| containerServiceStub  | True     | This is the handler to underlying blockchain service, should always be mapped to $flow.containerServiceStub |

## Input Schema
The json schema is automatically created based on settings

## Ouput Schema
The json schema is automatically created based on settings


