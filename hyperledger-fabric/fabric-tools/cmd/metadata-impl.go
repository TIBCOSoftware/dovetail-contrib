// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// TransactionSchemas describes data schemas in metadata of a transaction
type TransactionSchemas struct {
	Parameters string
	Transient  string
	Returns    string
}

// SchemaValue describes json schema of a parameter value
type SchemaValue struct {
	Ref        string          `json:"$id,omitempty"`
	Type       string          `json:"type"`
	Properties json.RawMessage `json:"properties,omitempty"`
	Items      json.RawMessage `json:"items,omitempty"`
}

var schemaCache map[string]SchemaValue
var sharedSchema map[string]string

func init() {
	schemaCache = make(map[string]SchemaValue)
	sharedSchema = make(map[string]string)
}

func generateMetadata(appfile, outpath string) error {
	appconfig, err := ioutil.ReadFile(appfile)
	if err != nil {
		return err
	}

	// extract contract info from flogo app config
	var contractData struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Triggers []struct {
			Ref      string `json:"ref"`
			Handlers []struct {
				Description string `json:"description"`
				Settings    struct {
					Name string `json:"name"`
				} `json:"settings"`
				Schemas map[string]json.RawMessage `json:"schemas"`
			} `json:"handlers"`
		} `json:"triggers"`
	}
	if err := json.Unmarshal(appconfig, &contractData); err != nil {
		return errors.Wrapf(err, "failed to extract data from flogo config file %s", appfile)
	}

	// select all handlers of the first fabric transaction trigger
	for _, trig := range contractData.Triggers {
		if trig.Ref == "#transaction" {
			for _, handler := range trig.Handlers {
				fmt.Printf("name %s\n", handler.Settings.Name)
				txnSchemas := extractTransactionSchema(handler.Schemas)
				fmt.Printf("schemas: %+v\n", txnSchemas)
			}
			break
		}
	}
	for k, v := range sharedSchema {
		fmt.Printf("shared schema type %s: %s\n", v, k)
	}
	return nil
}

func extractTransactionSchema(schemas map[string]json.RawMessage) TransactionSchemas {
	var parameters, transient, returns string
	if reply, ok := schemas["reply"]; ok {
		var r struct {
			Returns struct {
				Value string `json:"value"`
			} `json:"returns"`
		}
		if err := json.Unmarshal(reply, &r); err != nil {
			fmt.Printf("failed to unmarshal reply: %+v\n", err)
		} else {
			if r.Returns.Value != "" {
				returns = addSchemaToCache(r.Returns.Value)
			}
		}
	}
	if output, ok := schemas["output"]; ok {
		var o struct {
			Parameters struct {
				Value string `json:"value"`
			} `json:"parameters,omitempty"`
			Transient struct {
				Value string `json:"value"`
			} `json:"transient,omitempty"`
		}
		if err := json.Unmarshal(output, &o); err != nil {
			fmt.Printf("failed to unmarshal output: %+v\n", err)
		} else {
			if o.Parameters.Value != "" {
				parameters = addSchemaToCache(o.Parameters.Value)
			}
			if o.Transient.Value != "" {
				transient = addSchemaToCache(o.Transient.Value)
			}
		}
	}
	return TransactionSchemas{Parameters: parameters, Transient: transient, Returns: returns}
}

func addSchemaToCache(value string) string {
	var sv SchemaValue
	if err := json.Unmarshal([]byte(value), &sv); err != nil {
		fmt.Printf("failed to unmarshal transient value: %+v\n", err)
	} else {
		if svbytes, err := json.Marshal(sv); err != nil {
			fmt.Printf("failed to marshal transient schema: %+v\n", err)
		} else {
			skey := string(svbytes)
			// add to cache
			if _, ok := schemaCache[skey]; ok {
				sharedSchema[skey] = sv.Type
			} else {
				schemaCache[skey] = sv
			}
			return skey
		}
	}
	return ""
}
