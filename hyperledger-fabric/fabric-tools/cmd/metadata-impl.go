// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

// example for generating metadata and graphql files:
// go install
// cd ../samples/equinix/contract-metadata
// fabric-tools metadata -f ../equinix.json -o ./override.json

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// TODO: change this to http://tibcosoftware.github.io/dovetail/schemas/contract-schema.json
const metaSchema = "http://json-schema.org/draft-07/schema#"

var schemaOverride map[string]string

func setSchemaOverride(overridefile string) {
	schemaOverride = make(map[string]string)
	override, err := ioutil.ReadFile(overridefile)
	if err != nil {
		return
	}
	if err := json.Unmarshal(override, &schemaOverride); err != nil {
		fmt.Printf("failed to parse schema override file %s: %+v\n", overridefile, err)
	}
}

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
	Name         string                  `json:"name"`
	Transactions map[string]*Transaction `json:"transactions"`
}

// Info describes metadata info
type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// Metadata describes smart contracts of an app
type Metadata struct {
	Schema     string                 `json:"$schema"`
	Info       Info                   `json:"info"`
	Contract   *Contract              `json:"contract"`
	Components map[string]SchemaValue `json:"components,omitempty"`
}

var schemaCache map[string]*SchemaValue
var sharedSchema map[string]int
var schemaSeq int
var appSchema = make(map[string]string)

func init() {
	schemaCache = make(map[string]*SchemaValue)
	sharedSchema = make(map[string]int)
}

