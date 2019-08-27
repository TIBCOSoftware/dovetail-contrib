---
title: Collection
weight: 4603
---

# Collection
This activity operates on array of primitive or json objects, supported operations: "COUNT", "DISTINCT", "FILTER", "INDEXING", "MERGE", "REDUCE-JOIN"

## Schema
Inputs and Outputs:

```json
{
  "inputs": [
            {
                "name": "operation",
                "type": "string",
                "required": true,
                "allowed": ["COUNT", "DISTINCT", "FILTER", "INDEXING", "MERGE", "REDUCE-JOIN"]
           },
           {
                "name": "model",
                "type": "string",
                "required": true
            },
           {
               "name": "dataType",
               "type": "string",
               "required":false,
               "allowed": ["String"]
           },
           {
                "name": "delimiter",
                "type": "string",
                "required":false
            },
            {
                "name": "filterFieldType",
                "type": "string",
                "required":false,
                "allowed":[ "Boolean","Integer", "Long", "String"]
            },
            {
                "name": "filterFieldOp",
                "type": "string",
                "required":false,
                "allowed":[ "==",">", ">=", "<", "<=", "!=", "IN"]
            },
           {
                "name": "input",
                "type": "complex_object"
            },
            {
                "name": "userInput",
                "type": "complex_object"
            }
    ],
  
    "outputs": [
            {
                "name": "output",
                "type": "complex_object"
           }
    ]
}
```

## Settings
| Setting         | Required | Description |
|:----------------|:---------|:------------|
| operation       | True     | Supported collection functions: "COUNT", "DISTINCT", "FILTER", "INDEXING", "MERGE", "REDUCE-JOIN", see details below |
| model           | False    | Select the common data model, must be the same as the one selected in Trigger |
| dataType        | True     | Data type, the available data types depend on operation and if a common data model is selected. "User Defined..." data type allows user defined json schema |
| delimiter       | False    | Required for REDUCE-JOIN  |
| filterFieldType | False    | Required for FILTER, specify the data type of field field |
| filterFieldOp   | False    | Required for FILTER, specify the filter operation, supported operators: "==",">", ">=", "<", "<=", "!=", "IN"|

| Operation     | Description |
|:--------------|:------------|
| COUNT         | Return the size of an array of objects|
| DISTINCT      | Return a list of distinct string values from a collection object |
| FILTER        | Filter a collection into true set and false set based on filter conditions, filterField should be in the format of "$dataset.path.to.field, e.g. $dataset.myobject.myfield |
| INDEXING      | Assign a position number to each item in the collection starting with 0, the number is stored in _index_ field |
| MERGE         | Merge to two collections of the same object type into one |
| REDUCE-JOIN   | Concat a collection of strings with a delimiter |

## Input Schema
The json schema is automatically created based on data type selected.

## Ouput Schema
The json schema is automatically created based on data type selected.


