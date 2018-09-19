package request

import (
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"strings"
	"time"
)

const queryStrKeywordsSeparator = "+"

type groupedTracksContainer map[model.Track]map[string]int

type SearchWorker struct {
	dao      datalayer.TrackRecordDAO
	keywords []string
}

func NewSearchWorker(dao datalayer.TrackRecordDAO, query string) (SearchWorker, error) {
	if dao == nil {
		return SearchWorker{}, errors.New("dao must not be nil")
	}
	if query == "" {
		return SearchWorker{}, errors.New("query must not be empty")
	}
	lowercaseQuery := strings.ToLower(query)
	return SearchWorker{dao, strings.Split(lowercaseQuery, queryStrKeywordsSeparator)}, nil
}

func (worker SearchWorker) Search(startDate, endDate time.Time) (model.MatchedTracks, error) {
	trackRecords, err := worker.dao.GetTrackRecords(startDate, endDate)
	if err != nil {
		return model.MatchedTracks{}, err
	}

	matchedTrackRecords := worker.findMatchingTrackRecords(trackRecords)
	if len(matchedTrackRecords) == 0 {
		return model.MatchedTracks{
			model.MetaInfo{startDate, endDate},
			[]model.MatchedTrack{},
		}, nil
	}

	stationIDs := extractStationIDs(matchedTrackRecords)

	groupedTracks := make(groupedTracksContainer)

	for _, trackRecord := range matchedTrackRecords {
		if _, ok := groupedTracks[trackRecord.Track]; !ok {
			groupedTracks[trackRecord.Track] = newStationsMap(stationIDs)
		}
		groupedTracks[trackRecord.Track][trackRecord.StationId]++
	}

	return model.MatchedTracks{
		model.MetaInfo{startDate, endDate},
		buildResultStructure(groupedTracks),
	}, nil
}

func (worker SearchWorker) findMatchingTrackRecords(trackRecords []model.TrackRecord) []model.
	TrackRecord {
	matchedTrackRecords := make([]model.TrackRecord, 0)

	for _, trackRecord := range trackRecords {
		if worker.trackRecordMatchesQuery(trackRecord) {
			matchedTrackRecords = append(matchedTrackRecords, trackRecord)
		}
	}

	return matchedTrackRecords
}

func (worker SearchWorker) trackRecordMatchesQuery(trackRecord model.TrackRecord) bool {
	title := strings.ToLower(trackRecord.Title)
	artist := strings.ToLower(trackRecord.Artist)
	for _, keyword := range worker.keywords {
		if strings.Contains(title, keyword) || strings.Contains(artist, keyword) {
			return true
		}
	}
	return false
}

func extractStationIDs(trackRecords []model.TrackRecord) []string {
	groupedStationIDs := make(map[string]bool)
	for _, trackRecord := range trackRecords {
		groupedStationIDs[trackRecord.StationId] = true
	}

	stationIDs := make([]string, len(groupedStationIDs))
	i := 0
	for key := range groupedStationIDs {
		stationIDs[i] = key
		i++
	}
	return stationIDs
}

func newStationsMap(stationIDs []string) map[string]int {
	stationsMap := make(map[string]int)
	for _, stationID := range stationIDs {
		stationsMap[stationID] = 0
	}
	return stationsMap
}

func buildResultStructure(groupedTracks groupedTracksContainer) []model.MatchedTrack {
	matchedTracks := make([]model.MatchedTrack, len(groupedTracks))
	i := 0
	for track, countsByStation := range groupedTracks {
		matchedTracks[i] = model.MatchedTrack{countsByStation, track}
		i++
	}
	return matchedTracks
}
