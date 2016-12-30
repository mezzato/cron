package cron

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mezzato/astrotime"
)

// ConstantDelaySchedule represents a simple recurring duty cycle, e.g. "Every 5 minutes".
// It does not support jobs more frequent than once a second.
type AstroSchedule struct {
	Latitude  float64
	Longitude float64
	Offset    time.Duration
}

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (schedule *AstroSchedule) Next(t time.Time) (r time.Time) {

	// Start at the earliest possible time (the upcoming second).
	// t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

	// take the midnight plus 1 sec as a reference
	calcTime := t.Truncate(time.Hour * 24).Add(time.Second)

	for {
		//fmt.Printf("offset:%v\n", schedule.Offset)
		r = astrotime.CalcSunset(calcTime, schedule.Latitude, schedule.Longitude).Add(schedule.Offset)

		// second rounding
		r = r.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

		//duration - time.Duration(duration.Nanoseconds())%time.Second

		//fmt.Printf("astro time in schedule:%v\n", sunsetWithOffset)
		if r.After(t) {
			break
		}
		calcTime = calcTime.Add(time.Hour * 24)
	}

	return r
}

var astroSpecRegexp *regexp.Regexp = regexp.MustCompile(`@(?P<type>[a-zA-Z]+)\((?P<latitude>\d+\.?\d*)[\s,]+(?P<longitude>\d+\.?\d*)[\s,]*(?P<offset>.+)?\)`)

// ParseAstroSpec returns an AstroSchedule
// @sunset(latitude float64, longitude float64, offset time.Duration?)
// note that offset is optional and 0 by default
func ParseAstroSpec(spec string) (*AstroSchedule, error) {

	matches := astroSpecRegexp.FindStringSubmatch(spec)
	if matches == nil {
		return nil, fmt.Errorf("Failed to parse astronomical schedule string %s, no matches", spec)
	}

	astroType := strings.TrimSpace(matches[1])
	latitudeS := strings.TrimSpace(matches[2])
	longitudeS := strings.TrimSpace(matches[3])
	offsetS := strings.TrimSpace(matches[4])

	switch astroType {
	case "sunset":
		lat, err1 := strconv.ParseFloat(latitudeS, 64)
		if err1 != nil {
			return nil, fmt.Errorf("Failed to parse astronomical schedule string %s, invalid latitude", spec)
		}
		long, err2 := strconv.ParseFloat(longitudeS, 64)
		if err2 != nil {
			return nil, fmt.Errorf("Failed to parse astronomical schedule string %s, invalid latitude", spec)
		}
		var offset time.Duration
		if offsetS != "" {
			offset, err2 = time.ParseDuration(offsetS)
			if err2 != nil {
				return nil, fmt.Errorf("Failed to parse astronomical schedule string %s, invalid offset", spec)
			}
		}

		// round to the sec
		offset = offset - time.Duration(offset.Nanoseconds())%time.Second

		return &AstroSchedule{
			Latitude:  lat,
			Longitude: long,
			Offset:    offset,
		}, nil

	default:
		return nil, fmt.Errorf("Failed to parse astronomical schedule string %s. Invalid type %s", spec, astroType)
	}

	//astrotime.CalcSunset(t, latitude, longitude)
}
