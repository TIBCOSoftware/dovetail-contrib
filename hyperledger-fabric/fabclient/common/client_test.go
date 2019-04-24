/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	connectorName = "test"
	configFile    = "${GOPATH}/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/testdata/config_test.yaml"
	matcherFile   = "${GOPATH}/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/testdata/local_entity_matchers.yaml"
	channelID     = "mychannel"
	org           = "org1"
	user          = "User1"
	ccID          = "mycc"
)

func TestClient(t *testing.T) {
	networkConfig, err := ReadFile(configFile)
	require.NoError(t, err, "failed to read config file %s", configFile)

	entityMatcherOverride, err := ReadFile(matcherFile)
	require.NoError(t, err, "failed to read entity matcher file %s", matcherFile)

	fbClient, err := NewFabricClient(ConnectorSpec{
		Name:           connectorName,
		NetworkConfig:  networkConfig,
		EntityMatchers: entityMatcherOverride,
		OrgName:        org,
		UserName:       user,
		ChannelID:      channelID,
	})
	require.NoError(t, err, "failed to create fabric client %s", connectorName)
	fmt.Printf("created fabric client %+v\n", fbClient)

	// query original
	result, err := fbClient.QueryChaincode(ccID, "query", [][]byte{[]byte("a")})
	require.NoError(t, err, "failed to query %s", ccID)
	fmt.Printf("Query result: %s\n", string(result))
	origValue := result

	// update
	result, err = fbClient.ExecuteChaincode(ccID, "invoke", [][]byte{[]byte("a"), []byte("b"), []byte("10")})
	require.NoError(t, err, "failed to invoke %s", ccID)
	fmt.Printf("Invoke result: %s\n", string(result))

	// query after update
	result, err = fbClient.QueryChaincode(ccID, "query", [][]byte{[]byte("a")})
	require.NoError(t, err, "failed to query %s", ccID)
	fmt.Printf("Query result: %s\n", string(result))
	assert.NotEqual(t, origValue, result, "original %s should different from %s", string(origValue), string(result))
}
