package main

import (
	"reflect"
	"taxifare"
	"testing"
)

type iTaxiRide interface {
	ID() int
	String() string
}

func keyDataFn(in iTaxiRide) (int, string) {
	return in.ID(), in.String()
}

type result struct {
	id   int
	desc string
}

func TestKeyDataFn(t *testing.T) {
	cases := []struct {
		in   interface{}
		want result
	}{
		{
			taxifare.Parse("1,2013000001,2013000001,2013-01-01 00:00:00,CSH,0,0,21.5"),
			result{1, "Total: 21.50; Tip: 0.00"},
		},
	}

	for _, c := range cases {
		id, desc := keyDataFn(c.in)
		got := result{id, desc}

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("keyDataFn(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
