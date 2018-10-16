package taxiride

import (
	"github.com/golang/geo/s2"
	"reflect"
	"testing"
	"util"
)

func TestParse(t *testing.T) {
	cases := []struct {
		in   string
		want Ride
	}{
		{
			"6,START,2013-01-01 00:00:00,1970-01-01 00:00:00,-73.866135,40.771091,-73.961334,40.764912,6,2013000006,2013000006",
			Ride{
				RideID:       6,
				RideType:     "START",
				StartTime:    util.ParseTime("2013-01-01 00:00:00"),
				EndTime:      util.ParseTime("1970-01-01 00:00:00"),
				StartPoint:   s2.LatLng{40.771091, -73.866135},
				EndPoint:     s2.LatLng{40.764912, -73.961334},
				PassengerCnt: 6,
				TaxiID:       2013000006,
				DriverID:     2013000006,
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
			"6,START,2013-01-01 00:00:00,1970-01-01 00:00:00,-73.866135,40.771091,-73.961334,40.764912,6,2013000006,2013000006",
			"Duration: 376944h0m0s; Length: 5.4515798",
		},
	}

	for _, c := range cases {
		ride := Parse(c.in)
		got := ride.String()

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
			"6,START,2013-01-01 00:00:00,1970-01-01 00:00:00,-73.866135,40.771091,-73.961334,40.764912,6,2013000006,2013000006",
			6,
		},
	}

	for _, c := range cases {
		ride := Parse(c.in)
		got := ride.ID()

		if got != c.want {
			t.Errorf("Parse(%q).String() == %q, want %q", c.in, got, c.want)
		}
	}
}
