package gnfilter

type EventType int

const (
	EVT_CONN_ACCEPTED	EventType = iota
	EVT_CONN_CLOSSING

	EVT_REQ_BEFORE
	EVT_REQ_AFTER
)

type FilterResult int

const (
	FR_OK				FilterResult = iota		//ok, but will continue processing with other filters
	FR_END										//ok, stop processing with other filters
	FR_ABORT									//fail 
)

type Filter interface {
	Filter(evt EventType, obj interface{}) FilterResult
}
