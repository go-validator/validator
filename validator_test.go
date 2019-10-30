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
	"encoding/json"
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/heetch/walidator"
)

type Simple struct {
	A int `validate:"min=10"`
}

type I interface {
	Foo() string
}

type Impl struct {
	F string `validate:"len=3"`
}

func (this *Impl) Foo() string {
	return this.F
}

type TestStruct struct {
	A   int    `validate:"nonzero"`
	B   string `validate:"len=8,min=6,max=4"`
	Sub struct {
		A int `validate:"nonzero"`
		B string
		C float64 `validate:"nonzero,min=1"`
		D *string `validate:"nonzero"`
	}
	D *Simple `validate:"nonzero"`
	E I       `validate:nonzero`
}

func TestValidate(t *testing.T) {
	c := qt.New(t)
	v := TestStruct{
		A: 0,
		B: "12345",
	}
	v.Sub.A = 1
	v.Sub.B = ""
	v.Sub.C = 0.0
	v.D = &Simple{10}
	v.E = &Impl{"hello"}

	err := walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))

	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["A"], qt.Contains, walidator.ErrZeroValue)
	c.Assert(errs["B"], qt.Contains, walidator.ErrLen)
	c.Assert(errs["B"], qt.Contains, walidator.ErrMin)
	c.Assert(errs["B"], qt.Contains, walidator.ErrMax)
	c.Assert(errs["Sub.A"], qt.HasLen, 0)
	c.Assert(errs["Sub.B"], qt.HasLen, 0)
	c.Assert(errs["Sub.C"], qt.HasLen, 2)
	c.Assert(errs["Sub.D"], qt.Contains, walidator.ErrZeroValue)
	c.Assert(errs["E.F"], qt.Contains, walidator.ErrLen)
}

func TestValidSlice(t *testing.T) {
	c := qt.New(t)
	s := make([]int, 0, 10)
	err := walidator.Valid(s, "nonzero")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrZeroValue)

	for i := 0; i < 10; i++ {
		s = append(s, i)
	}

	err = walidator.Valid(s, "min=11,max=5,len=9,nonzero")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrMin)
	c.Assert(errs, qt.Contains, walidator.ErrMax)
	c.Assert(errs, qt.Contains, walidator.ErrLen)
	c.Assert(errs, qt.Not(qt.Contains), walidator.ErrZeroValue)
}

func TestValidMap(t *testing.T) {
	c := qt.New(t)
	m := make(map[string]string)
	err := walidator.Valid(m, "nonzero")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrZeroValue)

	err = walidator.Valid(m, "min=1")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrMin)

	m = map[string]string{"A": "a", "B": "a"}
	err = walidator.Valid(m, "max=1")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrMax)

	err = walidator.Valid(m, "min=2, max=5")
	c.Assert(err, qt.IsNil)

	m = map[string]string{
		"1": "a",
		"2": "b",
		"3": "c",
		"4": "d",
		"5": "e",
	}
	err = walidator.Valid(m, "len=4,min=6,max=1,nonzero")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrLen)
	c.Assert(errs, qt.Contains, walidator.ErrMin)
	c.Assert(errs, qt.Contains, walidator.ErrMax)
	c.Assert(errs, qt.Not(qt.Contains), walidator.ErrZeroValue)

}

func TestValidFloat(t *testing.T) {
	c := qt.New(t)
	err := walidator.Valid(12.34, "nonzero")
	c.Assert(err, qt.IsNil)

	err = walidator.Valid(0.0, "nonzero")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrZeroValue)
}

func TestValidInt(t *testing.T) {
	c := qt.New(t)
	i := 123
	err := walidator.Valid(i, "nonzero")
	c.Assert(err, qt.IsNil)

	err = walidator.Valid(i, "min=1")
	c.Assert(err, qt.IsNil)

	err = walidator.Valid(i, "min=124, max=122")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrMin)
	c.Assert(errs, qt.Contains, walidator.ErrMax)

	err = walidator.Valid(i, "max=10")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrMax)
}

func TestValidString(t *testing.T) {
	c := qt.New(t)
	s := "test1234"
	err := walidator.Valid(s, "len=8")
	c.Assert(err, qt.IsNil)

	err = walidator.Valid(s, "len=0")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.Contains, walidator.ErrLen)

	err = walidator.Valid(s, "regexp=^[tes]{4}.*")
	c.Assert(err, qt.IsNil)

	err = walidator.Valid(s, "regexp=^.*[0-9]{5}$")
	c.Assert(errs, qt.Not(qt.IsNil))

	err = walidator.Valid("", "nonzero,len=3,max=1")
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.HasLen, 2)
	c.Assert(errs, qt.Contains, walidator.ErrZeroValue)
	c.Assert(errs, qt.Contains, walidator.ErrLen)
	c.Assert(errs, qt.Not(qt.Contains), walidator.ErrMax)
}

