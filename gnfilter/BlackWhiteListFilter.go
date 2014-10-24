package gnfilter

import (
	"container/list"
	"errors"
)

type BlackWhiteListFilter struct {
	cp			string					//conf file path
	bl_rules	*list.List				//black list rules
	wl_rules	*list.List				//black list rules
}

func NewBlackWhiteListFilter() *BlackWhiteListFilter {
	return &BlackWhiteListFilter {
		cp:	"",
		bl_rules: nil,
		wl_rules: nil,
	}
}

func (this* BlackWhiteListFilter) loadBlackWhiteList() error {
	return errors.New("load black white list error")
}

func (this *BlackWhiteListFilter) Init(conf_path string) error {
	this.cp = conf_path
	this.bl_rules = list.New()
	this.wl_rules = list.New()
	return this.loadBlackWhiteList()
}

func (this *BlackWhiteListFilter) CareEvent(evt EventType) bool {
	return evt == EVT_CONN_ACCEPTED
}

func (this *BlackWhiteListFilter) DoFilter(evt EventType, obj interface{}) FilterResult {

	return FR_OK
}
