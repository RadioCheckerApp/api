package model

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"
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
	MetaInfo
	Tracks []Track `json:"tracks"`
}

type CountedTracks struct {
	MetaInfo
	CountedTracks []CountedTrack `json:"tracks"`
}

type MatchedTracks struct {
	MetaInfo
	MatchedTracks []MatchedTrack `json:"tracks"`
}

type MetaInfo struct {
	StartDate time.Time
	EndDate   time.Time
}

func (m MetaInfo) MarshalJSON() ([]byte, error) {
	dateFormat := "02.01.2006"
	equalDay := m.StartDate.Day() == m.EndDate.Day() &&
		m.StartDate.Month() == m.EndDate.Month() &&
		m.StartDate.Year() == m.EndDate.Year()
	if equalDay {
		return marshalSingleDayInterval(m, dateFormat)
	}
	return marshalMultiDayInterval(m, dateFormat)
}

func marshalSingleDayInterval(m MetaInfo, dateFormat string) ([]byte, error) {
	return json.Marshal(&struct {
		Date string `json:"date"`
	}{
		Date: m.StartDate.Format(dateFormat),
	})
}

func marshalMultiDayInterval(m MetaInfo, dateFormat string) ([]byte, error) {
	return json.Marshal(&struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{
		StartDate: m.StartDate.Format(dateFormat),
		EndDate:   m.EndDate.Format(dateFormat),
	})
}
