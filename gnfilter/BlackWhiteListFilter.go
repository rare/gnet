package gnfilter

import (
	"bufio"
	"container/list"
	"os"
	"strings"
	log "github.com/cihub/seelog"
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
				log.Tracef("black list rule: %s", rule)
				this.bl_rules.PushBack(rule)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Warnf("load black list error: (%v)", err)
			return err
		}
	}

	if wlf, err := os.Open(this.wl_file); err == nil {
		defer wlf.Close()

		scanner := bufio.NewScanner(wlf)
		for scanner.Scan() {
			rule := strings.Trim(scanner.Text(), " \t")
			if rule != "" {
				log.Tracef("white list rule: %s", rule)
				this.wl_rules.PushBack(rule)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Warnf("load white list error: (%v)", err)
			return err
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
