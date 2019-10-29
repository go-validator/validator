// Package validator implements value validations
//
// Copyright 2014 Roberto Teixeira <robteix@robteix.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package walidator_test

import (
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/heetch/walidator/v4"
)

type TestStruct struct {
	A          int    `validate:"nonzero"`
	B          string `validate:"len=8,min=6,max=4"`
	Sub        TestStructSub
	D          *Simple `validate:"nonzero"`
	E          I       `validate:"nonzero"`
	unexported int     `validate:"nonzero"`
}

type TestStructSub struct {
	A int `validate:"nonzero"`
	B string
	C float64 `validate:"nonzero,min=1"`
	D *string `validate:"nonzero"`
}

type Simple struct {
	A int `validate:"min=10"`
}

type I interface {
	Foo() string
}

type Impl struct {
	F string `validate:"len=3"`
}

func (i *Impl) Foo() string {
	return i.F
}

type tree struct {
	Val         int `validate:"min=1"`
	Left, Right *tree
}

type Mer interface {
	M()
}

type T1 struct{}

func (t *T1) M() {}

type T2 struct {
	Mer Mer
}

type Str string

var validateTests = []struct {
	testName      string
	tagValidators map[string]walidator.TagValidator
	value         interface{}
	expectError   walidator.Errors
}{{
	testName: "validate-TestStruct",
	value: TestStruct{
		A: 0,
		B: "12345",
		Sub: TestStructSub{
			A: 1,
		},
		D: &Simple{10},
		E: &Impl{"hello"},
	},
	expectError: walidator.Errors{
		"A":     {{"nonzero", "zero value"}},
		"B":     {{"len", "invalid length"}, {"min", "less than min"}, {"max", "greater than max"}},
		"Sub.C": {{"nonzero", "zero value"}, {"min", "less than min"}},
		"Sub.D": {{"nonzero", "zero value"}},
		"E.F":   {{"len", "invalid length"}},
	},
}, {
	testName: "nil value",
	value:    nil,
}, {
	testName: "zero-length-nonzero-slice",
	value: struct {
		S []int `validate:"nonzero"`
	}{
		S: make([]int, 0, 10),
	},
	expectError: walidator.Errors{
		"S": {{"nonzero", "zero value"}},
	},
}, {
	testName: "slice-other-validators",
	value: struct {
		S []int `validate:"min=11,max=5,len=9,nonzero"`
	}{
		S: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	},
	expectError: walidator.Errors{
		"S": {{"min", "less than min"}, {"max", "greater than max"}, {"len", "invalid length"}},
	},
}, {
	testName: "empty-validation-name",
	value: struct {
		S int `validate:",nonzero"`
	}{2},
	expectError: walidator.Errors{
		"S": {{"invalid", "empty validation name in tag"}},
	},
}, {
	testName: "empty-validate-tag",
	value: struct {
		S int `validate:""`
	}{2},
}, {
	testName: "nonzero-fail",
	value: struct {
		S0  int          `validate:"nonzero"`
		S1  int8         `validate:"nonzero"`
		S2  int16        `validate:"nonzero"`
		S3  int32        `validate:"nonzero"`
		S4  int64        `validate:"nonzero"`
		S5  uint         `validate:"nonzero"`
		S6  uint8        `validate:"nonzero"`
		S7  uint16       `validate:"nonzero"`
		S8  uint32       `validate:"nonzero"`
		S9  uint64       `validate:"nonzero"`
		S10 float64      `validate:"nonzero"`
		S11 float32      `validate:"nonzero"`
		S12 uintptr      `validate:"nonzero"`
		S13 bool         `validate:"nonzero"`
		S14 *int         `validate:"nonzero"`
		S15 *interface{} `validate:"nonzero"`
		S16 struct{}     `validate:"nonzero"`
		S17 chan int     `validate:"nonzero"`
	}{},
	expectError: walidator.Errors{
		"S0":  {{"nonzero", "zero value"}},
		"S1":  {{"nonzero", "zero value"}},
		"S2":  {{"nonzero", "zero value"}},
		"S3":  {{"nonzero", "zero value"}},
		"S4":  {{"nonzero", "zero value"}},
		"S5":  {{"nonzero", "zero value"}},
		"S6":  {{"nonzero", "zero value"}},
		"S7":  {{"nonzero", "zero value"}},
		"S8":  {{"nonzero", "zero value"}},
		"S9":  {{"nonzero", "zero value"}},
		"S10": {{"nonzero", "zero value"}},
		"S11": {{"nonzero", "zero value"}},
		"S12": {{"nonzero", "zero value"}},
		"S13": {{"nonzero", "zero value"}},
		"S14": {{"nonzero", "zero value"}},
		"S15": {{"nonzero", "zero value"}},
		"S17": {{"nonzero", "unsupported type"}},
	},
}, {
	testName: "nonzero-ok",
	value: struct {
		S0  int          `validate:"nonzero"`
		S1  int8         `validate:"nonzero"`
		S2  int16        `validate:"nonzero"`
		S3  int32        `validate:"nonzero"`
		S4  int64        `validate:"nonzero"`
		S5  uint         `validate:"nonzero"`
		S6  uint8        `validate:"nonzero"`
		S7  uint16       `validate:"nonzero"`
		S8  uint32       `validate:"nonzero"`
		S9  uint64       `validate:"nonzero"`
		S10 float64      `validate:"nonzero"`
		S11 float32      `validate:"nonzero"`
		S12 uintptr      `validate:"nonzero"`
		S13 bool         `validate:"nonzero"`
		S14 *int         `validate:"nonzero"`
		S15 *interface{} `validate:"nonzero"`
		S16 struct{}     `validate:"nonzero"`
	}{
		S0:  1,
		S1:  1,
		S2:  1,
		S3:  1,
		S4:  1,
		S5:  1,
		S6:  1,
		S7:  1,
		S8:  1,
		S9:  1,
		S10: 1,
		S11: 1,
		S12: 1,
		S13: true,
		S14: new(int),
		S15: new(interface{}),
	},
}, {
	testName: "nonzero-empty-map",
	value: struct {
		S map[string]string `validate:"nonzero"`
	}{
		S: make(map[string]string),
	},
	expectError: walidator.Errors{
		"S": {{"nonzero", "zero value"}},
	},
}, {
	testName: "nonzero-unsigned-ok",
	value: struct {
		S uint
	}{
		S: 1,
	},
}, {
	testName: "nonzero-unsigned-fail",
	value: struct {
		S uint `validate:"nonzero"`
	}{},
	expectError: walidator.Errors{
		"S": {{"nonzero", "zero value"}},
	},
}, {
	testName: "map-min-length",
	value: struct {
		S map[string]string `validate:"min=1"`
	}{
		S: make(map[string]string),
	},
	expectError: walidator.Errors{
		"S": {{"min", "less than min"}},
	},
}, {
	testName: "map-max-length",
	value: struct {
		S map[string]string `validate:"max=1"`
	}{
		S: map[string]string{"A": "a", "B": "a"},
	},
	expectError: walidator.Errors{
		"S": {{"max", "greater than max"}},
	},
}, {
	testName: "map-multiple-failures",
	value: struct {
		S map[string]string `validate:"len=4,min=6,max=1,nonzero"`
	}{
		S: map[string]string{
			"1": "a",
			"2": "b",
			"3": "c",
			"4": "d",
			"5": "e",
		},
	},
	expectError: walidator.Errors{
		"S": {{"len", "invalid length"}, {"min", "less than min"}, {"max", "greater than max"}},
	},
}, {
	testName: "nonzero-float",
	value: struct {
		S float64 `validate:"nonzero"`
	}{
		S: 12.34,
	},
}, {
	testName: "zero-float",
	value: struct {
		S float64 `validate:"nonzero"`
	}{},
	expectError: walidator.Errors{
		"S": {{"nonzero", "zero value"}},
	},
}, {
	testName: "nonzero-int",
	value: struct {
		S int `validate:"nonzero"`
	}{
		S: 123,
	},
}, {
	testName: "min-ok",
	value: struct {
		S0  int         `validate:"min=3"`
		S1  int8        `validate:"min=3"`
		S2  int16       `validate:"min=3"`
		S3  int32       `validate:"min=3"`
		S4  int64       `validate:"min=3"`
		S5  uint        `validate:"min=3"`
		S6  uint8       `validate:"min=3"`
		S7  uint16      `validate:"min=3"`
		S8  uint32      `validate:"min=3"`
		S9  uint64      `validate:"min=3"`
		S10 float64     `validate:"min=3"`
		S11 float32     `validate:"min=3"`
		S12 uintptr     `validate:"min=3"`
		S13 *int        `validate:"min=3"`
		S14 interface{} `validate:"min=3"`
		S15 interface{} `validate:"min=3"`
		S16 interface{} `validate:"min=3"`
	}{
		S0:  5,
		S1:  5,
		S2:  5,
		S3:  5,
		S4:  5,
		S5:  5,
		S6:  5,
		S7:  5,
		S8:  5,
		S9:  5,
		S10: 5,
		S11: 5,
		S12: 5,
		S13: newInt(5),
		S14: 5,
		S15: 5.0,
		S16: nil,
	},
}, {
	testName: "min-fail",
	value: struct {
		S0  int         `validate:"min=3"`
		S1  int8        `validate:"min=3"`
		S2  int16       `validate:"min=3"`
		S3  int32       `validate:"min=3"`
		S4  int64       `validate:"min=3"`
		S5  uint        `validate:"min=3"`
		S6  uint8       `validate:"min=3"`
		S7  uint16      `validate:"min=3"`
		S8  uint32      `validate:"min=3"`
		S9  uint64      `validate:"min=3"`
		S10 float64     `validate:"min=3"`
		S11 float32     `validate:"min=3"`
		S12 uintptr     `validate:"min=3"`
		S13 *int        `validate:"min=3"`
		S14 interface{} `validate:"min=3"`
		S15 interface{} `validate:"min=3"`
		S16 interface{} `validate:"min=3"`
		S17 interface{} `validate:"min=3"`
	}{
		S0:  2,
		S1:  2,
		S2:  2,
		S3:  2,
		S4:  2,
		S5:  2,
		S6:  2,
		S7:  2,
		S8:  2,
		S9:  2,
		S10: 2,
		S11: 2,
		S12: 2,
		S13: newInt(2),
		S14: 2,
		S15: 2.0,
		S16: newInt(2),
		S17: new(chan int),
	},
	expectError: walidator.Errors{
		"S0":  {{"min", "less than min"}},
		"S1":  {{"min", "less than min"}},
		"S2":  {{"min", "less than min"}},
		"S3":  {{"min", "less than min"}},
		"S4":  {{"min", "less than min"}},
		"S5":  {{"min", "less than min"}},
		"S6":  {{"min", "less than min"}},
		"S7":  {{"min", "less than min"}},
		"S8":  {{"min", "less than min"}},
		"S9":  {{"min", "less than min"}},
		"S10": {{"min", "less than min"}},
		"S11": {{"min", "less than min"}},
		"S12": {{"min", "less than min"}},
		"S13": {{"min", "less than min"}},
		"S14": {{"min", "less than min"}},
		"S15": {{"min", "less than min"}},
		"S16": {{"min", "less than min"}},
		"S17": {{"min", "unsupported type"}},
	},
}, {
	testName: "max-ok",
	value: struct {
		S0  int         `validate:"max=3"`
		S1  int8        `validate:"max=3"`
		S2  int16       `validate:"max=3"`
		S3  int32       `validate:"max=3"`
		S4  int64       `validate:"max=3"`
		S5  uint        `validate:"max=3"`
		S6  uint8       `validate:"max=3"`
		S7  uint16      `validate:"max=3"`
		S8  uint32      `validate:"max=3"`
		S9  uint64      `validate:"max=3"`
		S10 float64     `validate:"max=3"`
		S11 float32     `validate:"max=3"`
		S12 uintptr     `validate:"max=3"`
		S13 *int        `validate:"max=3"`
		S14 interface{} `validate:"max=3"`
		S15 interface{} `validate:"max=3"`
	}{
		S0:  2,
		S1:  2,
		S2:  2,
		S3:  2,
		S4:  2,
		S5:  2,
		S6:  2,
		S7:  2,
		S8:  2,
		S9:  2,
		S10: 2,
		S11: 2,
		S12: 2,
		S13: newInt(2),
		S14: 2,
		S15: 2.0,
	},
}, {
	testName: "max-fail",
	value: struct {
		S0  int         `validate:"max=3"`
		S1  int8        `validate:"max=3"`
		S2  int16       `validate:"max=3"`
		S3  int32       `validate:"max=3"`
		S4  int64       `validate:"max=3"`
		S5  uint        `validate:"max=3"`
		S6  uint8       `validate:"max=3"`
		S7  uint16      `validate:"max=3"`
		S8  uint32      `validate:"max=3"`
		S9  uint64      `validate:"max=3"`
		S10 float64     `validate:"max=3"`
		S11 float32     `validate:"max=3"`
		S12 uintptr     `validate:"max=3"`
		S13 *int        `validate:"max=3"`
		S14 interface{} `validate:"max=3"`
		S15 interface{} `validate:"max=3"`
		S16 uint        `validate:"max=-3"`
		S17 int         `validate:"max=2.3"`
		S18 float64     `validate:"max=x"`
		S19 chan int    `validate:"max=3"`
	}{
		S0:  5,
		S1:  5,
		S2:  5,
		S3:  5,
		S4:  5,
		S5:  5,
		S6:  5,
		S7:  5,
		S8:  5,
		S9:  5,
		S10: 5,
		S11: 5,
		S12: 5,
		S13: newInt(5),
		S14: 5,
		S15: 5.0,
	},
	expectError: walidator.Errors{
		"S0":  {{"max", "greater than max"}},
		"S1":  {{"max", "greater than max"}},
		"S2":  {{"max", "greater than max"}},
		"S3":  {{"max", "greater than max"}},
		"S4":  {{"max", "greater than max"}},
		"S5":  {{"max", "greater than max"}},
		"S6":  {{"max", "greater than max"}},
		"S7":  {{"max", "greater than max"}},
		"S8":  {{"max", "greater than max"}},
		"S9":  {{"max", "greater than max"}},
		"S10": {{"max", "greater than max"}},
		"S11": {{"max", "greater than max"}},
		"S12": {{"max", "greater than max"}},
		"S13": {{"max", "greater than max"}},
		"S14": {{"max", "greater than max"}},
		"S15": {{"max", "greater than max"}},
		"S16": {{"max", "bad parameter"}},
		"S17": {{"max", "bad parameter"}},
		"S18": {{"max", "bad parameter"}},
		"S19": {{"max", "unsupported type"}},
	},
}, {
	testName: "min-max-int",
	value: struct {
		S int `validate:"min=124, max=122"`
	}{
		S: 123,
	},
	expectError: walidator.Errors{
		"S": {{"min", "less than min"}, {"max", "greater than max"}},
	},
}, {
	testName: "max-int",
	value: struct {
		S int `validate:"max=10"`
	}{
		S: 123,
	},
	expectError: walidator.Errors{
		"S": {{"max", "greater than max"}},
	},
}, {
	testName: "string-len-equal",
	value: struct {
		S string `validate:"len=8"`
	}{
		S: "test1234",
	},
}, {
	testName: "string-len-equal-fail",
	value: struct {
		S string `validate:"len=0"`
	}{
		S: "test1234",
	},
	expectError: walidator.Errors{
		"S": {{"len", "invalid length"}},
	},
}, {
	testName: "string-regexp-match",
	value: struct {
		S string `validate:"regexp=^[tes]{4}.*"`
	}{
		S: "test1234",
	},
}, {
	testName: "string-regexp-no-match",
	value: struct {
		S string `validate:"regexp=^.*[0-9]{5}$"`
	}{
		S: "test1234",
	},
	expectError: walidator.Errors{
		"S": {{"regexp", "regular expression mismatch"}},
	},
}, {
	testName: "regexp-bad-pattern",
	value: struct {
		S string `validate:"regexp=)"`
	}{},
	expectError: walidator.Errors{
		"S": {{"regexp", "bad parameter"}},
	},
}, {
	testName: "string-multiple-failure",
	value: struct {
		S string `validate:"nonzero,len=3,max=1"`
	}{
		S: "",
	},
	expectError: walidator.Errors{
		"S": {{"nonzero", "zero value"}, {"len", "invalid length"}},
	},
}, {
	testName: "custom-struct-validator-ok",
	tagValidators: map[string]walidator.TagValidator{
		"struct": structValidator,
	},
	value: struct {
		S struct{ A int } `validate:"struct"`
	}{},
}, {
	testName: "custom-struct-validator-fail",
	tagValidators: map[string]walidator.TagValidator{
		"struct": structValidator,
	},
	value: struct {
		S int `validate:"struct"`
	}{},
	expectError: walidator.Errors{
		"S": {{"struct", "unsupported non-struct type"}},
	},
}, {
	testName: "custom-struct-validator-fail",
	tagValidators: map[string]walidator.TagValidator{
		"struct": structValidator,
	},
	value: struct {
		S int `validate:"struct"`
	}{},
	expectError: walidator.Errors{
		"S": {{"struct", "unsupported non-struct type"}},
	},
}, {
	testName: "unrecognized-validator",
	value: struct {
		S int `validate:"foo"`
	}{},
	expectError: walidator.Errors{
		"S": {{"invalid", `unknown validation name "foo" in tag`}},
	},
}, {
	testName: "custom-validator-nil",
	tagValidators: map[string]walidator.TagValidator{
		"nil": nilValidator,
	},
	value: struct {
		S *struct {
			A int
		} `validate:"nil"`
	}{},
}, {
	testName: "custom-validator-nil-fail",
	tagValidators: map[string]walidator.TagValidator{
		"nil": nilValidator,
	},
	value: struct {
		S *struct {
			A int
		} `validate:"nil"`
	}{
		S: &struct{ A int }{56},
	},
	expectError: walidator.Errors{
		"S": {{"nil", "unsupported non-nil value"}},
	},
}, {
	testName: "pointer-value-checks",
	value: struct {
		A *string `validate:"min=6"`
		B *[]int  `validate:"len=7"`
		C *int    `validate:"min=12"`
	}{
		A: newString("aaa"),
		B: new([]int),
	},
	expectError: walidator.Errors{
		"A": {{"min", "less than min"}},
		"B": {{"len", "invalid length"}},
	},
}, {
	testName: "struct-with-slice",
	value: struct {
		Slices []Simple
	}{
		Slices: []Simple{{1}, {20}, {5}},
	},
	expectError: walidator.Errors{
		"Slices[0].A": {{"min", "less than min"}},
		"Slices[2].A": {{"min", "less than min"}},
	},
}, {
	testName: "omitted-field",
	value: struct {
		A Simple `validate:"-"`
	}{
		// Note: this would be an invalid Simple
		// value if the "-" tag on the parent wasn't
		// respected.
		A: Simple{},
	},
}, {
	testName: "nested-slice",
	value: struct {
		Slices [][]Simple
	}{
		Slices: [][]Simple{{{1}, {20}}, {{5}}},
	},
	expectError: walidator.Errors{
		"Slices[0][0].A": {{"min", "less than min"}},
		"Slices[1][0].A": {{"min", "less than min"}},
	},
}, {
	testName: "map",
	value: struct {
		Map          map[string]Simple
		StructKeyMap map[Simple]Simple
	}{
		Map: map[string]Simple{
			"hello": {6},
		},
		StructKeyMap: map[Simple]Simple{
			{3}: {1},
		},
	},
	expectError: walidator.Errors{
		"Map[hello](value).A":          {{"min", "less than min"}},
		"StructKeyMap[{A:3}](key).A":   {{"min", "less than min"}},
		"StructKeyMap[{A:3}](value).A": {{"min", "less than min"}},
	},
}, {
	testName: "regexp-unsupported-type",
	value: struct {
		A int     `validate:"regexp=a.*b"`
		B float64 `validate:"regexp=.*"`
	}{},
	expectError: walidator.Errors{
		"A": {{"regexp", "unsupported type"}},
		"B": {{"regexp", "unsupported type"}},
	},
}, {
	testName: "bad-parameter",
	value: struct {
		A string `validate:"min="`
		B string `validate:"len=="`
		C string `validate:"max=foo"`
	}{},
	expectError: walidator.Errors{
		"A": {{"min", "bad parameter"}},
		"B": {{"len", "bad parameter"}},
		"C": {{"max", "bad parameter"}},
	},
}, {
	testName: "tag-escape-ok",
	value: struct {
		A string `validate:"min=0,regexp=^a{3\\,10}"`
	}{A: "aaaa"},
}, {
	testName: "tag-escape-fail",
	value: struct {
		A string `validate:"min=0,regexp=^a{3\\,10}"`
	}{A: "aa"},
	expectError: walidator.Errors{
		"A": {{"regexp", "regular expression mismatch"}},
	},
}, {
	testName: "json-tag",
	value: struct {
		A string `validate:"nonzero" json:"a,omitempty"`
		B string `validate:"nonzero" json:"b"`
	}{},
	expectError: walidator.Errors{
		"a": {{"nonzero", "zero value"}},
		"b": {{"nonzero", "zero value"}},
	},
}, {
	testName: "recursive-type-ok",
	value: &tree{
		Val: 1,
		Left: &tree{
			Val: 2,
			Right: &tree{
				Val: 3,
			},
		},
		Right: &tree{
			Val: 4,
		},
	},
}, {
	testName: "recursive-type-fail",
	value: &tree{
		Left: &tree{
			Right: &tree{},
		},
		Right: &tree{
			Val: 4,
		},
	},
	expectError: walidator.Errors{
		"Left.Right.Val": {{"min", "less than min"}},
		"Left.Val":       {{"min", "less than min"}},
		"Val":            {{"min", "less than min"}},
	},
}, {
	testName: "interface-field-ok",
	value: struct {
		X interface{} `validate:"max=1.0"`
	}{
		X: 0.3,
	},
}, {
	testName: "interface-ptr-ok",
	value: struct {
		X interface{} `validate:"max=1"`
	}{
		X: newInt(1),
	},
}, {
	testName: "interface-field-fail",
	value: struct {
		X interface{} `validate:"max=1.0"`
	}{
		X: "hello",
	},
	expectError: walidator.Errors{
		"X": {{"max", "bad parameter"}},
	},
}, {
	testName: "uuid-ok",
	value: struct {
		ID0 string `validator:"uuid"`
		ID1 Str    `validator:"uuid"`
	}{
		ID0: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		ID1: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	},
}, {
	testName: "uuid-fail-1",
	value: struct {
		ID string `validate:"uuid"`
	}{
		ID: "1234",
	},
	expectError: walidator.Errors{
		"ID": {{"uuid", "regular expression mismatch"}},
	},
}, {
	testName: "uuid-fail-2",
	value: struct {
		ID    string `validate:"uuid"`
		Other int    `validate:"uuid"`
	}{
		ID: "0VCE98AC-1326-4C79-8EBC-94908DA8B034",
	},
	expectError: walidator.Errors{
		"ID":    {{"uuid", "regular expression mismatch"}},
		"Other": {{"uuid", "unsupported type"}},
	},
}, {
	testName: "required-ok",
	value: struct {
		A string            `validate:"required"`
		B []int             `validate:"required"`
		C []int             `validate:"required"`
		D map[string]int    `validate:"required"`
		E int               `validate:"required"`
		F float64           `validate:"required"`
		G T1                `validate:"required"`
		H struct{ Foo int } `validate:"required"`
	}{
		A: "string",
		B: []int{1, 2, 3}[1:],
		C: []int{1, 2, 3},
		D: map[string]int{"a": 1, "b": 2},
		E: 12,
		F: 2.1,
		G: T1{},
		H: struct{ Foo int }{23},
	},
}, {
	testName: "required-fail",
	value: struct {
		A *uint          `validate:"required"`
		B *T1            `validate:"required"`
		C interface{}    `validate:"required"`
		D Mer            `validate:"required"`
		E chan int       `validate:"required"`
		F map[string]int `validate:"required"`
		G []int          `validate:"required"`
	}{
		D: T2{}.Mer,
	},
	expectError: walidator.Errors{
		"A": {{"required", "required value"}},
		"B": {{"required", "required value"}},
		"C": {{"required", "required value"}},
		"D": {{"required", "required value"}},
		"E": {{"required", "unsupported type"}},
		"F": {{"required", "required value"}},
		"G": {{"required", "required value"}},
	},
}, {
	testName: "latitude-ok",
	value: struct {
		A float64  `validate:"latitude"`
		B float64  `validate:"latitude"`
		C float64  `validate:"latitude"`
		D string   `validate:"latitude"`
		E string   `validate:"latitude"`
		F string   `validate:"latitude"`
		G *float64 `validate:"latitude"`
		H *string  `validate:"latitude"`
	}{
		A: 21.0,
		B: 0,
		C: -10.212,
		D: "12.0",
		E: "0.0",
		F: "-1.22",
		G: newFloat64(1.1),
		H: newString("-12.1"),
	},
}, {
	testName: "latitude-fail",
	value: struct {
		A float64  `validate:"latitude"`
		B float64  `validate:"latitude"`
		C string   `validate:"latitude"`
		D string   `validate:"latitude"`
		E *float64 `validate:"latitude"`
		F *string  `validate:"latitude"`
		G string   `validate:"latitude"`
		H int      `validate:"latitude"`
	}{
		A: 210.0,
		B: -1000.212,
		C: "1220",
		D: "-190.2",
		E: newFloat64(220.21),
		F: newString("-2121.1"),
		G: "x",
	},
	expectError: walidator.Errors{
		"A": {{"latitude", "210 is not a valid latitude"}},
		"B": {{"latitude", "-1000.212 is not a valid latitude"}},
		"C": {{"latitude", "1220 is not a valid latitude"}},
		"D": {{"latitude", "-190.2 is not a valid latitude"}},
		"E": {{"latitude", "220.21 is not a valid latitude"}},
		"F": {{"latitude", "-2121.1 is not a valid latitude"}},
		"G": {{"latitude", `"x" is not a valid latitude`}},
		"H": {{"latitude", `unsupported type`}},
	},
}, {
	testName: "longitude-ok",
	value: struct {
		A float64  `validate:"longitude"`
		B float64  `validate:"longitude"`
		C float64  `validate:"longitude"`
		D string   `validate:"longitude"`
		E string   `validate:"longitude"`
		F string   `validate:"longitude"`
		G *float64 `validate:"longitude"`
		H *string  `validate:"longitude"`
	}{
		A: 21.0,
		B: 0,
		C: -10.212,
		D: "12.0",
		E: "0.0",
		F: "-1.22",
		G: newFloat64(1.1),
		H: newString("-12.1"),
	},
}, {
	testName: "longitude-fail",
	value: struct {
		A float64  `validate:"longitude"`
		B float64  `validate:"longitude"`
		C string   `validate:"longitude"`
		D string   `validate:"longitude"`
		E *float64 `validate:"longitude"`
		F *string  `validate:"longitude"`
		G string   `validate:"longitude"`
		H int      `validate:"longitude"`
	}{
		A: 210.0,
		B: -1000.212,
		C: "1220",
		D: "-190.2",
		E: newFloat64(220.21),
		F: newString("-2121.1"),
		G: "x",
	},
	expectError: walidator.Errors{
		"A": {{"longitude", "210 is not a valid longitude"}},
		"B": {{"longitude", "-1000.212 is not a valid longitude"}},
		"C": {{"longitude", "1220 is not a valid longitude"}},
		"D": {{"longitude", "-190.2 is not a valid longitude"}},
		"E": {{"longitude", "220.21 is not a valid longitude"}},
		"F": {{"longitude", "-2121.1 is not a valid longitude"}},
		"G": {{"longitude", `"x" is not a valid longitude`}},
		"H": {{"longitude", `unsupported type`}},
	},
}, {
	testName: "len-invalid-type",
	value: struct {
		S int `validate:"len=3"`
	}{},
	expectError: walidator.Errors{
		"S": {{"len", "unsupported type"}},
	},
}, {
	testName: "min-bad-params",
	value: struct {
		S0 int         `validate:"min=2.3"`
		S1 float64     `validate:"min=x"`
		S2 uint        `validate:"min=x"`
		S3 chan string `validate:"min=1"`
	}{},
	expectError: walidator.Errors{
		"S0": {{"min", "bad parameter"}},
		"S1": {{"min", "bad parameter"}},
		"S2": {{"min", "bad parameter"}},
		"S3": {{"min", "unsupported type"}},
	},
}}

