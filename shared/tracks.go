package shared

import (
	"encoding/json"
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"log"
	"sort"
	"strings"
	"time"
)

type Filter int

const (
	Err Filter = iota
	All
	Top
)

func Tracks(dao datalayer.TrackRecordDAO, pathParams, queryStringParams map[string]string) (
	string, error) {
	if len(pathParams) != 1 || len(queryStringParams) != 2 {
		return "", errors.New("invalid number of path/query string parameters provided")
	}

	station, err := getStation(pathParams)
	if err != nil {
		return "", err
	}

	filter, err := getFilter(queryStringParams)
	if err != nil {
		return "", err
	}

	date, err := getDate(queryStringParams)
	if err == nil {
		var jsonStr []byte
		var err error
		if filter == Top {
			jsonStr, err = json.Marshal(topTracksForDay(dao, station, date))
		} else if filter == All {
			jsonStr, err = json.Marshal(allTracksForDay(dao, station, date))
		}
		return string(jsonStr), err
	}

	var jsonStr []byte

	date, err = getFirstDayOfWeek(queryStringParams)
	if err != nil {
		return "", err
	}
	if filter == Top {
		jsonStr, err = json.Marshal(topTracksForWeek(dao, station, date))
	} else if filter == All {
		jsonStr, err = json.Marshal(allTracksForWeek(dao, station, date))
	}
	return string(jsonStr), err
}

func getStation(pathParams map[string]string) (string, error) {
	station, ok := pathParams["station"]
	if !ok || station == "" {
		return "", errors.New("path parameter `station` missing/invalid")
	}
	return strings.ToLower(station), nil
}

func getFilter(queryStringParams map[string]string) (Filter, error) {
	filterStr, _ := queryStringParams["filter"]
	switch strings.ToLower(filterStr) {
	case "top", "":
		return Top, nil
	case "all":
		return All, nil
	default:
		return Err, errors.New("invalid filter provided")
	}
}

func getDate(queryStringParams map[string]string) (time.Time, error) {
	dateStr, ok := queryStringParams["date"]
	if !ok {
		return time.Time{}, errors.New("date param not provided")
	}
	return time.ParseInLocation("2006-01-02", dateStr, getLocation())
}

func getFirstDayOfWeek(queryStringParams map[string]string) (time.Time, error) {
	dateStr, ok := queryStringParams["week"]
	if !ok {
		return time.Time{}, errors.New("week param not provided")
	}

	date, err := time.ParseInLocation("2006-01-02", dateStr, getLocation())
	if err != nil {
		return time.Time{}, err
	}

	return date.AddDate(0, 0, -normalizeWeekdayNumber(date)), nil
}

func getLocation() *time.Location {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Println("unable to load timezone location `Europe/Berlin`, using UTC instead")
		location, _ = time.LoadLocation("UTC")
	}
	return location
}

func normalizeWeekdayNumber(date time.Time) int {
	// Sunday = 0, ..., Saturday = 6
	usWeekdayNumber := date.Weekday()
	return int((usWeekdayNumber + 6) % 7)
}

func topTracksForDay(dao datalayer.TrackRecordDAO, station string, date time.Time) interface{} {
	startDate, endDate := calculateDayBoundaries(date)
	return topTracks(dao, station, startDate, endDate)
}

func topTracksForWeek(dao datalayer.TrackRecordDAO, station string, date time.Time) interface{} {
	startDate, endDate := calculateWeekBoundaries(date)
	return topTracks(dao, station, startDate, endDate)
}

func topTracks(dao datalayer.TrackRecordDAO, station string, startDate,
	endDate time.Time) []model.CountedTrack {
	trackRecords, err := dao.GetTrackRecords(station, startDate, endDate)
	if err != nil {
		log.Printf("topTracks(%q, %q, %q, %q): %q", dao, station, startDate, endDate, err)
		return []model.CountedTrack{}
	}

	groupedTracks := make(map[model.Track]int)
	for _, trackRecord := range trackRecords {
		groupedTracks[trackRecord.Track]++
	}

	orderedTracks := make([]model.CountedTrack, len(groupedTracks))
	i := 0
	for track, count := range groupedTracks {
		orderedTracks[i] = model.CountedTrack{Counter: count, Track: track}
		i++
	}

	sort.Slice(orderedTracks, func(i, j int) bool {
		return orderedTracks[i].Counter > orderedTracks[j].Counter
	})

	return orderedTracks
}

func allTracksForDay(dao datalayer.TrackRecordDAO, station string, date time.Time) interface{} {
	startDate, endDate := calculateDayBoundaries(date)
	return allTracks(dao, station, startDate, endDate)
}

func allTracksForWeek(dao datalayer.TrackRecordDAO, station string, date time.Time) interface{} {
	startDate, endDate := calculateWeekBoundaries(date)
	return allTracks(dao, station, startDate, endDate)
}

func calculateDayBoundaries(date time.Time) (time.Time, time.Time) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0,
		date.Location())
	endDate := startDate.AddDate(0, 0, 1).Add(-1 * time.Second)
	return startDate, endDate
}

func calculateWeekBoundaries(date time.Time) (time.Time, time.Time) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0,
		date.Location())
	endDate := startDate.AddDate(0, 0, 7).Add(-1 * time.Second)
	return startDate, endDate
}

func allTracks(dao datalayer.TrackRecordDAO, station string,
	startDate, endDate time.Time) []model.Track {
	trackRecords, err := dao.GetTrackRecords(station, startDate, endDate)
	if err != nil {
		log.Printf("allTracks(%q, %q, %q, %q): %q", dao, station, startDate, endDate, err)
		return []model.Track{}
	}

	distinctTracks := make(map[model.Track]bool, 0)
	for _, trackRecord := range trackRecords {
		distinctTracks[trackRecord.Track] = true
	}

	tracks := make([]model.Track, len(distinctTracks))
	i := 0
	for track := range distinctTracks {
		tracks[i] = track
		i++
	}

	return tracks
}
