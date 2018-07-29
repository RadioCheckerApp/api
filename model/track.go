package model

type Track struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type CountedTrack struct {
	Counter int   `json:"times_played"`
	Track   Track `json:"track"`
}
