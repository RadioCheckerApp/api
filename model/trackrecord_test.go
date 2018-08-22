package model

import (
	"reflect"
	"testing"
	"time"
)

type testTrackRecordSuccess struct {
	input    *TrackRecord
	expected *TrackRecord
}

var timestamp = time.Now().Add(-13 * time.Minute).Unix()
var timestampFutureValid = time.Now().Add(29 * time.Minute).Unix()

var testsTrackRecordSuccess = []testTrackRecordSuccess{
	// stationId
	{
		&TrackRecord{"&nbsp;station-a", timestamp, "track", Track{"RHCP", "Californication"}},
		&TrackRecord{"station-a", timestamp, "track", Track{"rhcp", "californication"}},
	},
	{
		&TrackRecord{"AB", timestamp, "track", Track{"Felix Jaehn Feat. Jasmin Thompson", "Ain't Nobody (Loves Me Better)"}},
		&TrackRecord{"ab", timestamp, "track", Track{"felix jaehn feat. jasmin thompson", "ain't nobody (loves me better)"}},
	},
	{
		&TrackRecord{"hitradio-oe3", timestamp, "track", Track{"Axwell /\\ Ingrosso", "+++ The Shit +++"}},
		&TrackRecord{"hitradio-oe3", timestamp, "track", Track{"axwell /\\ ingrosso", "+++ the shit +++"}},
	},
	{
		&TrackRecord{"station24", timestamp, "TRACK", Track{"RHCP", "Californication"}},
		&TrackRecord{"station24", timestamp, "track", Track{"rhcp", "californication"}},
	},
	// timestamp
	{
		&TrackRecord{"hitradio-oe3", timestampFutureValid, "track", Track{"DOLLAR $IGN", "MØNE¥"}},
		&TrackRecord{"hitradio-oe3", timestampFutureValid, "track", Track{"dollar $ign", "møne¥"}},
	},
	// type
	{
		&TrackRecord{"station-a", timestamp, "TRACK", Track{"Nico &amp; Vinz feat. Kid Ink &amp; Bebe Rexha", "That's How You Know"}},
		&TrackRecord{"station-a", timestamp, "track", Track{"nico & vinz feat. kid ink & bebe rexha", "that's how you know"}},
	},
}

func TestTrackRecord_Sanitize_Success(t *testing.T) {
	for i, test := range testsTrackRecordSuccess {
		if err := test.input.Sanitize(); err != nil {
			t.Errorf("#%d (%q) Sanitize(): Expected no error, got `%s`.",
				i, test.input, err.Error())
			continue
		}
		if !reflect.DeepEqual(*test.input, *test.expected) {
			t.Errorf("#%d Sanitize(): Expected `%q`, got `%q`.",
				i, test.expected, test.input,
			)
		}
	}
}

var testsTrackRecordErr = []TrackRecord{
	// stationId
	{"station%20a", timestamp, "TRACK", Track{"RHCP", "Californication"}},
	{"A", timestamp, "TRACK", Track{"RHCP", "Californication"}},
	{"", timestamp, "TRACK", Track{"RHCP", "Californication"}},
	// timestamp
	{"station-a", time.Now().Add(time.Hour).Unix(), "TRACK", Track{"RHCP", "Californication"}},
	{"station-a", time.Now().Add(31 * time.Minute).Unix(), "TRACK", Track{"RHCP", "Californication"}},
	{"station-a", time.Now().AddDate(-10, 0, 0).Unix(), "TRACK", Track{"RHCP", "Californication"}},
	{"station-a", 5432955, "TRACK", Track{"RHCP", "Californication"}},
	// type
	{"station-a", timestamp, "", Track{"RHCP", "Californication"}},
	{"station-a", timestamp, "song", Track{"RHCP", "Californication"}},
	{"station-a", timestamp, "so--ng", Track{"RHCP", "Californication"}},
	// track
	{"station", timestamp, "track", Track{" ", ""}},
}

func TestTrackRecord_Sanitize_Err(t *testing.T) {
	for i, test := range testsTrackRecordErr {
		if err := test.Sanitize(); err == nil {
			t.Errorf("#%d Sanitize(): Expected error, got `%v` for Title: `%s`, Artist: `%s`.",
				i, err, test.Title, test.Artist)
		}
	}
}
