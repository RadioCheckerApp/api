package model

type Station struct {
	ID          string `json:"stationId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}
