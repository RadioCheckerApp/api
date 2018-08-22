package model

import (
	"errors"
	"html"
	"regexp"
	"strings"
	"time"
)

type TrackRecord struct {
	StationId string `json:"stationId"`
	Timestamp int64  `json:"airtime"`
	Type      string `json:"type"`
	Track
}

func (record *TrackRecord) Sanitize() error {
	if err := record.sanitizeStationId(); err != nil {
		return err
	}
	if err := record.sanitizeTimestamp(); err != nil {
		return err
	}
	if err := record.sanitizeType(); err != nil {
		return err
	}
	return record.Track.Sanitize()
}

func (record *TrackRecord) sanitizeStationId() error {
	record.StationId = cleanString(record.StationId)
	r := regexp.MustCompile(`^[a-z][a-z0-9-]+$`)
	if !r.MatchString(record.StationId) {
		return errors.New("stationId contains invalid format")
	}
	return nil
}

func (record *TrackRecord) sanitizeTimestamp() error {
	airtime := time.Unix(record.Timestamp, 0)
	// tracks are allowed to lie max. 30min in the future, since some APIs also return the tracks of
	// the near future -- and we won't trash them, right?
	if airtime.After(time.Now().Add(30 * time.Minute)) {
		return errors.New("timestamp lies in the future")
	}
	if airtime.Before(time.Date(2016, 1, 1, 0, 0, 0, 0, airtime.Location())) {
		return errors.New("timestamp is older than RadioChecker")
	}
	return nil
}

func (record *TrackRecord) sanitizeType() error {
	record.Type = cleanString(record.Type)
	if record.Type != "track" {
		return errors.New("type is not `track`")
	}
	return nil
}

func cleanString(str string) string {
	str = strings.ToLower(str)     // html.UnescapeString() requires lowercase input
	str = html.UnescapeString(str) // unescapes HTML special characters to their original form
	str = discardWhitespaces(str)
	return str
}

// dicardWhitespaces removes leading/trailing whitespaces as well as redundant whitespaces
// between words.
func discardWhitespaces(str string) string {
	str = strings.TrimSpace(str)
	return strings.Join(strings.Fields(str), " ")
}
