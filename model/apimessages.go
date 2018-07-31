package model

type APIResponseMessage struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func NewAPIResponseMessage(data interface{}, err error) APIResponseMessage {
	if err != nil {
		return APIResponseMessage{false, nil, err.Error()}
	}
	return APIResponseMessage{true, data, ""}
}
