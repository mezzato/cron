package cron

import "testing"

func TestAstroNext(t *testing.T) {
	tests := []struct {
		spec, time string
		expected   string
	}{
		// Simple cases
		{"@sunset(45.6309, 11.7065)", "Mon Jul 9 14:45 2012", "Mon Jul 9 19:03:01 2012"},
		{"@sunset(45.6309, 11.7065)", "Sun Jul 8 23:45 2012", "Mon Jul 9 19:03:01 2012"},
		{"@sunset(45.6309, 11.7065, 1h)", "Sat Feb 27 00:45 2016", "Sat Feb 27 17:56:13 2016"},
		{"@sunset(45.6309, 11.7065, -4m)", "Sat Feb 27 00:45 2016", "Sat Feb 27 16:52:13 2016"},
		/*
			{"Mon Jul 9 14:59 2012", 15 * time.Minute, "Mon Jul 9 15:14 2012"},
			{"Mon Jul 9 14:59:59 2012", 15 * time.Minute, "Mon Jul 9 15:14:59 2012"},

			// Wrap around hours
			{"Mon Jul 9 15:45 2012", 35 * time.Minute, "Mon Jul 9 16:20 2012"},

			// Wrap around days
			{"Mon Jul 9 23:46 2012", 14 * time.Minute, "Tue Jul 10 00:00 2012"},
			{"Mon Jul 9 23:45 2012", 35 * time.Minute, "Tue Jul 10 00:20 2012"},
			{"Mon Jul 9 23:35:51 2012", 44*time.Minute + 24*time.Second, "Tue Jul 10 00:20:15 2012"},
			{"Mon Jul 9 23:35:51 2012", 25*time.Hour + 44*time.Minute + 24*time.Second, "Thu Jul 11 01:20:15 2012"},

			// Wrap around months
			{"Mon Jul 9 23:35 2012", 91*24*time.Hour + 25*time.Minute, "Thu Oct 9 00:00 2012"},

			// Wrap around minute, hour, day, month, and year
			{"Mon Dec 31 23:59:45 2012", 15 * time.Second, "Tue Jan 1 00:00:00 2013"},

			// Round to nearest second on the delay
			{"Mon Jul 9 14:45 2012", 15*time.Minute + 50*time.Nanosecond, "Mon Jul 9 15:00 2012"},

			// Round up to 1 second if the duration is less.
			{"Mon Jul 9 14:45:00 2012", 15 * time.Millisecond, "Mon Jul 9 14:45:01 2012"},

			// Round to nearest second when calculating the next time.
			{"Mon Jul 9 14:45:00.005 2012", 15 * time.Minute, "Mon Jul 9 15:00 2012"},

			// Round to nearest second for both.
			{"Mon Jul 9 14:45:00.005 2012", 15*time.Minute + 50*time.Nanosecond, "Mon Jul 9 15:00 2012"},
		*/
	}

	for _, c := range tests {

		sched, err := ParseAstroSpec(c.spec)
		if err != nil {
			t.Error(err)
			continue
		}

		actual := sched.Next(getTime(c.time))
		/*
			astroT := astrotime.CalcSunset(actual, sched.Latitude, sched.Longitude)
			t.Logf("astro time:%v", astroT)
		*/
		expected := getTime(c.expected)
		if !actual.Equal(expected) {
			t.Errorf("%s, \"%s\": (expected) %v != %v (actual)", c.time, c.spec, expected, actual)
		}
	}

}
