package apic

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"github.com/elastic/beats/libbeat/common"
)

func GetFvTenantNames (s *Session)([]string, error) {
	datas, err := s.Get("class/fvTenant.json")
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}
	var names []string
	for _, data := range datas {
		names = append(names, data.(map[string]interface{})["fvTenant"].(map[string]interface{})["attributes"].(map[string]interface{})["name"].(string))
	}
	return names, nil
}

func GetFvTenantEndPoints (s * Session) ([]common.MapStr, error) {
	parser := func (data map[string]interface{}) (common.MapStr, common.MapStr, error) {
		attr, dn, tenant, app, epg := GetAttributesAndRNs(data, "fvCEp")
		if tenant == "infra" { return nil, nil, fmt.Errorf("infra data") }
		desc := common.MapStr {
			"class" : "fvCEp",
			"tenant" : tenant,
			"app" : app,
			"epg" : epg,
			"dn" : dn,
		}
		return desc, attr, nil
	}
	
	events := []common.MapStr{}
	datas, err := s.Get("class/fvAEPg.json?query-target=subtree&target-subtree-class=fvCEp")
	if err != nil {
		fmt.Println(err)
		return events, err
	}
	
	timestamp := common.Time(time.Now())
	
	for _, data := range datas {
		desc, attr, err := parser(data.(map[string]interface{}))
		if err != nil { continue }
		event := common.MapStr{
			"@timestamp": timestamp,
			"type": "tenant_endpoint",
			"desc": desc,
			"tenant_endpoint": attr,
		}
		events = append(events, event)
	}
	return events, nil
}

func GetFvTenantEndPointDNs (s * Session) ([]common.MapStr, error) {
	parser := func (data map[string]interface{}) (common.MapStr, string, error) {
		_, dn, tenant, app, epg := GetAttributesAndRNs(data, "fvCEp")
		if tenant == "infra" { return nil, "", fmt.Errorf("infra data") }
		desc := common.MapStr {
			"class" : "fvCEp",
			"tenant" : tenant,
			"app" : app,
			"epg" : epg,
			"dn" : dn,
		}
		return desc, dn, nil
	}
	
	events := []common.MapStr{}
	datas, err := s.Get("class/fvAEPg.json?query-target=subtree&target-subtree-class=fvCEp")
	if err != nil {
		fmt.Println(err)
		return events, err
	}
	
	timestamp := common.Time(time.Now())
	
	for _, data := range datas {
		desc, dn, err := parser(data.(map[string]interface{}))
		if err != nil { continue }
		event := common.MapStr{
			"@timestamp": timestamp,
			"type": "tenant_endpoint_dn",
			"desc": desc,
			"tenant_endpoint_dn": dn,
		}
		events = append(events, event)
	}
	return events, nil
}

