package main

import (
	"context"
	"fmt"
	"os/exec"
	"taxifare"
	"taxiride"
	"util"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/log"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/filter"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
)

const (
	faresFile = "data/nycTaxiFares"
	ridesFile = "data/nycTaxiRides"
	output    = "joinedTaxiRidesAndFares.txt"
)

func keyData(s beam.Scope, rides, fares beam.PCollection) (beam.PCollection, beam.PCollection) {
	s = s.Scope("keyData")
	keyedRides := beam.ParDo(s, func(in taxiride.Ride) (int, string) { return in.ID(), in.String() }, rides)
	keyedFares := beam.ParDo(s, func(in taxifare.Fare) (int, string) { return in.ID(), in.String() }, fares)

	return keyedRides, keyedFares
}

func parseData(s beam.Scope, rides, fares beam.PCollection) (beam.PCollection, beam.PCollection) {
	s = s.Scope("parseData")
	parsedRides := beam.ParDo(s, taxiride.Parse, rides)
	parsedFares := beam.ParDo(s, taxifare.Parse, fares)

	return parsedRides, parsedFares
}

func filterData(s beam.Scope, rides, fares beam.PCollection) (beam.PCollection, beam.PCollection) {
	s = s.Scope("filterData")
	filteredRides := filter.Include(s, rides, func(r taxiride.Ride) bool {
		return r.IsDurationLess(util.WorkDay)
	})

	return filteredRides, fares
}

func processOutside(s beam.Scope, rides, fares beam.PCollection) (beam.PCollection, beam.PCollection) {
	s = s.Scope("processOutside")
	extProcessFn := func(f taxifare.Fare) taxifare.Fare {
		cmd := exec.Command("sleep", "0.1")
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
		return f
	}

	beam.ParDo(s, extProcessFn, fares)
	return rides, fares
}

func processInside(s beam.Scope, rides, fares beam.PCollection) (beam.PCollection, beam.PCollection) {
	s = s.Scope("processInside")
	extProcessFn := func(f taxifare.Fare) taxifare.Fare {
		// time.Sleep(100 * time.Millisecond)
		return f
	}

	beam.ParDo(s, extProcessFn, fares)
	return rides, fares
}

func joinEvents(s beam.Scope, rides, fares beam.PCollection) beam.PCollection {
	s = s.Scope("joinEvents")
	joined := beam.CoGroupByKey(s, rides, fares)
	result := beam.ParDo(s, processFn, joined)
	return beam.ParDo(s, formatFn, result)
}

func processFn(rideID int, rides, fares func(*string) bool, emit func(int, string)) {
	fare := "none"
	fares(&fare) // grab first (and only) country name, if any

	var ride string
	for rides(&ride) {
		emit(rideID, fmt.Sprintf("Ride info: %v, Fare info: %v", fare, ride))
	}
}

func formatFn(rideID int, info string) string {
	return fmt.Sprintf("Ride ID: %v, %v", rideID, info)
}

func getSrcData(s beam.Scope) (beam.PCollection, beam.PCollection) {
	s = s.Scope("getSrcData")
	srcRides := textio.Read(s, ridesFile)
	srcFares := textio.Read(s, faresFile)

	return srcRides, srcFares
}

func deduplicate(s beam.Scope, rides, fares beam.PCollection) (beam.PCollection, beam.PCollection) {
	s = s.Scope("deduplicate")
	uniqRides := filter.Distinct(s, rides)
	uniqFares := filter.Distinct(s, fares)

	return uniqRides, uniqFares
}

func main() {
	beam.Init()

	ctx := context.Background()
	p, s := beam.NewPipelineWithRoot()

	log.Info(ctx, "Reading files")
	srcRides, srcFares := getSrcData(s)

	log.Info(ctx, "Removing duplicates")
	uniqRides, uniqFares := deduplicate(s, srcRides, srcFares)

	log.Info(ctx, "Parsing data")
	parsedRides, parsedFares := parseData(s, uniqRides, uniqFares)

	log.Info(ctx, "Filtering data")
	filteredRides, filteredFares := filterData(s, parsedRides, parsedFares)

	// log.Info(ctx, "Processing data")
	// processedRides, processedFares := processInside(s, filteredRides, filteredFares)

	log.Info(ctx, "Preparing data for join")
	// keyedRides, keyedFares := keyData(s, processedRides, processedFares)
	keyedRides, keyedFares := keyData(s, filteredRides, filteredFares)

	log.Info(ctx, "Running join")
	formatted := joinEvents(s, keyedRides, keyedFares)

	s = s.Scope("writeData")
	textio.Write(s, output, formatted)

	if err := beamx.Run(ctx, p); err != nil {
		log.Exitf(ctx, "Failed to execute job: %v", err)
	}
}
