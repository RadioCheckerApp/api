package model

type Track struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
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
