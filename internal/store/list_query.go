package store

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func ParseListFilterExpression(expression string) (ListFilter, error) {
	var filter ListFilter
	if strings.TrimSpace(expression) == "" {
		return filter, nil
	}
	clauses := splitListFilterClauses(expression)
	if len(clauses) == 0 {
		return filter, errors.New("filter must not be empty")
	}
	for _, clause := range clauses {
		field, operator, rawValues, err := splitListFilterClause(clause)
		if err != nil {
			return ListFilter{}, err
		}
		values, err := parseListFilterValues(rawValues)
		if err != nil {
			return ListFilter{}, err
		}
		switch field {
		case "displayName":
			if !listFilterOperatorAllowed(operator, "=", "!=", "contains") {
				return ListFilter{}, unsupportedListFilterOperator(field, "=", "!=", "contains")
			}
			filter.Clauses = append(filter.Clauses, listFilterClause(field, "", operator, values))
			if operator == "=" {
				filter.DisplayName = append(filter.DisplayName, values...)
			}
		case "type":
			if !listFilterOperatorAllowed(operator, "=", "!=", "contains") {
				return ListFilter{}, unsupportedListFilterOperator(field, "=", "!=", "contains")
			}
			filter.Clauses = append(filter.Clauses, listFilterClause(field, "", operator, values))
			if operator == "=" {
				filter.Types = append(filter.Types, values...)
			}
		case "publisherId":
			if !listFilterOperatorAllowed(operator, "=", "!=", "contains") {
				return ListFilter{}, unsupportedListFilterOperator(field, "=", "!=", "contains")
			}
			filter.Clauses = append(filter.Clauses, listFilterClause(field, "", operator, values))
			if operator == "=" {
				filter.PublisherIDs = append(filter.PublisherIDs, values...)
			}
		case "tags":
			if !listFilterOperatorAllowed(operator, "=", "!=", "contains") {
				return ListFilter{}, unsupportedListFilterOperator(field, "=", "!=", "contains")
			}
			filter.Clauses = append(filter.Clauses, listFilterClause(field, "", operator, values))
			if operator == "=" {
				filter.Tags = append(filter.Tags, values...)
			}
		case "capabilities":
			if !listFilterOperatorAllowed(operator, "=", "!=", "contains") {
				return ListFilter{}, unsupportedListFilterOperator(field, "=", "!=", "contains")
			}
			filter.Clauses = append(filter.Clauses, listFilterClause(field, "", operator, values))
			if operator == "=" {
				filter.Capabilities = append(filter.Capabilities, values...)
			}
		case "createdAfter":
			if !listFilterOperatorAllowed(operator, ">", ">=") {
				return ListFilter{}, unsupportedListFilterOperator(field, ">", ">=")
			}
			timestamp, err := singleListFilterTime(field, values)
			if err != nil {
				return ListFilter{}, err
			}
			filter.Clauses = append(filter.Clauses, listFilterTimeClause(field, operator, timestamp))
			if operator == ">" {
				filter.CreatedAfter = &timestamp
			}
		case "updatedAfter":
			if !listFilterOperatorAllowed(operator, ">", ">=") {
				return ListFilter{}, unsupportedListFilterOperator(field, ">", ">=")
			}
			timestamp, err := singleListFilterTime(field, values)
			if err != nil {
				return ListFilter{}, err
			}
			filter.Clauses = append(filter.Clauses, listFilterTimeClause(field, operator, timestamp))
			if operator == ">" {
				filter.UpdatedAfter = &timestamp
			}
		default:
			if strings.HasPrefix(field, "metadata.") {
				if !listFilterOperatorAllowed(operator, "=", "!=", "contains") {
					return ListFilter{}, unsupportedListFilterOperator(field, "=", "!=", "contains")
				}
				key, err := listMetadataKey(field)
				if err != nil {
					return ListFilter{}, err
				}
				filter.Clauses = append(filter.Clauses, listFilterClause("metadata", key, operator, values))
				if filter.Metadata == nil {
					filter.Metadata = map[string][]string{}
				}
				if operator == "=" {
					filter.Metadata[key] = append(filter.Metadata[key], values...)
				}
				continue
			}
			return ListFilter{}, fmt.Errorf("unsupported filter field %q", field)
		}
	}
	return filter, nil
}

func ParseListOrderBy(raw string) (ListOrder, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ListOrder{}, nil
	}
	parts := strings.Fields(raw)
	if len(parts) < 1 || len(parts) > 2 {
		return ListOrder{}, errors.New("orderBy must be a field optionally followed by ASC or DESC")
	}
	field, err := normalizeListOrderField(parts[0])
	if err != nil {
		return ListOrder{}, err
	}
	direction := "ASC"
	if len(parts) == 2 {
		direction = strings.ToUpper(parts[1])
		if direction != "ASC" && direction != "DESC" {
			return ListOrder{}, errors.New("orderBy direction must be ASC or DESC")
		}
	}
	return ListOrder{Field: field, Direction: direction}, nil
}

func splitListFilterClauses(expression string) []string {
	clauses := []string{}
	start := 0
	quoted := rune(0)
	for index, char := range expression {
		if quoted != 0 {
			if char == quoted {
				quoted = 0
			}
			continue
		}
		if char == '\'' || char == '"' {
			quoted = char
			continue
		}
		if hasListFilterANDAt(expression, index) {
			clauses = append(clauses, strings.TrimSpace(expression[start:index]))
			start = index + len(" AND ")
		}
	}
	clauses = append(clauses, strings.TrimSpace(expression[start:]))
	return clauses
}

