package gnfilter

import (
	"container/list"
	"errors"
)

type WhiteListFilter struct {
	cp		string				//conf file path
	rules	*List
}

func NewWhiteListFilter(conf_path string) *WhiteListFilter {
	return &{
		cp:	nil,
		rules: nil,
	}
}

func (this* WhiteListFilter) loadWhiteList() error {
}

func (this *WhiteListFilter) Init(conf_path string) error {
	this.cp = conf_path
	if err := loadWhiteList(); err != nil {
		return errors.New()
	}
}

func (this *WhiteListFilter) Filter(evt EventType, obj interface{}) {
	if EVT_CONN_ACCEPTED == evt {

	}

	return FR_OK
}
