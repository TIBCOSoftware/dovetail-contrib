/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package transaction

type TransactionService interface {
	ResolveTransactionInput(txnInputsMetadata []TxnInputAttribute) (map[string]interface{}, error)
	GetInitiatorCertAttribute(attr string) (value string, found bool, err error)
	GetTransactionName() string
	GetTransactionInitiator() (string, error)
	TransactionSecuritySupported() bool
}