func GetFvTenantHealths (s *Session) ([]common.MapStr, error) {
	
	parser := func (data map[string]interface{}) ([]common.MapStr, []common.MapStr, error) {
		attributes := data["fvAEPg"].(map[string]interface{})["attributes"].(map[string]interface{})
		dn := attributes["dn"].(string)
		
		if strings.Contains(dn, "uni/tn-") {
			tenant, app, epg := GetRNsFromDN(dn)
			if tenant == "infra" { return nil, nil, fmt.Errorf("infra data") } 
			
			descs := []common.MapStr{}
			attrs := []common.MapStr{}
			
			children := data["fvAEPg"].(map[string]interface{})["children"].([]interface{})
			for _, child := range children {
				var attr map[string]interface{}
				var class string
				var rn string
				if child.(map[string]interface{})["healthInst"] != nil {
					attr = child.(map[string]interface{})["healthInst"].(map[string]interface{})["attributes"].(map[string]interface{})
					class = "healthInst"
					rn = child.(map[string]interface{})["healthInst"].(map[string]interface{})["attributes"].(map[string]interface{})["rn"].(string)
				} else if child.(map[string]interface{})["healthNodeInst"] != nil {
					attr = child.(map[string]interface{})["healthNodeInst"].(map[string]interface{})["attributes"].(map[string]interface{})
					class = "healthNodeInst"
					rn = child.(map[string]interface{})["healthNodeInst"].(map[string]interface{})["attributes"].(map[string]interface{})["rn"].(string)
				} else { continue }
				
				desc := common.MapStr {
					"class" : class,
					"tenant" : tenant,
					"app" : app,
					"epg" : epg,
					"dn" : dn,
					"rn" : rn,
				}
				descs = append(descs, desc)
				attrs = append(attrs, attr)
			}
			return descs, attrs, nil
		}
		return nil, nil, fmt.Errorf("non uni/tn data") 
	}
	
	events := []common.MapStr{}
	datas, err := s.Get("class/fvAEPg.json?rsp-subtree-include=health")
	if err != nil {
		fmt.Println(err)
		return events, err
	}
	
	timestamp := common.Time(time.Now())
	
	for _, data := range datas {
		descs, attrs, err := parser(data.(map[string]interface{}))
		if err != nil { continue }
		for idx, desc := range descs {
			event := common.MapStr{
				"@timestamp": timestamp,
				"type": "tenant_health",
				"desc": desc,
				"tenant_health": attrs[idx],
			}
			events = append(events, event)
		}
	}
	return events, nil
}

func GetFvTenantHealthCurs (s* Session) ([]common.MapStr, error) {
	parser := func (data map[string]interface{}) ([]common.MapStr, []int, error) {
		attributes := data["fvAEPg"].(map[string]interface{})["attributes"].(map[string]interface{})
		dn := attributes["dn"].(string)
		
		if strings.Contains(dn, "uni/tn-") {
			tenant, app, epg := GetRNsFromDN(dn)
			if tenant == "infra" { return nil, nil, fmt.Errorf("infra data") } 
		
			descs := []common.MapStr{}
			healths := []int{}
		
			children := data["fvAEPg"].(map[string]interface{})["children"].([]interface{})
			for _, child := range children {
				var health int
				var class string
				var rn string
				if child.(map[string]interface{})["healthInst"] != nil {
					health, _ = strconv.Atoi(child.(map[string]interface{})["healthInst"].(map[string]interface{})["attributes"].(map[string]interface{})["cur"].(string));
					class = "healthInst"
					rn = child.(map[string]interface{})["healthInst"].(map[string]interface{})["attributes"].(map[string]interface{})["rn"].(string)
				} else if child.(map[string]interface{})["healthNodeInst"] != nil {
					health, _ = strconv.Atoi(child.(map[string]interface{})["healthNodeInst"].(map[string]interface{})["attributes"].(map[string]interface{})["cur"].(string));
					class = "healthNodeInst"
					rn = child.(map[string]interface{})["healthNodeInst"].(map[string]interface{})["attributes"].(map[string]interface{})["rn"].(string)
				} else { continue }
				
				desc := common.MapStr {
					"class" : class,
					"tenant" : tenant,
					"app" : app,
					"epg" : epg,
					"dn" : dn,
					"rn" : rn,
				}
				descs = append(descs, desc)
				healths = append(healths, health)
			}
			return descs, healths, nil
		}
		return nil, nil, fmt.Errorf("non uni/tn data") 
	}
	
	events := []common.MapStr{}
	datas, err := s.Get("class/fvAEPg.json?rsp-subtree-include=health")
	if err != nil {
		fmt.Println(err)
		return events, err
	}
	
	timestamp := common.Time(time.Now())
	
	for _, data := range datas {
		descs, healths, err := parser(data.(map[string]interface{}))
		if err != nil { continue }
		for idx, desc := range descs {
			event := common.MapStr{
				"@timestamp": timestamp,
				"type": "tenant_health_cur",
				"desc": desc,
				"tenant_health_cur": healths[idx],
			}
			events = append(events, event)
		}
	}
	return events, nil
}