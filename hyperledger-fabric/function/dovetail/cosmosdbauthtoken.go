package dovetail

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/support/log"
)

// CosmosdbAuthToken dummy struct
type CosmosdbAuthToken struct {
}

func init() {
	function.Register(&CosmosdbAuthToken{})
}

// Name of function
func (s *CosmosdbAuthToken) Name() string {
	return "cosmosdbAuthToken"
}

// Sig - function signature
func (s *CosmosdbAuthToken) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString, data.TypeString, data.TypeString, data.TypeString}, false
}

// Eval - function implementation
func (s *CosmosdbAuthToken) Eval(params ...interface{}) (interface{}, error) {

	log.RootLogger().Debugf("Start cosmosdbAuthToken function with params %+v", params)

	verb, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("verb %T is not a string", params[0])
	}

	reqURI, ok := params[1].(string)
	if !ok {
		return nil, fmt.Errorf("reqURI %T is not a string", params[1])
	}
	resourceType, resourceID := parseRequestURI(reqURI)

	utc, ok := params[2].(string)
	if !ok {
		return nil, fmt.Errorf("utc %T is not a date string", params[2])
	}

	masterKey, ok := params[3].(string)
	if !ok {
		return nil, fmt.Errorf("master-key %T is not a string", params[3])
	}

	// decode master key from base64
	key, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode master key: %v", err)
	}

	// sign request using master key, and encode signature in base64
	text := fmt.Sprintf("%s\n%s\n%s\n%s\n\n", strings.ToLower(verb), strings.ToLower(resourceType), strings.ToLower(resourceID), strings.ToLower(utc))
	//fmt.Println("text=", text)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(text))
	sig := base64.StdEncoding.EncodeToString(h.Sum(nil))
	//fmt.Println("sig=", sig)

	// construct authorization token
	auth := fmt.Sprintf("type=master&ver=1.0&sig=%s", sig)
	return url.QueryEscape(auth), nil
}

// return resourceType and resourceID from request URI of a CosmosDB REST API
func parseRequestURI(uri string) (string, string) {
	trimedURI := strings.TrimSpace(uri)
	tokens := strings.Split(trimedURI, "/")
	size := len(tokens)
	if size == 0 {
		return "", ""
	}

	// for even number, second to last token is resourceType
	if size%2 == 0 {
		return tokens[size-2], trimedURI
	}

	// for odd number, last token is resourceType
	if size == 1 {
		return tokens[0], ""
	}
	indx := strings.LastIndex(trimedURI, "/")
	return tokens[size-1], trimedURI[0:indx]
}
