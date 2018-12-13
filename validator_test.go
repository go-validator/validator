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

	"github.com/heetch/walidator"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})

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

func (ms *MySuite) TestValidate(c *C) {
	t := TestStruct{
		A: 0,
		B: "12345",
	}
	t.Sub.A = 1
	t.Sub.B = ""
	t.Sub.C = 0.0
	t.D = &Simple{10}
	t.E = &Impl{"hello"}

	err := walidator.Validate(t)
	c.Assert(err, NotNil)

	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasError, walidator.ErrZeroValue)
	c.Assert(errs["B"], HasError, walidator.ErrLen)
	c.Assert(errs["B"], HasError, walidator.ErrMin)
	c.Assert(errs["B"], HasError, walidator.ErrMax)
	c.Assert(errs["Sub.A"], HasLen, 0)
	c.Assert(errs["Sub.B"], HasLen, 0)
	c.Assert(errs["Sub.C"], HasLen, 2)
	c.Assert(errs["Sub.D"], HasError, walidator.ErrZeroValue)
	c.Assert(errs["E.F"], HasError, walidator.ErrLen)
}

func (ms *MySuite) TestValidSlice(c *C) {
	s := make([]int, 0, 10)
	err := walidator.Valid(s, "nonzero")
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrZeroValue)

	for i := 0; i < 10; i++ {
		s = append(s, i)
	}

	err = walidator.Valid(s, "min=11,max=5,len=9,nonzero")
	c.Assert(err, NotNil)
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrMin)
	c.Assert(errs, HasError, walidator.ErrMax)
	c.Assert(errs, HasError, walidator.ErrLen)
	c.Assert(errs, Not(HasError), walidator.ErrZeroValue)
}

func (ms *MySuite) TestValidMap(c *C) {
	m := make(map[string]string)
	err := walidator.Valid(m, "nonzero")
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrZeroValue)

	err = walidator.Valid(m, "min=1")
	c.Assert(err, NotNil)
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrMin)

	m = map[string]string{"A": "a", "B": "a"}
	err = walidator.Valid(m, "max=1")
	c.Assert(err, NotNil)
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrMax)

	err = walidator.Valid(m, "min=2, max=5")
	c.Assert(err, IsNil)

	m = map[string]string{
		"1": "a",
		"2": "b",
		"3": "c",
		"4": "d",
		"5": "e",
	}
	err = walidator.Valid(m, "len=4,min=6,max=1,nonzero")
	c.Assert(err, NotNil)
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrLen)
	c.Assert(errs, HasError, walidator.ErrMin)
	c.Assert(errs, HasError, walidator.ErrMax)
	c.Assert(errs, Not(HasError), walidator.ErrZeroValue)

}

func (ms *MySuite) TestValidFloat(c *C) {
	err := walidator.Valid(12.34, "nonzero")
	c.Assert(err, IsNil)

	err = walidator.Valid(0.0, "nonzero")
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrZeroValue)
}

func (ms *MySuite) TestValidInt(c *C) {
	i := 123
	err := walidator.Valid(i, "nonzero")
	c.Assert(err, IsNil)

	err = walidator.Valid(i, "min=1")
	c.Assert(err, IsNil)

	err = walidator.Valid(i, "min=124, max=122")
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrMin)
	c.Assert(errs, HasError, walidator.ErrMax)

	err = walidator.Valid(i, "max=10")
	c.Assert(err, NotNil)
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrMax)
}

func (ms *MySuite) TestValidString(c *C) {
	s := "test1234"
	err := walidator.Valid(s, "len=8")
	c.Assert(err, IsNil)

	err = walidator.Valid(s, "len=0")
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, walidator.ErrLen)

	err = walidator.Valid(s, "regexp=^[tes]{4}.*")
	c.Assert(err, IsNil)

	err = walidator.Valid(s, "regexp=^.*[0-9]{5}$")
	c.Assert(errs, NotNil)

	err = walidator.Valid("", "nonzero,len=3,max=1")
	c.Assert(err, NotNil)
	errs, ok = err.(walidator.ErrorArray)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs, HasError, walidator.ErrZeroValue)
	c.Assert(errs, HasError, walidator.ErrLen)
	c.Assert(errs, Not(HasError), walidator.ErrMax)
}

func (ms *MySuite) TestValidateStructVar(c *C) {
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
	c.Assert(err, IsNil)

	type test2 struct {
		B int
	}
	type test1 struct {
		A test2 `validate:"struct"`
	}

	err = walidator.Validate(test1{})
	c.Assert(err, IsNil)

	type test4 struct {
		B int `validate:"foo"`
	}
	type test3 struct {
		A test4
	}
	err = walidator.Validate(test3{})
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A.B"], HasError, walidator.ErrUnknownTag)
}

