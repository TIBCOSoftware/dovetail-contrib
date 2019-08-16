---
title: time window
weight: 4603
---

# Time Window
Use the activity to set the time window when a transaction can be performed. The datetime string must represent a valid instant in UTC and is parsed using DateTimeFormatter.ISO_INSTANT, e.g. 2007-12-03T10:15:30.00Z.

For "Only valid for the duration of...", if "from" time is not set, current system time will be used.

## Settings
| Setting       | Required | Description                                                                       |
|:--------------|:---------|:----------------------------------------------------------------------------------|
| window        | True     | Select window constraint type                                                     |



