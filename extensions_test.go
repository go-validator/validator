// Package validator_test test value validations
//
// Copyright 2018 Heetch
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
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/heetch/walidator"
)

func TestUUIDOK(t *testing.T) {
	c := qt.New(t)
	err := walidator.Valid("6ba7b810-9dad-11d1-80b4-00c04fd430c8", "uuid")
	c.Assert(err, qt.IsNil)
}

func TestUUIDNOK(t *testing.T) {
	c := qt.New(t)
	cases := []string{
		"1234",
		"0VCE98AC-1326-4C79-8EBC-94908DA8B034",
	}
	for _, s := range cases {
		err := walidator.Valid(s, "uuid")
		c.Assert(err, qt.Not(qt.Equals), nil)
		errs, ok := err.(walidator.ErrorArray)
		c.Assert(ok, qt.Equals, true)
		c.Assert(errs, qt.HasLen, 1)
		c.Assert(errs, qt.Contains, walidator.ErrRegexp)
	}
}

type Mer interface {
	M()
}

type T1 struct{}

func (t *T1) M() {}

type T2 struct {
	Mer Mer
}

func TestRequiredOK(t *testing.T) {
	c := qt.New(t)
	a := []int{1, 2, 3}
	cases := []interface{}{
		"string",
		a[1:],
		a,
		map[string]int{"a": 1, "b": 2},
		12,
		2.1,
		T1{},
		struct{ Foo int }{23},
	}
	for _, s := range cases {
		err := walidator.Valid(s, "required")
		c.Assert(err, qt.IsNil)
	}
}

func TestRequiredNOK(t *testing.T) {
	c := qt.New(t)
	var ptr *uint
	var t1 *T1
	t2 := T2{}
	cases := []interface{}{
		ptr,
		t1,
		nil,
		t2.Mer,
	}
	for i, s := range cases {
		c.Logf("test %d: %#v", i, s)
		err := walidator.Valid(s, "required")
		c.Assert(err, qt.Not(qt.Equals), nil)
		errs, ok := err.(walidator.ErrorArray)
		c.Assert(ok, qt.Equals, true)
		c.Assert(errs, qt.HasLen, 1)
		c.Assert(errs, qt.Contains, walidator.ErrRequired)
	}
}

func TestLatitudeOK(t *testing.T) {
	c := qt.New(t)
	v := 1.1
	s := "-12.1"
	cases := []interface{}{
		21.0,
		0.0,
		-10.212,
		"12.0",
		"0.0",
		"-1.22",
		&v,
		&s,
	}
	for _, l := range cases {
		err := walidator.Valid(l, "latitude")
		c.Assert(err, qt.IsNil)
	}
}

func TestLatitudeNOK(t *testing.T) {
	c := qt.New(t)
	v := 220.21
	s := "-2121.1"
	cases := []interface{}{
		210.0,
		-1000.212,
		"1220",
		"-190.2",
		&v,
		&s,
	}
	for _, l := range cases {
		err := walidator.Valid(l, "latitude")
		c.Assert(err, qt.Not(qt.Equals), nil)
		errs, ok := err.(walidator.ErrorArray)
		c.Assert(ok, qt.Equals, true)
		c.Assert(errs, qt.HasLen, 1)
		switch loc := l.(type) {
		case *float64:
			c.Assert(errs[0].Error(), qt.Equals, fmt.Sprintf("%v is not a valid latitude", *loc))
		case *string:
			c.Assert(errs[0].Error(), qt.Equals, fmt.Sprintf("%v is not a valid latitude", *loc))
		default:
			c.Assert(errs[0].Error(), qt.Equals, fmt.Sprintf("%v is not a valid latitude", loc))
		}
	}
}

func TestLongitudeOK(t *testing.T) {
	c := qt.New(t)
	v := 1.1
	s := "-12.1"
	cases := []interface{}{
		21.0,
		0.0,
		-10.212,
		"12.0",
		"0.0",
		"-1.22",
		&v,
		&s,
	}
	for _, l := range cases {
		err := walidator.Valid(l, "latitude")
		c.Assert(err, qt.IsNil)
	}
}

func TestLongitudeNOK(t *testing.T) {
	c := qt.New(t)
	v := 220.21
	s := "-2121.1"
	cases := []interface{}{
		210.0,
		-1000.212,
		"1220",
		"-190.2",
		&v,
		&s,
	}
	for _, l := range cases {
		err := walidator.Valid(l, "longitude")
		c.Assert(err, qt.Not(qt.Equals), nil)
		errs, ok := err.(walidator.ErrorArray)
		c.Assert(ok, qt.Equals, true)
		c.Assert(errs, qt.HasLen, 1)
		switch loc := l.(type) {
		case *float64:
			c.Assert(errs[0].Error(), qt.Equals, fmt.Sprintf("%v is not a valid longitude", *loc))
		case *string:
			c.Assert(errs[0].Error(), qt.Equals, fmt.Sprintf("%v is not a valid longitude", *loc))
		default:
			c.Assert(errs[0].Error(), qt.Equals, fmt.Sprintf("%v is not a valid longitude", loc))
		}
	}
}
