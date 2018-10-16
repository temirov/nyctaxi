package taxiride

import (
	"fmt"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"strconv"
	"strings"
	"time"
	"util"
)

type Ride struct {
	RideID     int       // 0
	RideType   string    // 1
	StartTime  time.Time // 2
	EndTime    time.Time // 3
	StartPoint s2.LatLng // 0
	EndPoint   s2.LatLng // 0
	// startLon float64  // 4
	// startLat float64  // 5
	// endLon float64    // 6
	// endLat float64    // 7
	PassengerCnt int // 8
	TaxiID       int // 9
	DriverID     int // 10
}

func (r *Ride) String() string {
	return fmt.Sprintf("Duration: %s; Length: %s", r.Duration(), r.length())
}

func (r *Ride) ID() int {
	return r.RideID
}

func Parse(in string) Ride {
	tmp := strings.Split(in, ",")

	rideId, _ := strconv.Atoi(tmp[0])
	rideType := tmp[1]
	startTime := util.ParseTime(tmp[2])
	endTime := util.ParseTime(tmp[3])
	startLon, _ := strconv.ParseFloat(tmp[4], 64)
	startLat, _ := strconv.ParseFloat(tmp[5], 64)
	endLon, _ := strconv.ParseFloat(tmp[6], 64)
	endLat, _ := strconv.ParseFloat(tmp[7], 64)
	startPoint := s2.LatLng{Lat: s1.Angle(startLat), Lng: s1.Angle(startLon)}
	endPoint := s2.LatLng{Lat: s1.Angle(endLat), Lng: s1.Angle(endLon)}
	passengerCnt, _ := strconv.Atoi(tmp[8])
	taxiId, _ := strconv.Atoi(tmp[9])
	driverId, _ := strconv.Atoi(tmp[10])

	ride := Ride{
		rideId,
		rideType,
		startTime,
		endTime,
		startPoint,
		endPoint,
		passengerCnt,
		taxiId,
		driverId,
	}
	return ride
}

func (r *Ride) IsDurationMore(d time.Duration) bool {
	return r.Duration() > d
}

func (r *Ride) IsDurationLess(d time.Duration) bool {
	return !r.IsDurationMore(d)
}

func (r *Ride) Duration() time.Duration {
	return r.StartTime.Sub(r.EndTime)
}

func (r *Ride) length() s1.Angle {
	return r.EndPoint.Distance(r.StartPoint)
}
