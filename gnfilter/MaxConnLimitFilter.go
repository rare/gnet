package gnfilter

import (
	log "github.com/cihub/seelog"
)

var (
	DEFAULT_MAX_CONN	=	10000
)

type MaxConnFilter struct {
	cn		uint32					//cureent connection number
	limit	uint32					//max connection limit, 0 means no limit
}

func NewMaxConnFilter(limit uint32) *MaxConnFilter {
	return &MaxConnFilter {
		cn:		0,
		limit:	limit,
	}
}

func (this MaxConnFilter) CareEvent(evt EventType) bool {
	return evt == EVT_CONN_ACCEPTED || evt == EVT_CONN_CLOSING
}

func (this MaxConnFilter) DoFilter(evt EventType, obj interface{}) FilterResult {
	if evt == EVT_CONN_ACCEPTED {
		if this.limit != 0 && this.cn >= this.limit  {
			log.Warn("max conn limit reached")
			return FR_ABORT
		}
		this.cn++
	} else {	//EVT_CONN_CLOSING
		this.cn--
	}

	return FR_OK
}
