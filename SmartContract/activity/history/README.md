---
title: Hyperledger Fabric Historical Records
weight: 4603
---

# History
This activity retrieves all historical transactios for a particular asset identifier, Hyperledger Fabric only.

## Schema
Inputs and Outputs:

```json
{
  "inputs": [
            {
                "name": "model",
                "type": "string",
                "required": true
            },
           {
                "name": "assetName",
                "type": "string",
                "required": true
           },
           {
            "name": "input",
            "type": "complex_object",
            "required": true
           },
           {
            "name": "containerServiceStub",
            "type": "any",
            "required":true
           }
    ],
  
    "outputs": [
        {
            "name": "output",
            "type": "complex_object",
            "required": true
        }
    ]
}
```

## Settings
| Setting              | Required | Description |
|:---------------------|:---------|:------------|
| model                | True    | Select the common data model, must be the same as the one selected in Trigger |
| assetName            | True    | Select the asset to query|
| containerServiceStub | True    | This is the handler to underlying blockchain service, should always be mapped to $flow.containerServiceStub |

## Input Schema
The json schema is automatically created

## Ouput Schema
The json schema is automatically created 