func TestValidate(t *testing.T) {
	c := qt.New(t)
	for i, test := range validateTests {
		c.Run(test.testName, func(c *qt.C) {
			v := walidator.New()
			for name, f := range test.tagValidators {
				v.AddValidation(name, f)
			}
			err := v.Validate(test.value)
			if test.expectError != nil {
				c.Assert(err, qt.DeepEquals, test.expectError, qt.Commentf("%T %T", err, test.expectError))
			} else {
				c.Assert(err, qt.Equals, nil, qt.Commentf("test %d", i))
			}
		})
	}
}

func TestCopy(t *testing.T) {
	c := qt.New(t)
	v := walidator.New()
	// WithTag copies the validator, so shouldn't pollute the original.
	v2 := v.WithTag("validate")
	// now we add a custom func only to the second one, it shouldn't get added
	// to the first
	v2.AddValidation("custom", func(reflect.Type, string) (walidator.ValidationFunc, error) {
		return func(reflect.Value, *walidator.ErrorReporter) {}, nil
	})
	type test struct {
		A string `validate:"custom"`
	}
	err := v2.Validate(test{})
	c.Assert(err, qt.Equals, nil)

	err = v.Validate(test{})
	c.Assert(err, qt.DeepEquals, walidator.Errors{
		"A": {{"invalid", `unknown validation name "custom" in tag`}},
	})
}

