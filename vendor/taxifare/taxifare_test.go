package taxifare

import (
	"reflect"
	"testing"
	"util"
)

func TestParse(t *testing.T) {
	cases := []struct {
		in   string
		want Fare
	}{
		{
			"1,2013000001,2013000001,2013-01-01 00:00:00,CSH,0,0,21.5",
			Fare{
				RideID:      1,
				TaxiID:      2013000001,
				DriverID:    2013000001,
				StartTime:   util.ParseTime("2013-01-01 00:00:00"),
				PaymentType: "CSH",
				Tip:         0.0,
				Tolls:       0.0,
				TotalFare:   21.5,
			},
		},
	}
	for _, c := range cases {
		got := Parse(c.in)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Parse(%q) == %#v, \nwant %#v", c.in, got, c.want)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{
			"1,2013000001,2013000001,2013-01-01 00:00:00,CSH,0,0,21.5",
			"Total: 21.50; Tip: 0.00",
		},
	}

	for _, c := range cases {
		fare := Parse(c.in)
		got := fare.String()

		if got != c.want {
			t.Errorf("Parse(%q).String() == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestID(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{
			"1,2013000001,2013000001,2013-01-01 00:00:00,CSH,0,0,21.5",
			1,
		},
	}

	for _, c := range cases {
		fare := Parse(c.in)
		got := fare.ID()

		if got != c.want {
			t.Errorf("Parse(%q).String() == %q, want %q", c.in, got, c.want)
		}
	}
}
