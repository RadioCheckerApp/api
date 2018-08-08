package request

type Worker interface {
	HandleRequest() (interface{}, error)
}
