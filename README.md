# dovetail-contrib

[![Build Status](https://travis-ci.org/TIBCOSoftware/flogo-contrib.svg?branch=master)](https://travis-ci.org/TIBCOSoftware/flogo-contrib.svg?branch=master)

Collection of Dovetail™ activities, triggers and models.

## Contributions

### Activities


### Triggers

* [json-smartcontract](trigger/json-smartcontract): Start flow via a Json Smart Contract
 
### Models


## Installation

#### Install Activity
Example: install **log** activity

```bash
dovetail add activity github.com/TIBCOSoftware/dovetail-contrib/activity/log
```
#### Install Trigger
Example: install **rest** trigger

```bash
dovetail add trigger github.com/TIBCOSoftware/dovetail-contrib/trigger/rest
```


## Contributing and support

### Contributing

New activites, triggers and models are welcome. If you would like to submit one, contact us via [Slack](https://tibco-cloud.slack.com/messages/dovetail-general/).  Contributions should follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.

## License
dovetail-contrib is licensed under a BSD-type license. See TIBCO LICENSE.txt for license text.

### Support
For Q&A you can post your questions on [Slack](https://tibco-cloud.slack.com/messages/dovetail-general/)