func (ms *MySuite) TestValidatePointerVar(c *C) {
	// just verifies that a the given val is a struct
	walidator.SetValidationFunc("struct", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Struct {
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
	c.Assert(err, IsNil)

	type test2 struct {
		B int
	}
	type test1 struct {
		A *test2 `validate:"struct"`
	}

	err = walidator.Validate(&test1{&test2{}})
	c.Assert(err, IsNil)

	type test4 struct {
		B int `validate:"foo"`
	}
	type test3 struct {
		A test4
	}
	err = walidator.Validate(&test3{})
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A.B"], HasError, walidator.ErrUnknownTag)

	err = walidator.Valid((*test)(nil), "nil")
	c.Assert(err, IsNil)

	type test5 struct {
		A *test2 `validate:"nil"`
	}
	err = walidator.Validate(&test5{})
	c.Assert(err, IsNil)

	type test6 struct {
		A *test2 `validate:"nonzero"`
	}
	err = walidator.Validate(&test6{})
	errs, ok = err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasError, walidator.ErrZeroValue)

	err = walidator.Validate(&test6{&test2{}})
	c.Assert(err, IsNil)

	type test7 struct {
		A *string `validate:"min=6"`
		B *int    `validate:"len=7"`
		C *int    `validate:"min=12"`
	}
	s := "aaa"
	b := 8
	err = walidator.Validate(&test7{&s, &b, nil})
	errs, ok = err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasError, walidator.ErrMin)
	c.Assert(errs["B"], HasError, walidator.ErrLen)
	c.Assert(errs["C"], Not(HasError), walidator.ErrMin)
}

func (ms *MySuite) TestValidateOmittedStructVar(c *C) {
	type test2 struct {
		B int `validate:"min=1"`
	}
	type test1 struct {
		A test2 `validate:"-"`
	}

	t := test1{}
	err := walidator.Validate(t)
	c.Assert(err, IsNil)

	errs := walidator.Valid(test2{}, "-")
	c.Assert(errs, IsNil)
}

func (ms *MySuite) TestUnknownTag(c *C) {
	type test struct {
		A int `validate:"foo"`
	}
	t := test{}
	err := walidator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs["A"], HasError, walidator.ErrUnknownTag)
}

func (ms *MySuite) TestValidateStructWithSlice(c *C) {
	type test2 struct {
		Num    int    `validate:"max=2"`
		String string `validate:"nonzero"`
	}

	type test struct {
		Slices []test2 `validate:"len=1"`
	}

	t := test{
		Slices: []test2{{
			Num:    6,
			String: "foo",
		}},
	}
	err := walidator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["Slices[0].Num"], HasError, walidator.ErrMax)
	c.Assert(errs["Slices[0].String"], IsNil) // sanity check
}

func (ms *MySuite) TestValidateStructWithNestedSlice(c *C) {
	type test2 struct {
		Num int `validate:"max=2"`
	}

	type test struct {
		Slices [][]test2
	}

	t := test{
		Slices: [][]test2{{{Num: 6}}},
	}
	err := walidator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["Slices[0][0].Num"], HasError, walidator.ErrMax)
}

func (ms *MySuite) TestValidateStructWithMap(c *C) {
	type test2 struct {
		Num int `validate:"max=2"`
	}

	type test struct {
		Map          map[string]test2
		StructKeyMap map[test2]test2
	}

	t := test{
		Map: map[string]test2{
			"hello": {Num: 6},
		},
		StructKeyMap: map[test2]test2{
			{Num: 3}: {Num: 1},
		},
	}
	err := walidator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)

	c.Assert(errs["Map[hello](value).Num"], HasError, walidator.ErrMax)
	c.Assert(errs["StructKeyMap[{Num:3}](key).Num"], HasError, walidator.ErrMax)
}

func (ms *MySuite) TestUnsupported(c *C) {
	type test struct {
		A int     `validate:"regexp=a.*b"`
		B float64 `validate:"regexp=.*"`
	}
	t := test{}
	err := walidator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs["A"], HasError, walidator.ErrUnsupported)
	c.Assert(errs["B"], HasError, walidator.ErrUnsupported)
}

func (ms *MySuite) TestBadParameter(c *C) {
	type test struct {
		A string `validate:"min="`
		B string `validate:"len=="`
		C string `validate:"max=foo"`
	}
	t := test{}
	err := walidator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 3)
	c.Assert(errs["A"], HasError, walidator.ErrBadParameter)
	c.Assert(errs["B"], HasError, walidator.ErrBadParameter)
	c.Assert(errs["C"], HasError, walidator.ErrBadParameter)
}

func (ms *MySuite) TestCopy(c *C) {
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
	c.Assert(err, IsNil)

	err = v.Validate(test{})
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs["A"], HasError, walidator.ErrUnknownTag)
}

func (ms *MySuite) TestTagEscape(c *C) {
	type test struct {
		A string `validate:"min=0,regexp=^a{3\\,10}"`
	}
	t := test{"aaaa"}
	err := walidator.Validate(t)
	c.Assert(err, IsNil)

	t2 := test{"aa"}
	err = walidator.Validate(t2)
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasError, walidator.ErrRegexp)
}

func (ms *MySuite) TestJSONTag(c *C) {
	type test struct {
		A string `validate:"nonzero" json:"b,omitempty"`
	}

	var t test
	err := walidator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(walidator.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], IsNil)
	c.Assert(errs["b"], HasError, walidator.ErrZeroValue)
}

type hasErrorChecker struct {
	*CheckerInfo
}

func (c *hasErrorChecker) Check(params []interface{}, names []string) (bool, string) {
	var (
		ok    bool
		slice []error
		value error
	)
	slice, ok = params[0].(walidator.ErrorArray)
	if !ok {
		return false, "First parameter is not an Errorarray"
	}
	value, ok = params[1].(error)
	if !ok {
		return false, "Second parameter is not an error"
	}

	for _, v := range slice {
		if v == value {
			return true, ""
		}
	}
	return false, ""
}

func (c *hasErrorChecker) Info() *CheckerInfo {
	return c.CheckerInfo
}

var HasError = &hasErrorChecker{&CheckerInfo{Name: "HasError", Params: []string{"HasError", "expected to contain"}}}
