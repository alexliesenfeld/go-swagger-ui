package go_swagger_ui

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

func yamlOrJSONToJSON(data []byte) ([]byte, error) {
	vMap := make(map[string]any)
	if err := yaml.Unmarshal(data, &vMap); err != nil {
		if err = json.Unmarshal(data, &vMap); err != nil {
			return nil, fmt.Errorf("cannot unmarshal value as YAML or JSON: %w", err)
		}
	}

	newV, err := json.Marshal(vMap)
	if err != nil {
		return nil, fmt.Errorf("cannot convert value to JSON: %w", err)
	}

	return newV, nil
}
