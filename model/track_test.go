package model

import (
	"testing"
)

type testTrackSuccess struct {
	input    *Track
	expected *Track
}

var testsTrackSuccess = []testTrackSuccess{
	// URL tests
	{
		&Track{Title: "Hall Of Fame", Artist: "The Script / https://t.co/VpFYIHO9BJ"},
		&Track{Title: "hall of fame", Artist: "the script / ¯\\_(ツ)_/¯"},
	},
	{
		&Track{Title: "This Is Love", Artist: "https://t.co/VpFYII5L0j / Eva Simons"},
		&Track{Title: "this is love", Artist: "¯\\_(ツ)_/¯ / eva simons"},
	},
	{
		&Track{Title: "Thatpower", Artist: "https://t.co/VpFYII5L0j"},
		&Track{Title: "thatpower", Artist: "¯\\_(ツ)_/¯"},
	},
	{
		&Track{Title: " https://t.co/VpFYII5L0j ", Artist: "https://t.co/VpFYII5L0j"},
		&Track{Title: "¯\\_(ツ)_/¯", Artist: "¯\\_(ツ)_/¯"},
	},
	{
		&Track{Title: " https://t.co/VpFYII5L0j https://t.co/VpFYII5L0j ", Artist: "https://t.co/VpFYII5L0j"},
		&Track{Title: "¯\\_(ツ)_/¯ ¯\\_(ツ)_/¯", Artist: "¯\\_(ツ)_/¯"},
	},
	// HTML special chars tests
	{
		&Track{Title: "Súbeme La Radio", Artist: "Enrique Iglesias feat. Descemer Bueno, Zion &amp; Lenn"},
		&Track{Title: "súbeme la radio", Artist: "enrique iglesias feat. descemer bueno, zion & lenn"},
	},
	{
		&Track{Title: "Despacito (Remix)", Artist: "Luis Fonsi &amp; Daddy Yankee feat. Justin Bieber"},
		&Track{Title: "despacito (remix)", Artist: "luis fonsi & daddy yankee feat. justin bieber"},
	},
	{
		&Track{Title: "That's How You Know", Artist: "Nico &amp; Vinz feat. Kid Ink &amp; Bebe Rexha"},
		&Track{Title: "that's how you know", Artist: "nico & vinz feat. kid ink & bebe rexha"},
	},
	// whitespaces tests
	{
		&Track{Title: "That's  How You Know", Artist: "Nico &amp; Vinz feat. Kid    Ink &amp; Bebe Rexha"},
		&Track{Title: "that's how you know", Artist: "nico & vinz feat. kid ink & bebe rexha"},
	},
	{
		&Track{Title: "This Is   https://t.co/VpFYII5L0j   ", Artist: " https://t.co/VpFYII5L0j    / Eva  &amp;   Simons     "},
		&Track{Title: "this is ¯\\_(ツ)_/¯", Artist: "¯\\_(ツ)_/¯ / eva & simons"},
	},
	// `(Branding)` tests
	{
		&Track{Title: "That's  It", Artist: "Nico &amp; Vinz feat. Kid    Ink &amp; Bebe Rexha (Branding)"},
		&Track{Title: "that's it", Artist: "nico & vinz feat. kid ink & bebe rexha (branding)"},
	},
	{
		&Track{Title: "Wild Thoughts (Branding)", Artist: "DJ Khaled feat. Rihanna &amp; Bryson Tiller"},
		&Track{Title: "wild thoughts", Artist: "dj khaled feat. rihanna & bryson tiller"},
	},
	{
		&Track{Title: "Wild Thoughts (Branding) ", Artist: " https://t.co/VpFYII5L0j    / Eva  &amp;   Simons     "},
		&Track{Title: "wild thoughts", Artist: "¯\\_(ツ)_/¯ / eva & simons"},
	},
	{
		&Track{Title: " https://t.co/VpFYII5L0j   (Branding) ", Artist: " https://t.co/VpFYII5L0j    / Eva  &amp;   Simons     "},
		&Track{Title: "¯\\_(ツ)_/¯", Artist: "¯\\_(ツ)_/¯ / eva & simons"},
	},
	// lowercase tests
	{
		&Track{Title: "&NBSP; SÙBEME  LA RADIO", Artist: "ENRIQUE IGLESIAS Feat. DESCEMER BUENO, Zion &amp; Lenn"},
		&Track{Title: "sùbeme la radio", Artist: "enrique iglesias feat. descemer bueno, zion & lenn"},
	},
	{
		&Track{Title: "ROLLIN'", Artist: "SUM FREAKIN' ARTS'Y"},
		&Track{Title: "rollin'", Artist: "sum freakin' arts'y"},
	},
	{
		&Track{Title: "WHOOP!'", Artist: "TY DOLLA $IGN'"},
		&Track{Title: "whoop!'", Artist: "ty dolla $ign'"},
	},
	{
		&Track{Title: "U", Artist: "ADELE "},
		&Track{Title: "u", Artist: "adele"},
	},
	{
		&Track{Title: "SUMMER OF '69 (LIVE)", Artist: "BRYAN ADAMS "},
		&Track{Title: "summer of '69 (live)", Artist: "bryan adams"},
	},
	{
		&Track{Title: "STRONGER (WHAT DOESN'T KILL YOU)", Artist: "KELLY CLARKSON"},
		&Track{Title: "stronger (what doesn't kill you)", Artist: "kelly clarkson"},
	},
	{
		&Track{Title: "Y.M.C.A", Artist: "Village People"},
		&Track{Title: "y.m.c.a", Artist: "village people"},
	},
	{
		&Track{Title: "#THATPOWER", Artist: "WILL.I.AM"},
		&Track{Title: "#thatpower", Artist: "will.i.am"},
	},
	// period tests
	{
		&Track{Title: "You Don't Know Me", Artist: " JAX JONES FEAT  RAYE"},
		&Track{Title: "you don't know me", Artist: "jax jones feat. raye"},
	},
	{
		&Track{Title: "GREAT SPIRIT", Artist: "Armin Van Buuren Vs Vini Vici Feat Hilight &Tribe"},
		&Track{Title: "great spirit", Artist: "armin van buuren vs. vini vici feat. hilight &tribe"},
	},
	{
		&Track{Title: "Please Tell Rosie", Artist: "Alle Farben   ft Younotus"},
		&Track{Title: "please tell rosie", Artist: "alle farben ft. younotus"},
	},
	{
		&Track{Title: "Ain't Nobody (Loves Me Better)", Artist: "Felix Jaehn Feat. Jasmin Thompson"},
		&Track{Title: "ain't nobody (loves me better)", Artist: "felix jaehn feat. jasmin thompson"},
	},
}