func TestAddValidation(t *testing.T) {
	c := qt.New(t)
	nopValidator := func(reflect.Type, string) (walidator.ValidationFunc, error) {
		return func(reflect.Value, *walidator.ErrorReporter) {}, nil
	}
	v := walidator.New()
	c.Assert(func() {
		v.AddValidation("", nopValidator)
	}, qt.PanicMatches, `empty validation name`)
	c.Assert(func() {
		v.AddValidation("x", nil)
	}, qt.PanicMatches, `nil validation function`)
	c.Assert(func() {
		v.AddValidation("len", nopValidator)
	}, qt.PanicMatches, `validation function len already registered`)
}

func TestErrorsError(t *testing.T) {
	c := qt.New(t)
	c.Assert(make(walidator.Errors), qt.ErrorMatches, "no validation errors")

	c.Assert(walidator.Errors{
		"":  {{"x", "y"}, {"z", "a"}},
		"b": {{"1", "4"}},
	}, qt.ErrorMatches, "validation error: y")
}

func structValidator(t reflect.Type, _ string) (walidator.ValidationFunc, error) {
	return func(v reflect.Value, r *walidator.ErrorReporter) {
		if v.Kind() != reflect.Struct {
			r.Errorf("unsupported non-struct type")
		}
	}, nil
}

func nilValidator(t reflect.Type, _ string) (walidator.ValidationFunc, error) {
	return func(v reflect.Value, r *walidator.ErrorReporter) {
		if !v.IsNil() {
			r.Errorf("unsupported non-nil value")
		}
	}, nil
}

func newString(x string) *string {
	return &x
}

func newInt(x int) *int {
	return &x
}

func newFloat64(x float64) *float64 {
	return &x
}
