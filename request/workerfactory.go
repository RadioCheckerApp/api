package request

import (
	"encoding/json"
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Filter int

const (
	Err Filter = iota
	All
	Top
)

const (
	queryStrDateParam      = "date"
	queryStrWeekParam      = "week"
	queryStrFilterParam    = "filter"
	queryStrStationParam   = "station"
	queryStrQueryParam     = "q"
	queryStrTimestampParam = "timestamp"
)

func CreateMetaWorker() Worker {
	return MetaWorker{}
}

func CreateStationsWorker(dao datalayer.StationDAO) (Worker, error) {
	return NewStationsWorker(dao)
}

func CreateTracksWorker(dao datalayer.TrackRecordDAO, pathParams,
	queryStringParams map[string]string) (Worker, error) {
	station, err := getStation(pathParams)
	if err != nil {
		return nil, err
	}

	filter, err := getFilter(queryStringParams)
	if err != nil {
		return nil, err
	}

	if formattedDateStr, ok := queryStringParams[queryStrDateParam]; ok {
		date, err := createDate(formattedDateStr)
		if err != nil {
			return nil, err
		}
		return NewDayTracksWorker(dao, station, date, filter)
	}

	if formattedDateStr, ok := queryStringParams[queryStrWeekParam]; ok {
		date, err := createDate(formattedDateStr)
		if err != nil {
			return nil, err
		}
		return NewWeekTracksWorker(dao, station, date, filter)
	}

	return nil, errors.New("invalid/insufficient parameter(s) provided")
}

func getStation(pathParams map[string]string) (string, error) {
	station, ok := pathParams[queryStrStationParam]
	if !ok || station == "" {
		return "", errors.New("path parameter `station` missing/invalid")
	}
	return strings.ToLower(station), nil
}

func getFilter(queryStringParams map[string]string) (Filter, error) {
	filterStr, _ := queryStringParams[queryStrFilterParam]
	switch strings.ToLower(filterStr) {
	case "top", "":
		return Top, nil
	case "all":
		return All, nil
	default:
		return Err, errors.New("invalid filter provided")
	}
}

func createDate(formattedDateStr string) (time.Time, error) {
	date, err := time.ParseInLocation("2006-01-02", formattedDateStr, getLocation())
	if err != nil {
		return time.Time{}, errors.New("invalid date format provided")
	}
	return date, err
}

func CreateSearchWorker(dao datalayer.TrackRecordDAO, queryStringParams map[string]string) (Worker, error) {
	query, err := getQuery(queryStringParams)
	if err != nil {
		return nil, err
	}

	if formattedDateStr, ok := queryStringParams[queryStrDateParam]; ok {
		date, err := createDate(formattedDateStr)
		if err != nil {
			return nil, err
		}
		return NewDaySearchWorker(dao, query, date)
	}

	if formattedDateStr, ok := queryStringParams[queryStrWeekParam]; ok {
		date, err := createDate(formattedDateStr)
		if err != nil {
			return nil, err
		}
		return NewWeekSearchWorker(dao, query, date)
	}

	return nil, errors.New("invalid/insufficient parameter(s) provided")
}

func getQuery(queryStringParams map[string]string) (string, error) {
	query, ok := queryStringParams[queryStrQueryParam]
	if !ok || query == "" {
		return "", errors.New("no query provided")
	}
	return query, nil
}

func CreateCreateTrackWorker(trDAO datalayer.TrackRecordDAO, sDAO datalayer.StationDAO,
	pathParams map[string]string, body []byte) (Worker, error) {
	station, err := getStation(pathParams)
	if err != nil {
		return nil, err
	}

	timestamp, err := getTimestamp(pathParams)
	if err != nil {
		return nil, err
	}

	track, err := getTrack(body)
	if err != nil {
		return nil, err
	}

	trackRecord := model.TrackRecord{station, timestamp, "track", track}
	return NewCreateTrackWorker(trDAO, sDAO, trackRecord)
}

func getTimestamp(pathParams map[string]string) (int64, error) {
	timestamp, ok := pathParams[queryStrTimestampParam]
	if !ok || timestamp == "" {
		return 0, errors.New("no timestamp provided")
	}
	return strconv.ParseInt(timestamp, 10, 64)
}

func getTrack(body []byte) (model.Track, error) {
	// TODO: Implement Unmarshaller interface in model.Track
	var track model.Track
	if err := json.Unmarshal(body, &track); err != nil {
		return model.Track{}, errors.New("request body contains invalid JSON")
	}
	if reflect.DeepEqual(track, model.Track{}) {
		return model.Track{}, errors.New("request body contains invalid data")
	}
	return track, nil
}
