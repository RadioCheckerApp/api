package request

import (
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"sort"
	"time"
)

type TracksWorker struct {
	dao     datalayer.TrackRecordDAO
	station string
}

func NewTracksWorker(dao datalayer.TrackRecordDAO, station string) (TracksWorker, error) {
	if dao == nil {
		return TracksWorker{}, errors.New("dao must not be nil")
	}
	if station == "" {
		return TracksWorker{}, errors.New("station must not be empty")
	}
	return TracksWorker{dao, station}, nil
}

func (worker TracksWorker) TopTracks(startDate, endDate time.Time) (model.CountedTracks, error) {
	trackRecords, err := worker.dao.GetTrackRecordsByStation(worker.station, startDate, endDate)
	if err != nil {
		return model.CountedTracks{}, err
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

	resultLimitIdx := findResultLimitIdx(orderedTracks)

	return model.CountedTracks{orderedTracks[:resultLimitIdx]}, nil
}

func (worker TracksWorker) AllTracks(startDate, endDate time.Time) (model.Tracks, error) {
	trackRecords, err := worker.dao.GetTrackRecordsByStation(worker.station, startDate, endDate)
	if err != nil {
		return model.Tracks{}, err
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

	return model.Tracks{tracks}, nil
}

func (worker TracksWorker) MostRecentTrackRecord() (model.TrackRecord, error) {
	trackRecord, err := worker.dao.GetMostRecentTrackRecordByStation(worker.station)
	if err != nil {
		return model.TrackRecord{}, err
	}
	return trackRecord, nil
}

func (worker TracksWorker) HandleRequest() (interface{}, error) {
	return worker.MostRecentTrackRecord()
}

func findResultLimitIdx(tracksOrderedDescendinglyByCounter []model.CountedTrack) int {
	if len(tracksOrderedDescendinglyByCounter) <= 3 || tracksOrderedDescendinglyByCounter[0].Counter <= 3 {
		return len(tracksOrderedDescendinglyByCounter)
	}

	prevCounter := tracksOrderedDescendinglyByCounter[0].Counter
	foundRanks, limitIdx := 1, 1
	for ; limitIdx < len(tracksOrderedDescendinglyByCounter) && foundRanks <= 3; limitIdx++ {
		if prevCounter != tracksOrderedDescendinglyByCounter[limitIdx].Counter {
			foundRanks++
			prevCounter = tracksOrderedDescendinglyByCounter[limitIdx].Counter
		}
	}

	if foundRanks > 3 {
		return limitIdx - 1
	}
	return limitIdx
}
