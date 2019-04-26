package common

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// GetSettings returns map of parameter values of a connector
func GetSettings(connection interface{}) (map[string]interface{}, error) {
	connectionObject, err := data.CoerceToObject(connection)
	if err != nil {
		return nil, err
	}
	if connectionObject == nil {
		return nil, errors.New("Connection object is nil")
	}

	settings, ok := connectionObject["settings"].([]interface{})
	configs := make(map[string]interface{})
	if ok {
		attrs := make([]*data.Attribute, len(settings))

		attrsBinary, err := json.Marshal(settings)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(attrsBinary, &attrs)
		if err != nil {
			return nil, err
		}

		for _, v := range attrs {
			val, _ := data.CoerceToValue(v.Value(), v.Type())
			configs[v.Name()] = val
		}
	}
	return configs, nil
}

// ExtractFileContent returns content of a fileselector field
func ExtractFileContent(fileSelectorValue interface{}) ([]byte, error) {
	if fileSelectorValue == nil {
		return nil, nil
	}
	fileValue, err := data.CoerceToObject(fileSelectorValue)
	if err != nil {
		return nil, err
	}
	content, ok := fileValue["content"].(string)
	if !ok {
		return nil, errors.New("Failed extracting file content")
	}
	index := strings.Index(content, "base64,")
	if index > -1 {
		// remove prefix `data:application/octet-stream;base64,`
		content = content[index+7:]
	}
	return base64.StdEncoding.DecodeString(content)
}
