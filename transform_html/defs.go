package transform_html

import (
	"container/list"
	"regexp"
	"strconv"
)

type ParserTransfromRule struct {
	Selector               string
	Mapping                string
	Attribute_Name         string
	Regex_Sub_Value        [2]string
	Children               []*ParserTransfromRule
	Grouping               string
	Exception_on_not_found bool
}

var py_sub_rx *regexp.Regexp = regexp.MustCompile(`\\\d+`)
var py_sub_rx_plcmt = func(s string) string {
	offset := 0
	x, _ := strconv.Atoi(s[1:])
	x = x + offset
	return "$" + strconv.Itoa(x)
}

type TransformMap = map[string]any
type TransformList = []any

type TransformData interface {
	list.List | map[string]interface{}
}

type IListMapUnion interface {
	is_list() bool
	is_dict() bool
	pack(value any, key string)
	extract() any
}

type ListWrapper struct {
	listw *TransformList
}

type MapWrapper struct {
	mapw *TransformMap
}

func (u ListWrapper) is_list() bool {
	return true
}


func (u ListWrapper) is_dict() bool {
	return !u.is_list()
}


func (u MapWrapper) is_list() bool {
	return false
}

func (u ListWrapper) get() *TransformList {
	if (u.listw == nil) {
		u.listw = &TransformList{}
	}
	return u.listw
}


func (u MapWrapper) get() *TransformMap {
	if (u.mapw == nil) {
		u.mapw = &TransformMap{}
	}
	return u.mapw
}

func (u MapWrapper) is_dict() bool {
	return !u.is_list()
}


func (u ListWrapper) extract() any {
	return u.get()
}

func (u MapWrapper) extract() any {
	return u.get()
}

func (u ListWrapper) pack(value any, key string) {
	lst := append( *u.get(), value)
	u.listw = &lst
}

func (u MapWrapper) pack(value any, key string) {
	(*u.get())[key] = value
}
