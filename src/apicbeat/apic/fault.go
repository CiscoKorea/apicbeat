package apic

import (
	"fmt"
	"time"
	"github.com/elastic/beats/libbeat/common"
)

func GetFaultInfos (s *Session) ([]common.MapStr, error) {
	
	parser := func (data map[string]interface{}, class string) (common.MapStr, common.MapStr, error) {
		attr, dn := GetAttributesAndDN(data, class)
		desc := common.MapStr {
			"class" : class,
			"dn" : dn,
		}
		return desc, attr, nil
	}
	
	events := []common.MapStr{}
	datas, err := s.Get("class/faultInfo.json")
	if err != nil {
		fmt.Println(err)
		return events, err
	}
	
	timestamp := common.Time(time.Now())
	
	for _, data := range datas {
		if data.(map[string]interface{})["faultInst"] != nil {
			desc, attr, err := parser(data.(map[string]interface{}), "faultInst")
			if err != nil { continue }
			event := common.MapStr{
				"@timestamp": timestamp,
				"type": "faultinfo",
				"desc": desc,
				"faultinfo": attr,
			}
			events = append(events, event)
		} else if data.(map[string]interface{})["faultDelegate"] != nil {
			desc, attr, err := parser(data.(map[string]interface{}), "faultDelegate")
			if err != nil { continue }
			event := common.MapStr{
				"@timestamp": timestamp,
				"type": "faultinfo",
				"desc": desc,
				"faultinfo": attr,
			}
			events = append(events, event)
		} else {
			continue
		}
	}
	return events, nil
}
