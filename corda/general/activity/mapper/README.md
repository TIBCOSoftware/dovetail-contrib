---
title: Mapper
weight: 4603
---

# Mapper
This activity performs generic mapping activity, also support array type switching, from primitive array, such as ["a", "b"] to object array [{"field": "a"}, {"field", "b"}], or vice versa, to overcome UI array mapping limitation.

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
            "name": "dataType",
            "type": "string",
            "required": true,
            "allowed": ["Boolean","Datetime", "Double", "Integer", "Long", "String"]
        },
        {
            "name": "isArray",
            "type": "boolean"
           },
           {
            "name": "inputArrayType",
            "type": "string",
            "allowed": ["Object Array", "Primitive Array"]
           },
           {
            "name": "outputArrayType",
            "type": "string",
            "allowed": ["Object Array", "Primitive Array"]
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
            "type": "complex_object",
            "required": true
        }
    ]
}
```

## Settings
| Setting         | Required | Description |
|:----------------|:---------|:------------|
| model           | False    | Select the common data model, must be the same as the one selected in Trigger |
| dataType        | True     | Data type, primitive data types plus types defined in common data model if a model is selected. "User Defined..." data type allows user defined json schema |
| isArray         | True     | True if input is an array  |
| inputArrayType  | False    | Required if isArray is set to true, example of "Object Array": [{"field":"value"1},{"field", "value2"}], example of "Primitive Array": ["value1", "value2"]
| outputArrayType | False    | Required if isArray is set to true, see examples above

## Input Schema
The json schema is automatically created based on settings

## Ouput Schema
The json schema is automatically created based on settings


