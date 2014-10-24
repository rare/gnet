package gnfilter

import (
	"net"
)

var (
	DEFAULT_MAX_CONN	=	10000
)

type MaxConnFilter struct {
	cn		uint32		//cureent connection number
	limit	uint32		//max connection limit, 0 means no limit
}

func NewMaxConnFilter(limit uint32) *MaxConnFilter {
	return &{
		cn:		0,
		limit:	limit,
	}
}

func (this *MaxConnFilter) Filter(evt EventType, obj interface{}) FilterResult {
	if conn, ok := obj.(*net.TCPConn); ok {
		if evt == EVT_CONN_ACCEPTED {
			if this.limit != 0 && this.cn >= this.limit  {
				return FR_ABORT
			}
			this.cn++
		} else if EVT_CONN_CLOSING {
			this.cn--
		} 


		return FR_OK
	} else {
		//bad paramater type
		//TODO
		return FR_ABORT
	}
}
