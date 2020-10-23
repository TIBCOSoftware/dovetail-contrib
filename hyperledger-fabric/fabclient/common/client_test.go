/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package common

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test requires to start byfn fabric network using "byfn.sh up -a -s couchdb"
//   and set the fabPath below to the absolute path of "fabric-samples"
const (
	connectorName = "test"
	fabPath       = "/Users/yxu/work/dovetail/fabric-samples"
	configFile    = "${HOME}/work/dovetail/dovetail-contrib/hyperledger-fabric/testdata/config_test.yaml"
	matcherFile   = "${HOME}/work/dovetail/dovetail-contrib/hyperledger-fabric/testdata/local_entity_matchers.yaml"
	channelID     = "mychannel"
	org           = "org1"
	user          = "User1"
	ccID          = "mycc"
)

func TestClient(t *testing.T) {
	os.Setenv("FAB_PATH", fabPath)
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
	result, _, err := fbClient.QueryChaincode(ccID, "query", [][]byte{[]byte("a")}, nil)
	require.NoError(t, err, "failed to query %s", ccID)
	fmt.Printf("Query result: %s\n", string(result))
	origValue := result

	// update
	result, _, err = fbClient.ExecuteChaincode(ccID, "invoke", [][]byte{[]byte("a"), []byte("b"), []byte("10")}, nil)
	require.NoError(t, err, "failed to invoke %s", ccID)
	fmt.Printf("Invoke result: %s\n", string(result))

	// query after update
	result, _, err = fbClient.QueryChaincode(ccID, "query", [][]byte{[]byte("a")}, nil)
	require.NoError(t, err, "failed to query %s", ccID)
	fmt.Printf("Query result: %s\n", string(result))
	assert.NotEqual(t, origValue, result, "original %s should different from %s", string(origValue), string(result))
}