func hasListFilterANDAt(expression string, index int) bool {
	if index+len(" AND ") > len(expression) {
		return false
	}
	return strings.EqualFold(expression[index:index+len(" AND ")], " AND ")
}

func splitListFilterClause(clause string) (string, string, string, error) {
	for _, operator := range []string{">=", "!=", ">", "="} {
		if index := indexOutsideQuotes(clause, operator); index >= 0 {
			field := strings.TrimSpace(clause[:index])
			value := strings.TrimSpace(clause[index+len(operator):])
			if field == "" || value == "" {
				return "", "", "", fmt.Errorf("invalid filter clause %q", clause)
			}
			return field, operator, value, nil
		}
	}
	if field, value, ok := splitListFilterKeywordOperator(clause, "contains"); ok {
		return field, "contains", value, nil
	}
	return "", "", "", fmt.Errorf("invalid filter clause %q", clause)
}

func splitListFilterKeywordOperator(clause string, operator string) (string, string, bool) {
	quoted := rune(0)
	operatorLength := len(operator)
	for index, char := range clause {
		if quoted != 0 {
			if char == quoted {
				quoted = 0
			}
			continue
		}
		if char == '\'' || char == '"' {
			quoted = char
			continue
		}
		if index+operatorLength > len(clause) || !strings.EqualFold(clause[index:index+operatorLength], operator) {
			continue
		}
		beforeOK := index == 0 || isListFilterWhitespace(rune(clause[index-1]))
		afterIndex := index + operatorLength
		afterOK := afterIndex == len(clause) || isListFilterWhitespace(rune(clause[afterIndex]))
		if !beforeOK || !afterOK {
			continue
		}
		field := strings.TrimSpace(clause[:index])
		value := strings.TrimSpace(clause[afterIndex:])
		return field, value, field != "" && value != ""
	}
	return "", "", false
}

func isListFilterWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func indexOutsideQuotes(value string, needle string) int {
	quoted := rune(0)
	for index, char := range value {
		if quoted != 0 {
			if char == quoted {
				quoted = 0
			}
			continue
		}
		if char == '\'' || char == '"' {
			quoted = char
			continue
		}
		if strings.HasPrefix(value[index:], needle) {
			return index
		}
	}
	return -1
}

func parseListFilterValues(raw string) ([]string, error) {
	parts := splitCommaSeparatedValues(raw)
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if len(value) >= 2 {
			first := value[0]
			last := value[len(value)-1]
			if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
				value = value[1 : len(value)-1]
			}
		}
		value = strings.TrimSpace(value)
		if value == "" {
			return nil, errors.New("filter values must not be empty")
		}
		values = append(values, value)
	}
	if len(values) == 0 {
		return nil, errors.New("filter values must not be empty")
	}
	return values, nil
}

func splitCommaSeparatedValues(raw string) []string {
	values := []string{}
	start := 0
	quoted := rune(0)
	for index, char := range raw {
		if quoted != 0 {
			if char == quoted {
				quoted = 0
			}
			continue
		}
		if char == '\'' || char == '"' {
			quoted = char
			continue
		}
		if char == ',' {
			values = append(values, raw[start:index])
			start = index + 1
		}
	}
	values = append(values, raw[start:])
	return values
}

func singleListFilterTime(field string, values []string) (time.Time, error) {
	if len(values) != 1 {
		return time.Time{}, fmt.Errorf("filter field %q requires exactly one timestamp", field)
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"} {
		if parsed, err := time.Parse(layout, values[0]); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("filter field %q requires an ISO 8601 timestamp", field)
}

func listMetadataKey(field string) (string, error) {
	key := strings.TrimPrefix(field, "metadata.")
	if key == "" {
		return "", errors.New("metadata filter key must not be empty")
	}
	for _, char := range key {
		if char >= 'a' && char <= 'z' {
			continue
		}
		if char >= 'A' && char <= 'Z' {
			continue
		}
		if char >= '0' && char <= '9' {
			continue
		}
		switch char {
		case '_', '-', '.':
			continue
		default:
			return "", fmt.Errorf("metadata filter key %q contains unsupported character %q", key, char)
		}
	}
	return key, nil
}

func listFilterClause(field string, metadataKey string, operator string, values []string) ListFilterClause {
	return ListFilterClause{
		Field:       field,
		MetadataKey: metadataKey,
		Operator:    operator,
		Values:      append([]string(nil), values...),
	}
}

func listFilterTimeClause(field string, operator string, timestamp time.Time) ListFilterClause {
	return ListFilterClause{
		Field:    field,
		Operator: operator,
		Time:     &timestamp,
	}
}

func listFilterOperatorAllowed(operator string, allowed ...string) bool {
	for _, candidate := range allowed {
		if operator == candidate {
			return true
		}
	}
	return false
}

func unsupportedListFilterOperator(field string, allowed ...string) error {
	return fmt.Errorf("filter field %q only supports %s", field, strings.Join(allowed, ", "))
}

func normalizeListOrderField(field string) (string, error) {
	switch field {
	case "displayName", "display_name", "name":
		return "displayName", nil
	case "type":
		return "type", nil
	case "createdAt", "created_at":
		return "createdAt", nil
	case "updatedAt", "updated_at":
		return "updatedAt", nil
	case "publisherId", "publisher_id":
		return "publisherId", nil
	default:
		return "", fmt.Errorf("unsupported orderBy field %q", field)
	}
}