func TestValidateStructVar(t *testing.T) {
	c := qt.New(t)
	// just verifies that a the given val is a struct
	walidator.SetValidationFunc("struct", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Struct {
			return nil
		}
		return walidator.ErrUnsupported
	})

	type test struct {
		A int
	}
	err := walidator.Valid(test{}, "struct")
	c.Assert(err, qt.IsNil)

	type test2 struct {
		B int
	}
	type test1 struct {
		A test2 `validate:"struct"`
	}

	err = walidator.Validate(test1{})
	c.Assert(err, qt.IsNil)

	type test4 struct {
		B int `validate:"foo"`
	}
	type test3 struct {
		A test4
	}
	err = walidator.Validate(test3{})
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["A.B"], qt.Contains, walidator.ErrUnknownTag)
}

func TestValidatePointerVar(t *testing.T) {
	c := qt.New(t)
	// just verifies that a the given val is a struct
	walidator.SetValidationFunc("struct", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
			return nil
		}
		return walidator.ErrUnsupported
	})
	walidator.SetValidationFunc("nil", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.IsNil() {
			return nil
		}
		return walidator.ErrUnsupported
	})

	type test struct {
		A int
	}
	err := walidator.Valid(&test{}, "struct")
	c.Assert(err, qt.IsNil)

	type test2 struct {
		B int
	}
	type test1 struct {
		A *test2 `validate:"struct"`
	}

	err = walidator.Validate(&test1{&test2{}})
	c.Assert(err, qt.IsNil)

	type test4 struct {
		B int `validate:"foo"`
	}
	type test3 struct {
		A test4
	}
	err = walidator.Validate(&test3{})
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["A.B"], qt.Contains, walidator.ErrUnknownTag)

	err = walidator.Valid((*test)(nil), "nil")
	c.Assert(err, qt.IsNil)

	type test5 struct {
		A *test2 `validate:"nil"`
	}
	err = walidator.Validate(&test5{})
	c.Assert(err, qt.IsNil)

	type test6 struct {
		A *test2 `validate:"nonzero"`
	}
	err = walidator.Validate(&test6{})
	errs, ok = err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["A"], qt.Contains, walidator.ErrZeroValue)

	err = walidator.Validate(&test6{&test2{}})
	c.Assert(err, qt.IsNil)

	type test7 struct {
		A *string `validate:"min=6"`
		B *int    `validate:"len=7"`
		C *int    `validate:"min=12"`
	}
	s := "aaa"
	b := 8
	err = walidator.Validate(&test7{&s, &b, nil})
	errs, ok = err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["A"], qt.Contains, walidator.ErrMin)
	c.Assert(errs["B"], qt.Contains, walidator.ErrLen)
	c.Assert(errs["C"], qt.Not(qt.Contains), walidator.ErrMin)
}

func TestValidateOmittedStructVar(t *testing.T) {
	c := qt.New(t)
	type test2 struct {
		B int `validate:"min=1"`
	}
	type test1 struct {
		A test2 `validate:"-"`
	}

	v := test1{}
	err := walidator.Validate(v)
	c.Assert(err, qt.IsNil)

	errs := walidator.Valid(test2{}, "-")
	c.Assert(errs, qt.IsNil)
}

func TestUnknownTag(t *testing.T) {
	c := qt.New(t)
	type test struct {
		A int `validate:"foo"`
	}
	v := test{}
	err := walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.HasLen, 1)
	c.Assert(errs["A"], qt.Contains, walidator.ErrUnknownTag)
}

func TestValidateStructWithSlice(t *testing.T) {
	c := qt.New(t)
	type test2 struct {
		Num    int    `validate:"max=2"`
		String string `validate:"nonzero"`
	}

	type test struct {
		Slices []test2 `validate:"len=1"`
	}

	v := test{
		Slices: []test2{{
			Num:    6,
			String: "foo",
		}},
	}
	err := walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["Slices[0].Num"], qt.Contains, walidator.ErrMax)
	c.Assert(errs["Slices[0].String"], qt.IsNil) // sanity check
}

func TestValidateStructWithNestedSlice(t *testing.T) {
	c := qt.New(t)
	type test2 struct {
		Num int `validate:"max=2"`
	}

	type test struct {
		Slices [][]test2
	}

	v := test{
		Slices: [][]test2{{{Num: 6}}},
	}
	err := walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["Slices[0][0].Num"], qt.Contains, walidator.ErrMax)
}

