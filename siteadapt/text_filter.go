package siteadapt

import (
	"fmt"
)

type textFilter interface {
	// 过滤文本，过滤不成功请返回原始文本，交给后续 Filter 处理
	doFilter() (string, error)
}

type fieldFilterType string

const (
	fieldFilterTypeResearch    fieldFilterType = "re_search"
	fieldFilterTypeSplit       fieldFilterType = "split"
	fieldFilterTypeReplace     fieldFilterType = "replace"
	fieldFilterTypeStrip       fieldFilterType = "strip"
	fieldFilterTypeAppendLeft  fieldFilterType = "append_left"
	fieldFilterTypeQueryString fieldFilterType = "querystring"
	fieldFilterTypeRegex       fieldFilterType = "regex"
	fieldFilterTypeByteSize    fieldFilterType = "byte_size"
	fieldFilterTypeTimestamp   fieldFilterType = "timestamp"
	fieldFilterTypeEq          fieldFilterType = "eq"
	fieldFilterTypeCase        fieldFilterType = "case"
	fieldFilterTypeNotBlank    fieldFilterType = "not_blank"
	fieldFilterTypeBlank       fieldFilterType = "blank"
	fieldFilterTypeConstant    fieldFilterType = "constant"
)

var textFilterRegistry = map[string]func(string, *Filter) (string, error){
	string(fieldFilterTypeResearch): func(text string, filter *Filter) (string, error) {
		return researchFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeSplit): func(text string, filter *Filter) (string, error) {
		return splitFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeReplace): func(text string, filter *Filter) (string, error) {
		return replaceFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeStrip): func(text string, filter *Filter) (string, error) {
		return stripFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeAppendLeft): func(text string, filter *Filter) (string, error) {
		return appendLeftFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeQueryString): func(text string, filter *Filter) (string, error) {
		return queryStringFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeRegex): func(text string, filter *Filter) (string, error) {
		return regexFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeByteSize): func(text string, filter *Filter) (string, error) {
		return byteSizeFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeTimestamp): func(text string, filter *Filter) (string, error) {
		return timestampFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeEq): func(text string, filter *Filter) (string, error) {
		return eqFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeCase): func(text string, filter *Filter) (string, error) {
		return caseFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeNotBlank): func(text string, filter *Filter) (string, error) {
		return notBlankFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeBlank): func(text string, filter *Filter) (string, error) {
		return blankFilter{text, filter}.doFilter()
	},
	string(fieldFilterTypeConstant): func(text string, filter *Filter) (string, error) {
		return constantFilter{text, filter}.doFilter()
	},
}

func filter(text string, filter *Filter) (string, error) {
	if f, ok := textFilterRegistry[filter.Name]; ok {
		return f(text, filter)
	} else {
		return "", fmt.Errorf("unknown text filter: %s", filter.Name)
	}
}
