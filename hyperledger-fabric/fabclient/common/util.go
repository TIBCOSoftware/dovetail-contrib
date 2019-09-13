package common

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/pkg/errors"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/schema"
	"github.com/project-flogo/core/support/log"
	jschema "github.com/xeipuuv/gojsonschema"
)

// Create a new logger
var logger = log.ChildLogger(log.RootLogger(), "client-common")

// GetActivityInputSchema returns schema of an activity input attribute
func GetActivityInputSchema(ctx activity.Context, name string) (string, error) {
	if sIO, ok := ctx.(schema.HasSchemaIO); ok {
		s := sIO.GetInputSchema(name)
		if s != nil {
			logger.Debugf("schema for attribute '%s': %T, %s\n", name, s, s.Value())
			return s.Value(), nil
		}
	}
	return "", errors.Errorf("schema not found for attribute %s", name)
}

// ParameterIndex stores transaction parameters and its location in raw JSON schema string
// start and end location is used to sort the parameter list to match the parameter order in schema
type ParameterIndex struct {
	Name     string
	JSONType string
	start    int
	end      int
}

// OrderedParameters returns parameters of a JSON schema object sorted by their position in schema definition
// This is necessary because Golang JSON parser does not maintain the sequence of object parameters.
func OrderedParameters(schemaData []byte) ([]ParameterIndex, error) {
	if schemaData == nil || len(schemaData) == 0 {
		logger.Debug("schema data is empty")
		return nil, nil
	}
	// extract root object properties from JSON schema
	var rawProperties struct {
		Data json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(schemaData, &rawProperties); err != nil {
		logger.Errorf("failed to extract properties from metadata: %+v", err)
		return nil, err
	}

	// extract parameter names from raw object properties
	var params map[string]json.RawMessage
	if err := json.Unmarshal(rawProperties.Data, &params); err != nil {
		logger.Errorf("failed to extract parameters from object schema: %+v", err)
		return nil, err
	}

	// collect parameter locations in the raw object schema
	var paramIndex []ParameterIndex
	for p, v := range params {
		// encode parameter name with quotes
		key, _ := json.Marshal(p)
		// key may exist in raw schema multiple times,
		// so check each occurence to determine its correct location in the raw schema
		items := bytes.Split(rawProperties.Data, key)
		pos := 0
		for _, seg := range items {
			if pos == 0 {
				// first segment should not be the key definition
				pos += len(seg)
				continue
			}
			vpos := bytes.Index(seg, v)
			if vpos >= 0 {
				// the segment contains the key definition, so collect its position in raw schema
				endPos := pos + len(key) + vpos + len(v)
				// extract JSON type of the parameter
				var paramDef struct {
					RawType string `json:"type"`
				}
				if err := json.Unmarshal(v, &paramDef); err != nil {
					logger.Errorf("failed to extract JSON type of parameter %s: %+v", p, err)
				}
				paramType := jschema.TYPE_OBJECT
				if paramDef.RawType != "" {
					paramType = paramDef.RawType
				}
				logger.Debugf("add index parameter '%s' type '%s'\n", p, paramType)
				paramIndex = addIndex(paramIndex, ParameterIndex{Name: p, JSONType: paramType, start: pos, end: endPos})
			}
			pos += len(key) + len(seg)
		}
	}

	// sort parameter index by start location in raw schema
	if len(paramIndex) > 1 {
		sort.Slice(paramIndex, func(i, j int) bool {
			return paramIndex[i].start < paramIndex[j].start
		})
	}
	return paramIndex, nil
}

// addIndex adds a new parameter position to the index, ignore or merge index if index region overlaps.
func addIndex(parameters []ParameterIndex, param ParameterIndex) []ParameterIndex {
	for i, v := range parameters {
		if param.start > v.start && param.start < v.end {
			// ignore if new param's start postion falls in region covered by a known parameter
			return parameters
		} else if v.start > param.start && v.start < param.end {
			// replace old parameter region if its start position falls in the region covered by the new parameter
			updated := append(parameters[:i], param)
			if len(parameters) > i+1 {
				// check the remaining knonw parameters
				for _, p := range parameters[i+1:] {
					if !(p.start > param.start && p.start < param.end) {
						updated = append(updated, p)
					}
				}
			}
			return updated
		}
	}
	// append new parameter
	return append(parameters, param)
}
