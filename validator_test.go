// Package validator implements value validations
//
// Copyright (C) 2014-2016 Roberto Teixeira <robteix@robteix.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package validator_test

import (
	"testing"

	. "gopkg.in/check.v1"
	"gopkg.in/validator.v2"
)

func Test(t *testing.T) {
	TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})

type Simple struct {
	A int `validate:"min=10"`
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

	valid, err := validator.Validate(t)
	c.Assert(valid, Equals, false, Commentf("err: %v", err))
	c.Assert(validator.IsZeroValue(err, "A"), Equals, true)
	c.Assert(validator.Errors(err, "B"), HasLen, 3)
	c.Assert(validator.IsMin(err, "B"), Equals, true)
	c.Assert(validator.IsMax(err, "B"), Equals, true)
	c.Assert(validator.Errors(err, "Sub.A"), HasLen, 0)
	c.Assert(validator.Errors(err, "Sub.B"), HasLen, 0)
	c.Assert(validator.Errors(err, "Sub.C"), HasLen, 2)
	c.Assert(validator.IsZeroValue(err, "Sub.D"), Equals, true)
}

func (ms *MySuite) TestValidSlice(c *C) {
	s := make([]int, 0, 10)
	valid, errs := validator.Valid(s, "nonzero")
	c.Assert(valid, Equals, false)
	c.Assert(errs, HasError, validator.ErrZeroValue)

	for i := 0; i < 10; i++ {
		s = append(s, i)
	}

	_, errs = validator.Valid(s, "min=11,max=5,len=9,nonzero")
	c.Assert(errs, HasError, validator.ErrMin)
	c.Assert(errs, HasError, validator.ErrMax)
	c.Assert(errs, HasError, validator.ErrLen)
	c.Assert(errs, Not(HasError), validator.ErrZeroValue)
}

func (ms *MySuite) TestValidMap(c *C) {
	m := make(map[string]string)
	valid, errs := validator.Valid(m, "nonzero")
	c.Assert(valid, Equals, false)
	c.Assert(errs, HasError, validator.ErrZeroValue)

	_, errs = validator.Valid(m, "min=1")
	c.Assert(errs, HasError, validator.ErrMin)

	m = map[string]string{"A": "a", "B": "a"}
	_, errs = validator.Valid(m, "max=1")
	c.Assert(errs, HasError, validator.ErrMax)

	valid, _ = validator.Valid(m, "min=2, max=5")
	c.Assert(valid, Equals, true)

	m = map[string]string{
		"1": "a",
		"2": "b",
		"3": "c",
		"4": "d",
		"5": "e",
	}
	valid, errs = validator.Valid(m, "len=4,min=6,max=1,nonzero")
	c.Assert(valid, Equals, false)
	c.Assert(errs, HasError, validator.ErrLen)
	c.Assert(errs, HasError, validator.ErrMin)
	c.Assert(errs, HasError, validator.ErrMax)
	c.Assert(errs, Not(HasError), validator.ErrZeroValue)

}

func (ms *MySuite) TestValidFloat(c *C) {
	valid, _ := validator.Valid(12.34, "nonzero")
	c.Assert(valid, Equals, true)

	valid, errs := validator.Valid(0.0, "nonzero")
	c.Assert(valid, Equals, false)
	c.Assert(errs, HasError, validator.ErrZeroValue)
}

func (ms *MySuite) TestValidInt(c *C) {
	i := 123
	valid, errs := validator.Valid(i, "nonzero")
	c.Assert(valid, Equals, true)
	c.Assert(errs, Not(HasError), validator.ErrZeroValue)

	valid, _ = validator.Valid(i, "min=1")
	c.Assert(valid, Equals, true)

	valid, errs = validator.Valid(i, "min=124, max=122")
	c.Assert(valid, Equals, false)
	c.Assert(errs, HasError, validator.ErrMin)
	c.Assert(errs, HasError, validator.ErrMax)

	_, errs = validator.Valid(i, "max=10")
	c.Assert(errs, HasError, validator.ErrMax)
}

func (ms *MySuite) TestValidString(c *C) {
	s := "test1234"
	valid, errs := validator.Valid(s, "len=8")
	c.Assert(valid, Equals, true)
	c.Assert(errs, HasLen, 0)

	_, errs = validator.Valid(s, "len=0")
	c.Assert(errs, HasError, validator.ErrLen)

	_, errs = validator.Valid(s, `"regexp=^[a-z]{4,5}.*"`)
	c.Assert(errs, HasLen, 0)

	_, errs = validator.Valid(s, "regexp=^[tes]{4}.*")
	c.Assert(errs, HasLen, 0)

	_, errs = validator.Valid(s, "regexp=^.*[0-9]{5}$")
	c.Assert(errs, NotNil)

	valid, errs = validator.Valid("", "nonzero,len=3,max=1")
	c.Assert(valid, Equals, false)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs, HasError, validator.ErrZeroValue)
	c.Assert(errs, HasError, validator.ErrLen)
	c.Assert(errs, Not(HasError), validator.ErrMax)
}

func (ms *MySuite) TestValidateStructVar(c *C) {
	type test struct {
		A int
	}
	t := test{}
	valid, errs := validator.Valid(t, "")
	c.Assert(valid, Equals, false)
	c.Assert(errs, HasError, validator.ErrUnsupported)
}

func (ms *MySuite) TestUnknownTag(c *C) {
	type test struct {
		A int `validate:"foo"`
	}
	t := test{}
	valid, err := validator.Validate(t)
	c.Assert(valid, Equals, false)
	c.Assert(validator.Errors(err, "A"), HasLen, 1)
	c.Assert(validator.IsUnknownTag(err, "A"), Equals, true)
}

func (ms *MySuite) TestUnsupported(c *C) {
	type test struct {
		A int     `validate:"regexp=a.*b"`
		B float64 `validate:"regexp=.*"`
	}
	t := test{}
	valid, err := validator.Validate(t)
	c.Assert(valid, Equals, false)
	c.Assert(validator.IsUnsupported(err, "A"), Equals, true)
	c.Assert(validator.IsUnsupported(err, "B"), Equals, true)
}

func (ms *MySuite) TestBadParameter(c *C) {
	type test struct {
		A string `validate:"min="`
		B string `validate:"len=="`
		C string `validate:"max=foo"`
	}
	t := test{}
	valid, err := validator.Validate(t)
	c.Assert(valid, Equals, false)
	c.Assert(validator.IsBadParameter(err, "A"), Equals, true)
	c.Assert(validator.IsBadParameter(err, "B"), Equals, true)
	c.Assert(validator.IsBadParameter(err, "C"), Equals, true)
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
	slice, ok = params[0].([]error)
	if !ok {
		return false, "First parameter is not a []error"
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
