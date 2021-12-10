// Copyright (c) 2021, Geert JM Vanderkelen

package xdoc

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"
)

type Functions map[string]*Function

func (m Functions) SortedKeys() []string {
	keys := make([]string, len(m))
	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

type Structs map[string]*Struct

func (m Structs) SortedKeys() []string {
	keys := make([]string, len(m))
	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

type Struct struct {
	Name         string
	Constructors Functions
	Methods      Functions
	Fields       []string
	Doc          string
}

func newStruct(name string) *Struct {
	return &Struct{
		Name:         name,
		Constructors: make(Functions),
		Methods:      make(Functions),
	}
}

type Function struct {
	Name string
	Doc  string

	type_ *ast.FuncType
	recv  *ast.FieldList
}

func (sf Function) Signature() string {
	params := argumentList(sf.type_.Params)
	results := argumentList(sf.type_.Results)

	var result string
	switch len(results) {
	case 0:
		result = ""
	case 1:
		result = " " + results[0]
	default:
		result = " (" + strings.Join(results, ", ") + ")"
	}

	var receiver string
	if sf.recv != nil {
		receiver = "(" + argumentList(sf.recv)[0] + ") "
	}

	return fmt.Sprintf("%s%s(%s)%s", receiver, sf.Name, strings.Join(params, ", "), result)
}

type Value struct {
	Name       string
	Value      string
	IsConstant bool
	Doc        string
}

func (sv Value) Signature() string {
	t := "var"
	if sv.IsConstant {
		t = "const"
	}

	return fmt.Sprintf("%s %s = %s", t, sv.Name, sv.Value)
}

type Values map[string]*Value

func (m Values) SortedKeys() []string {
	keys := make([]string, len(m))
	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func newValue(name, value string) *Value {
	return &Value{
		Name:       name,
		Value:      value,
		IsConstant: false,
	}
}
