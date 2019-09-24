---
title: Aggregate
weight: 4603
---

# Aggregate
This activity allows you to aggregate of numeric values. Supported operations are sum, avg, min and max

## Schema
Inputs and Outputs:

```json
{
  "inputs": [
            {
                "name": "operation",
                "type": "string",
                "required": true,
                "allowed": ["MIN", "MAX", "AVG", "SUM"]
           },
           {
               "name": "dataType",
               "type": "string",
               "required":false,
               "allowed": ["Integer", "Long", "Double"]
           },
           {
            "name": "precision",
            "type": "integer",
            "required":false
            },
           {
                "name": "scale",
                "type": "integer",
                "required":false
            },
            {
                "name": "rounding",
                "type": "string",
                "required":false,
                "allowed": ["UP", "DOWN", "CEILING", "FLOOR", "HALF_UP", "HALF_DOWN", "HALF_EVEN"]
            },
           {
                "name": "input",
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
| Setting     | Required | Description |
|:------------|:---------|:------------|
| operation   | True     | The aggregate function, sum, avg, min and max are supported |
| dataType    | True     | Data type, Integer, Long and Double are supported |
| precision   | False    | Required for Double data type, specify the precision of aggregate result |
| scale       | False    | Required for Double data type, specify the scale of aggregate result |
| rounding    | False    | Required for Double data type, specify the rounding mode of aggregate result, see below rounging table for example |

## Result of rounding input to one digit with the given rounding mode
	
| Input Number | UP	| DOWN | CEILING | FLOOR | HALF_UP | HALF_DOWN | HALF_EVEN |
|:-------------|----|------|---------|-------|---------|-----------|-----------|
| 5.5	         | 6	| 5	   | 6	     | 5	   | 6	     | 5	       | 6         |
| 2.5	         | 3	| 2	   | 3	     | 2	   | 3       | 2	       | 2	       |
| 1.6	         | 2	| 1	   | 2	     | 1	   | 2	     | 2	       | 2	       |
| 1.1	         | 2	| 1	   | 2	     | 1	   | 1       | 1	       | 1	       |
| 1.0	         | 1	| 1	   | 1	     | 1	   | 1	     | 1	       | 1	       |
| -1.0	       | -1	| -1	 | -1	     | -1	   | -1	     | -1	       | -1	       |
| -1.1	       | -2	| -1	 | -1	     | -2	   | -1	     | -1	       | -1	       |
| -1.6	       | -2	| -1	 | -1	     | -2	   | -2	     | -2	       | -2        |
| -2.5	       | -3	| -2	 | -2	     | -3	   | -3	     | -2	       | -2	       |
| -5.5	       | -6	| -5	 | -5	     | -6	   | -6	     | -5	       | -6	       |

## Input Schema
The json schema is automatically created based on data type selected.

## Ouput Schema
The json schema is automatically created based on data type selected.
