package model

import (
	"errors"
	"regexp"
)

type Track struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

func (track *Track) Sanitize() error {
	track.sanitizeArtist()
	track.sanitizeTitle()
	if track.Artist == "" {
		return errors.New("artist contains invalid data")
	}
	if track.Title == "" {
		return errors.New("title contains invalid data")
	}
	return nil
}

func (track *Track) sanitizeArtist() {
	track.Artist = cleanString(track.Artist)
	track.Artist = replaceURLs(track.Artist, "¯\\_(ツ)_/¯")
	track.Artist = addPeriod(track.Artist)
}

func (track *Track) sanitizeTitle() {
	track.Title = cleanString(track.Title)
	track.Title = replaceURLs(track.Title, "¯\\_(ツ)_/¯")
	track.Title = removeBranding(track.Title)
}

// filterURLs replaces HTTP URLs that may exist in the string.
func replaceURLs(str, substitute string) string {
	r := regexp.MustCompile(`http[s]?:\/\/[a-zA-Z0-9.:\/\?&-]+`)
	return r.ReplaceAllString(str, substitute)
}

// removeBranding discards the optional `(branding)` suffix some strings might contain.
func removeBranding(str string) string {
	r := regexp.MustCompile(` \((B|b)randing\)$`)
	return r.ReplaceAllString(str, "")
}

// addPeriod ensures that certain words are followed by a period (.), e. g. `feat` => `feat.`
func addPeriod(str string) string {
	// featuring (1)
	r := regexp.MustCompile(` feat `)
	str = r.ReplaceAllString(str, " feat. ")
	// featuring (2)
	r = regexp.MustCompile(` ft `)
	str = r.ReplaceAllString(str, " ft. ")
	// versus
	r = regexp.MustCompile(` vs `)
	return r.ReplaceAllString(str, " vs. ")
}

type CountedTrack struct {
	Counter int   `json:"times_played"`
	Track   Track `json:"track"`
}

type MatchedTrack struct {
	CountsByStation map[string]int `json:"plays_by_station"`
	Track           Track          `json:"track"`
}

type Tracks struct {
	Tracks []Track `json:"tracks"`
}

type CountedTracks struct {
	CountedTracks []CountedTrack `json:"tracks"`
}

type MatchedTracks struct {
	MatchedTracks []MatchedTrack `json:"tracks"`
}