func TestTrack_Sanitize_Success(t *testing.T) {
	for i, test := range testsTrackSuccess {
		oT := test.input.Title
		oA := test.input.Artist
		if err := test.input.Sanitize(); err != nil {
			t.Errorf("#%d Sanitize(): Expected no error, got `%s` for Title: `%s`, Artist: `%s`.",
				i, err.Error(), oT, oA)
			continue
		}
		if test.input.Title != test.expected.Title {
			t.Errorf("#%d Sanitize(): Title: Expected `%s` for `%s`, got `%s`.",
				i, test.expected.Title, oT, test.input.Title,
			)
		}
		if test.input.Artist != test.expected.Artist {
			t.Errorf("#%d Sanitize(): Artist: Expected `%s` for `%s`, got `%s`.",
				i, test.expected.Artist, oA, test.input.Artist,
			)
		}
	}
}

var testsErr = []Track{
	{Title: "", Artist: ""},
	{Title: "   ", Artist: "  "},
	{Title: "Tabs", Artist: "		"},
	{Title: " 　 ", Artist: "Different whitespace chars"},
}

func TestTrack_Sanitize_Err(t *testing.T) {
	for i, test := range testsErr {
		if err := test.Sanitize(); err == nil {
			t.Errorf("#%d Sanitize(): Expected no error, got `%v` for Title: `%s`, Artist: `%s`.",
				i, err, test.Title, test.Artist)
		}
	}
}
