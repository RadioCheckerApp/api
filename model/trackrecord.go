package model

type TrackRecord struct {
	StationId string `json:"stationId"`
	Timestamp int64  `json:"timestamp"`
	Track     Track  `json:"track"`
}
