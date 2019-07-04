package common

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
)

// GetSettings returns map of parameter values of a connector
func GetSettings(connector map[string]interface{}) (map[string]interface{}, error) {
	if connector == nil {
		return nil, errors.New("Connection object is nil")
	}

	settings, ok := connector["settings"].([]interface{})
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
			val, _ := coerce.ToType(v.Value(), v.Type())
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
	fileValue, err := coerce.ToObject(fileSelectorValue)
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
