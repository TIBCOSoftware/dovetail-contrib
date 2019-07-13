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

// Transaction describes data schemas in metadata of a transaction
type Transaction struct {
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Operation   string          `json:"operation,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	Transient   json.RawMessage `json:"transient,omitempty"`
	Returns     json.RawMessage `json:"returns,omitempty"`
}

// SchemaValue describes json schema of a parameter value
type SchemaValue struct {
	Ref        string          `json:"$id,omitempty"`
	Type       string          `json:"type"`
	Properties json.RawMessage `json:"properties,omitempty"`
	Items      json.RawMessage `json:"items,omitempty"`
}

// Contract describes json schema of a contract
type Contract struct {
	Name         string         `json:"name"`
	Transactions []*Transaction `json:"transactions"`
}

// Info describes metadata info
type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// Metadata describes smart contracts of an app
type Metadata struct {
	Schema     string               `json:"$schema"`
	Info       Info                 `json:"info"`
	Contracts  map[string]*Contract `json:"contracts"`
	Components struct {
		Schemas map[string]SchemaValue `json:"schemas,omitempty"`
	} `json:"components,omitempty"`
}

var schemaCache map[string]*SchemaValue
var sharedSchema map[string]string

func init() {
	schemaCache = make(map[string]*SchemaValue)
	sharedSchema = make(map[string]string)
}

// generate contract metadata from flogo model
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

	// construct output metadata
	metadata := Metadata{
		Schema:    "http://json-schema.org/draft-04/schema#",
		Info:      Info{Title: contractData.Name, Version: contractData.Version},
		Contracts: map[string]*Contract{},
	}
	contract := Contract{Name: contractData.Name}
	metadata.Contracts[contract.Name] = &contract

	// select all handlers of the first fabric transaction trigger
	for _, trig := range contractData.Triggers {
		if trig.Ref == "#transaction" {
			for _, handler := range trig.Handlers {
				fmt.Printf("transaction name %s\n", handler.Settings.Name)
				txn := extractTransactionSchema(handler.Schemas)
				txn.Name = handler.Settings.Name
				txn.Description = handler.Description
				// TODO: determine if it is invoke or query
				txn.Operation = "invoke"
				contract.Transactions = append(contract.Transactions, &txn)
			}
			break
		}
	}

	// set unique id for shared schema
	setSharedSchemaRef(contractData.Name)
	if len(sharedSchema) > 0 {
		metadata.Components.Schemas = map[string]SchemaValue{}
		for k := range sharedSchema {
			sv := schemaCache[k]
			metadata.Components.Schemas[sv.Ref] = *sv
		}
	}

	// use shared schema in transacions
	for _, txn := range contract.Transactions {
		replaceSharedSchema(txn)
	}
	metabytes, _ := json.MarshalIndent(metadata, "", "    ")
	fmt.Println("write metadata to", outpath)
	return writeFile(outpath, metabytes)
}

func replaceSharedSchema(txn *Transaction) {
	if txn.Parameters != nil {
		txn.Parameters = toSharedSchema(txn.Parameters)
	}
	if txn.Transient != nil {
		txn.Transient = toSharedSchema(txn.Transient)
	}
	if txn.Returns != nil {
		txn.Returns = toSharedSchema(txn.Returns)
	}
}

func toSharedSchema(schema []byte) []byte {
	svkey := string(schema)
	if _, ok := sharedSchema[svkey]; ok {
		sv := schemaCache[svkey]
		shared := fmt.Sprintf(`{
			"$ref": "#/components/schemas/%s"
		}`, sv.Ref)
		return []byte(shared)
	}
	return schema
}

func extractTransactionSchema(schemas map[string]json.RawMessage) Transaction {
	var parameters, transient, returns []byte
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
	return Transaction{Parameters: parameters, Transient: transient, Returns: returns}
}

func addSchemaToCache(value string) []byte {
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
				schemaCache[skey] = &sv
			}
			return svbytes
		}
	}
	return nil
}

func setSharedSchemaRef(appname string) {
	arraySeq := 0
	objectSeq := 0
	for k := range sharedSchema {
		sv := schemaCache[k]
		if sv.Type == "array" {
			sv.Ref = fmt.Sprintf("%s_array%d", appname, arraySeq)
			arraySeq++
		} else {
			if n, err := objectPropertyCount(sv.Properties); err == nil && n > 2 {
				sv.Ref = fmt.Sprintf("%s_%d", appname, objectSeq)
				objectSeq++
			} else {
				delete(sharedSchema, k)
				if err != nil {
					fmt.Printf("remove shared schema %s due to error: %+v\n", k, err)
				}
			}
		}
	}
}

func objectPropertyCount(schema []byte) (int, error) {
	var props map[string]interface{}
	if err := json.Unmarshal(schema, &props); err != nil {
		return -1, errors.Wrapf(err, "failed to unmarshal object properties: %s", string(schema))
	}
	return len(props), nil
}
