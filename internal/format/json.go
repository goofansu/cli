package format

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
)

func Output(data any, fields string, jqExpr string) error {
	var outputData any = data

	if fields != "" {
		filtered, err := filterFields(data, fields)
		if err != nil {
			return err
		}
		outputData = filtered
	}

	if jqExpr != "" {
		jsonData, err := json.Marshal(outputData)
		if err != nil {
			return err
		}
		var cleanData any
		if err := json.Unmarshal(jsonData, &cleanData); err != nil {
			return err
		}

		results, err := applyJQ(cleanData, jqExpr)
		if err != nil {
			return err
		}

		for _, result := range results {
			if result == nil {
				fmt.Println()
				continue
			}
			switch v := result.(type) {
			case string:
				fmt.Println(v)
			case float64, int, int64, bool:
				fmt.Println(v)
			default:
				output, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(output))
			}
		}
		return nil
	}

	output, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func filterFields(data any, fields string) (map[string]any, error) {
	if fields == "" {
		return nil, fmt.Errorf("no fields specified")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var raw any
	if err := json.Unmarshal(jsonData, &raw); err != nil {
		return nil, err
	}

	v, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected map with 'items' field")
	}

	items, ok := v["items"].([]any)
	if !ok {
		return nil, fmt.Errorf("expected 'items' to be an array")
	}

	fieldList := strings.Split(fields, ",")
	filteredItems := make([]map[string]any, len(items))
	for i, item := range items {
		if m, ok := item.(map[string]any); ok {
			filteredItems[i] = filterMapFields(m, fieldList)
		}
	}

	return map[string]any{"total": v["total"], "items": filteredItems}, nil
}

func filterMapFields(m map[string]any, fields []string) map[string]any {
	result := make(map[string]any)
	for _, field := range fields {
		if val, exists := m[field]; exists {
			result[field] = val
		}
	}
	return result
}

func applyJQ(data any, jqExpr string) ([]any, error) {
	query, err := gojq.Parse(jqExpr)
	if err != nil {
		return nil, fmt.Errorf("jq parse error: %w", err)
	}

	iter := query.Run(data)
	var results []any
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, fmt.Errorf("jq error: %s", err.Error())
		}
		results = append(results, v)
	}

	return results, nil
}
