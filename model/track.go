package model

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

const dateFormat = "2006-01-02"

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
	Station   string    `json:"station"`
	StartDate time.Time `json:"omit"`
	EndDate   time.Time `json:"omit"`
	Tracks    []Track   `json:"tracks"`
}

type CountedTracks struct {
	Station       string         `json:"station"`
	StartDate     time.Time      `json:"omit"`
	EndDate       time.Time      `json:"omit"`
	CountedTracks []CountedTrack `json:"tracks"`
}

type MatchedTracks struct {
	StartDate     time.Time      `json:"omit"`
	EndDate       time.Time      `json:"omit"`
	MatchedTracks []MatchedTrack `json:"tracks"`
}

func (tracks Tracks) MarshalJSON() ([]byte, error) {
	type Alias Tracks
	if equalDate(tracks.StartDate, tracks.EndDate) {
		return json.Marshal(&struct {
			Date string `json:"date"`
			Alias
		}{
			Date:  tracks.StartDate.Format(dateFormat),
			Alias: (Alias)(tracks),
		})
	}

	return json.Marshal(&struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Alias
	}{
		StartDate: tracks.StartDate.Format(dateFormat),
		EndDate:   tracks.EndDate.Format(dateFormat),
		Alias:     (Alias)(tracks),
	})
}

func (tracks CountedTracks) MarshalJSON() ([]byte, error) {
	type Alias CountedTracks
	if equalDate(tracks.StartDate, tracks.EndDate) {
		return json.Marshal(&struct {
			Date string `json:"date"`
			Alias
		}{
			Date:  tracks.StartDate.Format(dateFormat),
			Alias: (Alias)(tracks),
		})
	}

	return json.Marshal(&struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Alias
	}{
		StartDate: tracks.StartDate.Format(dateFormat),
		EndDate:   tracks.EndDate.Format(dateFormat),
		Alias:     (Alias)(tracks),
	})
}

func (tracks MatchedTracks) MarshalJSON() ([]byte, error) {
	type Alias MatchedTracks
	if equalDate(tracks.StartDate, tracks.EndDate) {
		return json.Marshal(&struct {
			Date string `json:"date"`
			Alias
		}{
			Date:  tracks.StartDate.Format(dateFormat),
			Alias: (Alias)(tracks),
		})
	}

	return json.Marshal(&struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Alias
	}{
		StartDate: tracks.StartDate.Format(dateFormat),
		EndDate:   tracks.EndDate.Format(dateFormat),
		Alias:     (Alias)(tracks),
	})
}

func equalDate(d1, d2 time.Time) bool {
	return d1.Day() == d2.Day() &&
		d1.Month() == d2.Month() &&
		d1.Year() == d2.Year()
}
