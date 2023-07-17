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

type ListMapUnion struct {
	lst *TransformList
	mp *TransformMap
	
	is_list_flag bool
}

func (u *ListMapUnion) is_list() bool {
	return u.is_list_flag
}


func (u *ListMapUnion) is_dict() bool {
	return !u.is_list()
}

func (u *ListMapUnion) use_dict() *TransformMap {
	if !u.is_dict() {
		panic(`data is not a dict`)
	}
	if u.mp == nil {
		u.mp = &TransformMap{}
	}
	return u.mp
}

func (u *ListMapUnion) use_list() *TransformList {
	if !u.is_list() {
		panic(`data is not a list`)
	}
	if u.lst == nil {
		u.lst = &TransformList{}
	}
	return u.lst
}


func (u *ListMapUnion) extract() any {
	if !u.is_list() {
		return u.use_dict()
	}
	return u.use_list()
}

func (u *ListMapUnion) pack(value any, key string) {
	if !u.is_list() {
		(*u.use_dict())[key] = value
		return
	}
	x := append( *u.use_list(), value)
	u.lst = &x
}