func TestValidateStructWithMap(t *testing.T) {
	c := qt.New(t)
	type test2 struct {
		Num int `validate:"max=2"`
	}

	type test struct {
		Map          map[string]test2
		StructKeyMap map[test2]test2
	}

	v := test{
		Map: map[string]test2{
			"hello": {Num: 6},
		},
		StructKeyMap: map[test2]test2{
			{Num: 3}: {Num: 1},
		},
	}
	err := walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)

	c.Assert(errs["Map[hello](value).Num"], qt.Contains, walidator.ErrMax)
	c.Assert(errs["StructKeyMap[{Num:3}](key).Num"], qt.Contains, walidator.ErrMax)
}

func TestUnsupported(t *testing.T) {
	c := qt.New(t)
	type test struct {
		A int     `validate:"regexp=a.*b"`
		B float64 `validate:"regexp=.*"`
	}
	v := test{}
	err := walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.HasLen, 2)
	c.Assert(errs["A"], qt.Contains, walidator.ErrUnsupported)
	c.Assert(errs["B"], qt.Contains, walidator.ErrUnsupported)
}

func TestBadParameter(t *testing.T) {
	c := qt.New(t)
	type test struct {
		A string `validate:"min="`
		B string `validate:"len=="`
		C string `validate:"max=foo"`
	}
	v := test{}
	err := walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.HasLen, 3)
	c.Assert(errs["A"], qt.Contains, walidator.ErrBadParameter)
	c.Assert(errs["B"], qt.Contains, walidator.ErrBadParameter)
	c.Assert(errs["C"], qt.Contains, walidator.ErrBadParameter)
}

func TestCopy(t *testing.T) {
	c := qt.New(t)
	v := walidator.NewValidator()
	// WithTag calls copy, so we just copy the validator with the same tag
	v2 := v.WithTag("validate")
	// now we add a custom func only to the second one, it shouldn't get added
	// to the first
	v2.SetValidationFunc("custom", func(_ interface{}, _ string) error { return nil })
	type test struct {
		A string `validate:"custom"`
	}
	err := v2.Validate(test{})
	c.Assert(err, qt.IsNil)

	err = v.Validate(test{})
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs, qt.HasLen, 1)
	c.Assert(errs["A"], qt.Contains, walidator.ErrUnknownTag)
}

func TestTagEscape(t *testing.T) {
	c := qt.New(t)
	type test struct {
		A string `validate:"min=0,regexp=^a{3\\,10}"`
	}
	v := test{"aaaa"}
	err := walidator.Validate(v)
	c.Assert(err, qt.IsNil)

	t2 := test{"aa"}
	err = walidator.Validate(t2)
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["A"], qt.Contains, walidator.ErrRegexp)
}

func TestJSONTag(t *testing.T) {
	c := qt.New(t)
	type test struct {
		A string `validate:"nonzero" json:"b,omitempty"`
	}

	var v test
	err := walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(errs["A"], qt.IsNil)
	c.Assert(errs["b"], qt.Contains, walidator.ErrZeroValue)
}

type tree struct {
	Val         int `validate:"min=1"`
	Left, Right *tree
}

func TestRecursiveType(t *testing.T) {
	c := qt.New(t)
	v := &tree{
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
	}
	err := walidator.Validate(v)
	c.Assert(err, qt.IsNil)
	v = &tree{
		Left: &tree{
			Right: &tree{},
		},
		Right: &tree{
			Val: 4,
		},
	}
	err = walidator.Validate(v)
	c.Assert(err, qt.Not(qt.IsNil))
	data, err := json.MarshalIndent(err, "", "\t")
	c.Assert(err, qt.IsNil)
	c.Assert(string(data), qt.Equals, `
{
	"Left.Right.Val": [
		"less than min"
	],
	"Left.Val": [
		"less than min"
	],
	"Val": [
		"less than min"
	]
}`[1:])
}

func TestInterfaceField(t *testing.T) {
	c := qt.New(t)
	var x struct {
		X interface{} `validate:"max=1.0"`
	}
	x.X = "hello"
	err := walidator.Validate(x)

	err1, ok := err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(err1["X"], qt.ErrorMatches, `bad parameter`)

	// Test the cached path.
	err = walidator.Validate(x)
	err1, ok = err.(walidator.ErrorMap)
	c.Assert(ok, qt.Equals, true)
	c.Assert(err1["X"], qt.ErrorMatches, `bad parameter`)

	// Test that it passes with an appropriate type
	x.X = 0.3
	err = walidator.Validate(x)
	c.Assert(err, qt.Equals, nil)

	// ... and a pointer to that type
	f := 0.9
	x.X = &f
	err = walidator.Validate(x)
	c.Assert(err, qt.Equals, nil)
}
