package walidator_test

import (
	"testing"

	"github.com/heetch/walidator"
)

func BenchmarkValidate(b *testing.B) {
	type t1 struct {
		ID        string  `json:"id" validate:"nonzero"`
		State     string  `json:"state" validate:"nonzero"`
		FooID     string  `json:"foo_id" validate:"nonzero"`
		UserID    string  `json:"user_id" validate:"nonzero"`
		Latitude  float64 `json:"origin_latitude" validate:"nonzero"`
		Longitude float64 `json:"origin_longitude" validate:"nonzero"`
		XXXID     string  `json:"xxx_id"`
	}
	type T2 struct {
		X *t1 `json:"ride" validate:"nonzero"`
	}
	v := &T2{
		X: &t1{
			ID:        "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			State:     "yes",
			FooID:     "somefoo",
			UserID:    "someuser",
			Latitude:  23,
			Longitude: 45,
		},
	}
	for i := 0; i < b.N; i++ {
		err := walidator.Validate(v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUUID(b *testing.B) {
	v := struct {
		ID string `json:"id" validate:"uuid"`
	}{
		ID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	}
	for i := 0; i < b.N; i++ {
		err := walidator.Validate(v)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLatitude(b *testing.B) {
	v := struct {
		OriginLatitude float64 `json:"origin_latitude" validate:"nonzero,latitude"`
	}{
		OriginLatitude: 23,
	}
	for i := 0; i < b.N; i++ {
		err := walidator.Validate(v)
		if err != nil {
			b.Fatal(err)
		}
	}
}
