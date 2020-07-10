package dovetail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var b = &CosmosdbAuthToken{}

func TestParseRequestURI(t *testing.T) {
	// test odd tokens
	res := "dbs/demodb/colls/democoll/docs"
	resType, resID := parseRequestURI(res)
	assert.Equal(t, "docs", resType, "resource type for uri %s: %s should be docs", res, resType)
	assert.Equal(t, "dbs/demodb/colls/democoll", resID, "resource ID for uri %s: %s should be dbs/demodb/colls/democoll", res, resID)

	// test single token
	res = "dbs"
	resType, resID = parseRequestURI(res)
	assert.Equal(t, "dbs", resType, "resource type for uri %s: %s should be dbs", res, resType)
	assert.Equal(t, "", resID, "resource ID for uri %s: %s should be \"\"", res, resID)

	// test even tokens
	res = "dbs/demodb/colls/democoll"
	resType, resID = parseRequestURI(res)
	assert.Equal(t, "colls", resType, "resource type for uri %s: %s should be colls", res, resType)
	assert.Equal(t, res, resID, "resource ID for uri %s: %s should be %s", res, resID, res)
}

func TestAuthToken(t *testing.T) {
	uri := "dbs/jbademodb"
	utc := "Thu, 09 Jul 2020 16:19:51 GMT"
	masterKey := "0ejqFuQrwF2xgUamnHXud3RFYIMXbq7kaTj0ysU0b9Z83X8IY710UInJIqxRbqXIlHBPbjaWTb3aBdrGDpya2w=="
	v, err := b.Eval("GET", uri, utc, masterKey)
	assert.NoError(t, err, "cosmosdbAuthToken should not return error %v", err)
	b, ok := v.(string)
	assert.True(t, ok, "cosmosdbAuthToken should return string, not type %T", v)
	assert.Equal(t, "type%3Dmaster%26ver%3D1.0%26sig%3DjY50HSpAMkJTKRFJ67l4g9Aj2hf1hl%2BtYVcJ0uxaNaM%3D", b, "cosmosdbAuthToken should not return value %s", b)
}