// generate contract metadata from flogo model
func generateMetadata(appfile, outpath string) error {
	appconfig, err := ioutil.ReadFile(appfile)
	if err != nil {
		return err
	}

	// collect flow activity info for determinng query vs invoke transactions
	activityData, err := collectFlowActivities(appconfig)
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
		Schemas map[string]struct {
			Value string `json:"value"`
		} `json:"schemas,omitempty"`
	}
	if err := json.Unmarshal(appconfig, &contractData); err != nil {
		return errors.Wrapf(err, "failed to extract contract data from flogo config file %s", appfile)
	}

	// cache application schemas
	if len(contractData.Schemas) > 0 {
		for k, v := range contractData.Schemas {
			appSchema["schema://"+k] = v.Value
		}
	}

	// construct output metadata
	metadata := Metadata{
		Schema:   metaSchema,
		Info:     Info{Title: contractData.Name, Version: contractData.Version},
		Contract: &Contract{Name: contractData.Name, Transactions: map[string]*Transaction{}},
	}

	// select all handlers of the first fabric transaction trigger
	for _, trig := range contractData.Triggers {
		if trig.Ref == "#transaction" {
			for _, handler := range trig.Handlers {
				fmt.Printf("transaction name %s\n", handler.Settings.Name)
				txn := extractTransactionSchema(handler.Schemas)
				txn.Name = handler.Settings.Name
				txn.Description = handler.Description
				if readOnly, ok := activityData[txn.Name]; ok && readOnly {
					txn.Operation = "query"
				} else {
					txn.Operation = "invoke"
				}
				metadata.Contract.Transactions[txn.Name] = &txn
			}
			break
		}
	}

	// set unique id for shared schema
	setSharedSchemaRef(contractData.Name)
	if len(sharedSchema) > 0 {
		metadata.Components = map[string]SchemaValue{}
		for k := range sharedSchema {
			sv := schemaCache[k]
			metadata.Components[sv.Ref] = *sv
		}
	}

	// use shared schema in transacions
	for _, txn := range metadata.Contract.Transactions {
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
			"$ref": "#/components/%s"
		}`, sv.Ref)
		return []byte(shared)
	}
	return schema
}

func extractSchemaValue(def json.RawMessage) []byte {
	defstr := string(def)
	if strings.HasPrefix(defstr, "\"schema://") {
		// replace with app-level schema
		k := defstr[1 : len(defstr)-1]
		if v, ok := appSchema[k]; ok {
			return addSchemaToCache(v)
		}
		return nil
	}

	// parse schema value
	var sch struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(def, &sch); err == nil && sch.Value != "" {
		return addSchemaToCache(sch.Value)
	}
	return nil
}

func extractTransactionSchema(schemas map[string]json.RawMessage) Transaction {
	var parameters, transient, returns []byte
	if reply, ok := schemas["reply"]; ok {
		var r struct {
			Returns json.RawMessage `json:"returns"`
		}
		if err := json.Unmarshal(reply, &r); err == nil && r.Returns != nil {
			returns = extractSchemaValue(r.Returns)
		}
	}
	if output, ok := schemas["output"]; ok {
		var o struct {
			Parameters json.RawMessage `json:"parameters,omitempty"`
			Transient  json.RawMessage `json:"transient,omitempty"`
		}
		if err := json.Unmarshal(output, &o); err == nil {
			if o.Parameters != nil {
				parameters = extractSchemaValue(o.Parameters)
			}
			if o.Transient != nil {
				transient = extractSchemaValue(o.Transient)
			}
		}
	}
	return Transaction{Parameters: parameters, Transient: transient, Returns: returns}
}

func addSchemaToCache(value string) []byte {
	var sv SchemaValue
	if err := json.Unmarshal([]byte(value), &sv); err != nil {
		fmt.Printf("failed to unmarshal schema value: %+v\n", err)
	} else {
		if svbytes, err := json.Marshal(sv); err != nil {
			fmt.Printf("failed to marshal schema value: %+v\n", err)
		} else {
			skey := string(svbytes)
			// add to cache
			if _, ok := schemaCache[skey]; ok {
				// found a reused data schema, record it as shared if it is not already known
				if _, ok := sharedSchema[skey]; !ok {
					sharedSchema[skey] = schemaSeq
					schemaSeq++
				}
			} else {
				schemaCache[skey] = &sv
			}
			return svbytes
		}
	}
	return nil
}

func setSharedSchemaRef(appname string) {
	for k, v := range sharedSchema {
		sv := schemaCache[k]
		if sv.Type == "array" {
			sv.Ref = fmt.Sprintf("%s_%d", appname, v)
			if oref, ok := schemaOverride[sv.Ref]; ok {
				sv.Ref = oref
			}
		} else {
			if n, err := objectPropertyCount(sv.Properties); err == nil && n > 2 {
				sv.Ref = fmt.Sprintf("%s_%d", appname, v)
				if oref, ok := schemaOverride[sv.Ref]; ok {
					sv.Ref = oref
				}
			} else {
				delete(sharedSchema, k)
				if err != nil {
					fmt.Printf("error when removing shared schema %s: %+v\n", k, err)
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

// determine if a transaction is read-only
func collectFlowActivities(appconfig []byte) (map[string]bool, error) {
	var flowData struct {
		Resources []struct {
			Data struct {
				Name  string `json:"name"`
				Tasks []struct {
					Name     string `json:"name"`
					Activity struct {
						Ref      string `json:"ref"`
						Settings struct {
							FlowURI string `json:"flowURI"`
						} `json:"settings"`
					} `json:"activity"`
				} `json:"tasks"`
			} `json:"data"`
		} `json:"resources"`
	}
	if err := json.Unmarshal(appconfig, &flowData); err != nil {
		return nil, errors.Wrap(err, "failed to extract activity data from flogo config")
	}

	// map flow-name => true if it calls #put, #putall or #delete
	result := make(map[string]bool)
	// map sub-flow-name => list of flows that calls the subflow
	subflowRef := make(map[string][]string)

	for _, resource := range flowData.Resources {
		name := resource.Data.Name
		readOnly := true
		for _, task := range resource.Data.Tasks {
			ref := task.Activity.Ref
			if ref == "#put" || ref == "#putall" || ref == "#delete" {
				readOnly = false
				break
			} else if ref == "#subflow" {
				// check subflow URI: "res://flow:name"
				flowURI := task.Activity.Settings.FlowURI
				colon := strings.LastIndex(flowURI, ":")
				flowName := flowURI[colon+1:]
				flowRefs, ok := subflowRef[flowName]
				if !ok {
					flowRefs = []string{}
				}
				fmt.Printf("add subflow ref %s -> %s\n", name, flowName)
				subflowRef[flowName] = append(flowRefs, name)
			}
		}
		result[name] = readOnly
	}

	// override caller result if subflow is not read-only
	for {
		updated := false
		for k, v := range subflowRef {
			if ro, ok := result[k]; ok && !ro {
				for _, n := range v {
					if r, ok := result[n]; !ok || r {
						fmt.Printf("set %s not read-only\n", n)
						updated = true
						result[n] = false
					}
				}
			}
		}
		if !updated {
			break
		}
	}
	return result, nil
}

// cache for all object types
var objectTypes map[string]*objectDef
var transTypes []*transactionDef
var componentNames map[string]string
var objectSeq int

// generate contract graphql file from metadata
func generateGqlfile(metafile, gqlfile string) error {
	metabytes, err := ioutil.ReadFile(metafile)
	if err != nil {
		return err
	}
	var meta Metadata
	if err := json.Unmarshal(metabytes, &meta); err != nil {
		return err
	}

	// extract shared object types
	objectTypes = make(map[string]*objectDef)
	componentNames = make(map[string]string)
	for k, v := range meta.Components {
		if v.Type == "object" {
			var props map[string]interface{}
			if err := json.Unmarshal(v.Properties, &props); err != nil {
				fmt.Printf("failed to unmarshal component schema %s: %+v\n", k, err)
				continue
			}
			cdef := getDefForObject(k, props)
			componentNames[k] = cdef.typeID
		} else if v.Type == "array" {
			var items map[string]interface{}
			if err := json.Unmarshal(v.Items, &items); err != nil {
				fmt.Printf("failed to unmarshal component array %s: %+v\n", k, err)
				continue
			}
			cdef := getDefForArray(k, items)
			componentNames[k] = cdef.typeID
		} else {
			fmt.Println("component should not be type of", v.Type)
		}
	}

	// collect transaction as graphQL operations
	transTypes = []*transactionDef{}
	for k, v := range meta.Contract.Transactions {
		transDef := transactionDef{
			name:       k,
			operation:  v.Operation,
			parameters: []*attributeDef{},
		}
		if v.Parameters != nil {
			if odef, err := extractObjectDef(k, v.Parameters); err == nil {
				transDef.parameters = append(transDef.parameters, odef.attributes...)
			}
		}
		if v.Transient != nil {
			if odef, err := extractObjectDef(k+"_t", v.Transient); err == nil {
				transDef.parameters = append(transDef.parameters, odef.attributes...)
			}
		}
		if v.Returns != nil {
			if odef, err := extractObjectDef(k+"Return", v.Returns); err == nil {
				transDef.returns = odef
			}
		}
		transTypes = append(transTypes, &transDef)
	}

	// write graphQL file
	writeGQL(gqlfile)
	return nil
}

func extractObjectDef(name string, data []byte) (*objectDef, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		fmt.Printf("failed to unmarshal object def %s: %+v\n", name, err)
		return nil, err
	}
	odef := getObjectDef(name, obj)
	return odef, nil
}

func writeGQL(gqlfile string) error {
	fmt.Println("write metadata to", gqlfile)
	p := Subst(gqlfile)
	d := filepath.Dir(p)
	if err := os.MkdirAll(d, 0755); err != nil {
		return err
	}
	os.Remove(gqlfile)
	f, err := os.OpenFile(gqlfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	printGQLSchemaDef(f)
	printGQLOperations(f, "query")
	printGQLOperations(f, "invoke")

	// collect object types used by transactions
	otypes := make(map[string]bool)
	for _, trans := range transTypes {
		for _, p := range trans.parameters {
			otypes[p.typeID] = true
		}
		r := strings.Replace(trans.returns.typeID, "[", "", -1)
		r = strings.Replace(r, "]", "", -1)
		otypes[r] = false
	}
	// collect attribute types of the collected objects
	for collectReferencedTypes(otypes) {
		fmt.Printf("collected %d types\n", len(otypes))
	}
	// print types referenced by transactions
	for _, v := range objectTypes {
		if input, ok := otypes[v.typeID]; ok {
			f.WriteString(objectGQL(v, input))
			f.WriteString("\n")
		}
	}
	f.Sync()
	return nil
}

func collectReferencedTypes(types map[string]bool) bool {
	size := len(types)
	for t, input := range types {
		if def, ok := objectTypes[t]; ok {
			for _, p := range def.attributes {
				r := strings.Replace(p.typeID, "[", "", -1)
				r = strings.Replace(r, "]", "", -1)
				types[r] = input
			}
		}
	}
	return len(types) > size
}

func printGQLSchemaDef(file *os.File) {
	file.WriteString(`schema {
	query: Query
	mutation: Mutation
}
`)
}

func printGQLOperations(file *os.File, op string) {
	if op == "query" {
		file.WriteString("type Query {\n")
	} else {
		file.WriteString("type Mutation {\n")
	}
	for _, tx := range transTypes {
		if tx.operation == op {
			file.WriteString("\t")
			file.WriteString(transactionGQL(tx))
			file.WriteString("\n")
		}
	}
	file.WriteString("}\n")
}

func transactionGQL(def *transactionDef) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s(", def.name)
	for i, p := range def.parameters {
		if i > 0 {
			fmt.Fprint(&b, ", ")
		}
		if p.isArray {
			fmt.Fprintf(&b, "%s: [%s]", p.name, p.typeID)
		} else {
			fmt.Fprintf(&b, "%s: %s", p.name, p.typeID)
		}
	}
	fmt.Fprintf(&b, "): %s", def.returns.typeID)
	return b.String()
}

func objectGQL(def *objectDef, input bool) string {
	var b strings.Builder
	if input {
		fmt.Fprintf(&b, "input %s {\n", def.typeID)
	} else {
		fmt.Fprintf(&b, "type %s {\n", def.typeID)
	}
	for _, attr := range def.attributes {
		if attr.isArray {
			fmt.Fprintf(&b, "\t%s: [%s]\n", attr.name, attr.typeID)
		} else {
			fmt.Fprintf(&b, "\t%s: %s\n", attr.name, attr.typeID)
		}
	}
	fmt.Fprintln(&b, "}")
	return b.String()
}

type transactionDef struct {
	name       string
	operation  string
	parameters []*attributeDef
	returns    *objectDef
}

type attributeDef struct {
	name    string
	typeID  string
	isArray bool
}

type objectDef struct {
	typeID     string
	attributes []*attributeDef
}

func refTypeID(ref string) string {
	// ref should be in format '#/components/typeID'
	slash := strings.LastIndex(ref, "/")
	if slash >= 0 {
		t := ref[slash+1:]
		if ot, ok := componentNames[t]; ok {
			t = ot
		}
		return t
	}
	return ref
}

func getAttributeDefs(props map[string]interface{}) []*attributeDef {
	attrs := []*attributeDef{}
	for k, v := range props {
		def := v.(map[string]interface{})
		isArray := false
		if ref, ok := def["$ref"]; ok {
			attrs = append(attrs, &attributeDef{
				name:    k,
				typeID:  refTypeID(ref.(string)),
				isArray: isArray,
			})
			continue
		}
		t := def["type"]
		if t.(string) == "array" {
			def = def["items"].(map[string]interface{})
			isArray = true
		}
		odef := getObjectDef(k, def)
		if odef != nil {
			attrs = append(attrs, &attributeDef{
				name:    k,
				typeID:  odef.typeID,
				isArray: isArray,
			})
		}
	}
	return attrs
}

func getDefForObject(name string, props map[string]interface{}) *objectDef {
	objectSeq++

	// check cached object defs
	attrs := getAttributeDefs(props)
	for _, cached := range objectTypes {
		if isDuplicateDef(cached, attrs) {
			return cached
		}
	}

	// check override name
	t := name
	if ot, ok := schemaOverride[name]; ok {
		t = ot
	}
	if _, ok := objectTypes[t]; ok {
		// name already used by a different object, so use a unique name
		t = fmt.Sprintf("%s%d", name, objectSeq)
		if ot, ok := schemaOverride[t]; ok {
			if _, ok := objectTypes[ot]; !ok {
				t = ot
			} else {
				fmt.Println("Ignore duplicate override to type", ot)
			}
		}
	}

	// add new object to cache
	odef := objectDef{typeID: t, attributes: attrs}
	objectTypes[odef.typeID] = &odef
	return &odef
}

func getDefForArray(name string, items map[string]interface{}) *objectDef {
	item := getObjectDef(name, items)
	return &objectDef{typeID: "[" + item.typeID + "]"}
}

func getObjectDef(name string, def map[string]interface{}) *objectDef {
	if ref, ok := def["$ref"]; ok {
		return &objectDef{typeID: refTypeID(ref.(string))}
	}
	t := def["type"]
	switch t.(string) {
	case "string":
		return &objectDef{typeID: "String"}
	case "integer":
		return &objectDef{typeID: "Int"}
	case "number":
		return &objectDef{typeID: "Float"}
	case "boolean":
		return &objectDef{typeID: "Boolean"}
	case "array":
		return getDefForArray(name, def["items"].(map[string]interface{}))
	case "object":
		nm := strings.Title(name + "Type")
		return getDefForObject(nm, def["properties"].(map[string]interface{}))
	default:
		fmt.Println("Unknown property type:", t)
	}
	return nil
}

func isDuplicateDef(def *objectDef, attrs []*attributeDef) bool {
	if len(def.attributes) != len(attrs) {
		return false
	}
	names := make(map[string]string)
	for _, v := range def.attributes {
		names[v.name] = v.typeID
	}
	for _, v := range attrs {
		if t, ok := names[v.name]; !ok || t != v.typeID {
			return false
		}
	}
	return true
}
