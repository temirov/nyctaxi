package util

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	cases := []struct {
		in   string
		want time.Time
	}{
		{
			"1970-01-01 00:00:00",
			time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			"2013-01-01 00:00:00",
			time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			"2013-01-01 00:00:01",
			time.Date(2013, 1, 1, 0, 0, 1, 0, time.UTC),
		},
	}

	for _, c := range cases {
		got := ParseTime(c.in)

		if !c.want.Equal(got) {
			t.Errorf("ParseTime(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseTimePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	ParseTime("nonsense")
}
