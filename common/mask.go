package common

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/pingcap/errors"

	"github.com/romberli/go-util/constant"
)

const (
	// sensitive keyword
	DefaultSensitivePassKeyword   = "pass"
	DefaultSensitiveSecretKeyword = "secret"
	DefaultSensitivePwdKeyword    = "pwd"
	// sensitive pattern
	DefaultSensitivePatternIdentifiedBy = `(?i)(IDENTIFIED BY\s*')([^']+)(')`
	DefaultReplacementString            = "${1}******${3}"
)

var (
	DefaultSensitiveKeywords = []string{
		DefaultSensitivePassKeyword,
		DefaultSensitiveSecretKeyword,
		DefaultSensitivePwdKeyword,
	}
	DefaultIdentifiedByPattern, _ = NewSensitivePattern(DefaultSensitivePatternIdentifiedBy, DefaultReplacementString)
	DefaultSensitivePatterns      = []*SensitivePattern{
		DefaultIdentifiedByPattern,
	}
)

// SensitivePattern 定义敏感值模式
type SensitivePattern struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// NewSensitivePattern 创建敏感模式
func NewSensitivePattern(pattern string, replacement string) (*SensitivePattern, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if replacement == constant.EmptyString {
		replacement = "${1}xxxxxx${3}"
	}

	return &SensitivePattern{
		Pattern:     re,
		Replacement: replacement,
	}, nil
}

// MaskJSON masks the sensitive fields in the json body
func MaskJSON(jsonBytes []byte, sensitiveFields []string, excludes ...string) ([]byte, error) {
	return MaskJSONWithPatterns(jsonBytes, sensitiveFields, nil, excludes...)
}

// MaskJSONWithPatterns masks the sensitive fields in the json body with the sensitive patterns
// Note: the sensitive patterns will be applied to the value not the key
func MaskJSONWithPatterns(jsonBytes []byte, sensitiveFields []string, sensitivePatterns []*SensitivePattern, excludes ...string) ([]byte, error) {
	if len(jsonBytes) == constant.ZeroInt {
		return jsonBytes, nil
	}

	var data interface{}
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return jsonBytes, errors.Trace(err)
	}

	maskedValue := maskValue(data, sensitiveFields, sensitivePatterns, excludes...)

	result, err := json.Marshal(maskedValue)
	if err != nil {
		return jsonBytes, errors.Trace(err)
	}

	return result, nil
}

func maskValue(value interface{}, sensitiveFields []string, sensitivePatterns []*SensitivePattern, excludes ...string) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if isSensitiveField(key, sensitiveFields, excludes...) {
				v[key] = constant.DefaultMaskedValue
			} else {
				// mask value recursively
				processedVal := maskValue(val, sensitiveFields, sensitivePatterns, excludes...)
				v[key] = maskSensitiveValue(processedVal, sensitivePatterns)
			}
		}
		return v

	case []interface{}:
		for i, item := range v {
			processedItem := maskValue(item, sensitiveFields, sensitivePatterns, excludes...)
			v[i] = maskSensitiveValue(processedItem, sensitivePatterns)
		}
		return v

	default:
		return maskSensitiveValue(v, sensitivePatterns)
	}
}

// isSensitiveField checks if the field name contains any of the sensitive fields
func isSensitiveField(fieldName string, sensitiveFields []string, excludes ...string) bool {
	lowerField := strings.ToLower(fieldName)
	for _, exclude := range excludes {
		if strings.Contains(lowerField, strings.ToLower(exclude)) {
			return false
		}
	}
	for _, sensitiveField := range sensitiveFields {
		if strings.Contains(lowerField, strings.ToLower(sensitiveField)) {
			return true
		}
	}

	return false
}

func maskSensitiveValue(value interface{}, sensitivePatterns []*SensitivePattern) interface{} {
	str, ok := value.(string)
	if !ok {
		return value
	}

	result := str
	for _, pattern := range sensitivePatterns {
		if pattern == nil || pattern.Pattern == nil {
			continue
		}
		result = pattern.Pattern.ReplaceAllString(result, pattern.Replacement)
	}

	return result
}
