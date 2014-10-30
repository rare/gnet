package gnfilter

import (
	"bufio"
	"container/list"
	"errors"
	"os"
	"strings"
)

type BlackWhiteListFilter struct {
	bl_file		string					//black list file
	wl_file		string					//white list file
	bl_rules	*list.List				//black list rules
	wl_rules	*list.List				//black list rules
}

func NewBlackWhiteListFilter() *BlackWhiteListFilter {
	return &BlackWhiteListFilter {
		bl_file: "",
		wl_file: "",
		bl_rules: nil,
		wl_rules: nil,
	}
}

func (this BlackWhiteListFilter) loadBlackWhiteList() error {
	//load black list
	if blf, err := os.Open(this.bl_file); err == nil {
		defer blf.Close()

		scanner := bufio.NewScanner(blf)
		for scanner.Scan() {
			rule := strings.Trim(scanner.Text(), " \t")
			if rule != "" {
				this.bl_rules.PushBack(rule)
			}
		}
		if err := scanner.Err(); err != nil {
			//TODO
			//trace
			return errors.New("load black list error")
		}
	}

	if wlf, err := os.Open(this.wl_file); err == nil {
		defer wlf.Close()

		scanner := bufio.NewScanner(wlf)
		for scanner.Scan() {
			rule := strings.Trim(scanner.Text(), " \t")
			if rule != "" {
				this.wl_rules.PushBack(rule)
			}
		}
		if err := scanner.Err(); err != nil {
			//TODO
			//trace
			return errors.New("load white list error")
		}
	}

	return nil
}

func (this BlackWhiteListFilter) Init(blf string, wlf string) error {
	this.bl_file = blf
	this.wl_file = wlf
	return this.loadBlackWhiteList()
}

func (this BlackWhiteListFilter) CareEvent(evt EventType) bool {
	return evt == EVT_CONN_ACCEPTED
}

func (this BlackWhiteListFilter) DoFilter(evt EventType, obj interface{}) FilterResult {
	//TODO

	return FR_OK
}
