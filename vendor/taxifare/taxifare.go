package taxifare

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"util"
)

type Fare struct {
	RideID      int
	TaxiID      int
	DriverID    int
	StartTime   time.Time
	PaymentType string
	Tip         float64
	Tolls       float64
	TotalFare   float64
}

func (f *Fare) String() string {
	return fmt.Sprintf("Total: %.2f; Tip: %.2f", f.TotalFare, f.Tip)
}

func (f *Fare) ID() int {
	return f.RideID
}

func Parse(in string) Fare {
	tmp := strings.Split(in, ",")

	rideId, _ := strconv.Atoi(tmp[0])
	taxiId, _ := strconv.Atoi(tmp[1])
	driverId, _ := strconv.Atoi(tmp[2])
	startTime := util.ParseTime(tmp[3])
	paymentType := tmp[4]
	tip, _ := strconv.ParseFloat(tmp[5], 64)
	tolls, _ := strconv.ParseFloat(tmp[6], 64)
	totalFare, _ := strconv.ParseFloat(tmp[7], 64)

	fare := Fare{
		rideId,
		taxiId,
		driverId,
		startTime,
		paymentType,
		tip,
		tolls,
		totalFare,
	}
	return fare
}
