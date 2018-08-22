package request

import (
	"errors"
	"fmt"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
)

var stationsCache = make(map[string]bool)

type CreateTrackWorker struct {
	trackRecordDAO datalayer.TrackRecordDAO
	stationsDAO    datalayer.StationDAO
	trackRecord    model.TrackRecord
}

func NewCreateTrackWorker(trDAO datalayer.TrackRecordDAO, sDAO datalayer.StationDAO,
	trackRecord model.TrackRecord) (CreateTrackWorker, error) {
	if trDAO == nil || sDAO == nil {
		return CreateTrackWorker{}, errors.New("daos must not be nil")
	}
	return CreateTrackWorker{trDAO, sDAO, trackRecord}, nil
}

func (worker CreateTrackWorker) HandleRequest() (interface{}, error) {
	if err := worker.trackRecord.Sanitize(); err != nil {
		return nil, err
	}

	if len(stationsCache) == 0 {
		worker.populateStationsCache()
	}

	if !stationsCache[worker.trackRecord.StationId] {
		return nil, errors.New("invalid stationId provided")
	}

	if err := worker.trackRecordDAO.CreateTrackRecord(worker.trackRecord); err != nil {
		return nil, err
	}

	return fmt.Sprintf("track created: /stations/%s/tracks/%d",
		worker.trackRecord.StationId, worker.trackRecord.Timestamp), nil
}

func (worker CreateTrackWorker) populateStationsCache() {
	stations, _ := worker.stationsDAO.GetAll()
	for _, station := range stations {
		stationsCache[station.ID] = true
	}
}
