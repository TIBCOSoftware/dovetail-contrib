/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package transaction

type TxnInputAttribute struct {
	Name          string
	DataType      string
	IsAssetRef    bool
	IsArray       bool
	AssetName     string
	Identifiers   string
	IsParticipant bool
}

type TxnACL struct {
	AuthorizedParty []string
	Conditions      map[string]string
}
