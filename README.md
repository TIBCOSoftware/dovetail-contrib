# dovetail-contrib

[![Build Status](https://travis-ci.org/TIBCOSoftware/dovetail-contrib.svg?branch=master)](https://travis-ci.org/TIBCOSoftware/dovetail-contrib.svg?branch=master)

Collection of Dovetail™ activities, triggers and models.

## Contributions

### Connectors

* [composer](SmartContract/composer)

### Activities

* [aggregate](SmartContract/aggregate)
* [collection](SmartContract/collection)
* [history](SmartContract/history)
* [ledger](SmartContract/ledger)
* [logger](SmartContract/logger)
* [mapper](SmartContract/mapper)
* [publisher](SmartContract/publisher)
* [query](SmartContract/query)
* [txnreply](SmartContract/txnreply)

### Triggers

* [transaction](SmartContract/transaction)

## Installation

For step by step instructions on how to install a new activity and trigger please go to the [documentation page](https://tibcosoftware.github.io/dovetail/)


### Contributing

New activites, triggers and models are welcome. If you would like to submit one, contact us via email at tibcolabs@tibco.com .  Contributions should follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.

## License
dovetail-contrib is licensed under a BSD-type license. See [LICENSE](https://github.com/TIBCOSoftware/dovetail-contrib/blob/master/LICENSE) for license text.

### Support
For Q&A you can contact us at tibcolabs@tibco.com.
