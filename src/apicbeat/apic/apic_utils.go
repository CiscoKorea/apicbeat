package apic

import (
//	"bytes"
	"strings"
//	"net/http"
	"github.com/elastic/beats/libbeat/common"
)

func GetAttributesAndDN(data map[string]interface{}, class string) (common.MapStr, string) {
	attributes := data[class].(map[string]interface{})["attributes"].(map[string]interface{})
	return attributes, attributes["dn"].(string)
}

func GetAttributesAndRNs(data map[string]interface{}, class string) (common.MapStr, string, string, string, string) {
	attributes := data[class].(map[string]interface{})["attributes"].(map[string]interface{})
	dn := attributes["dn"].(string)
	split_dn := strings.Split(dn, "/")
	return attributes, dn, split_dn[1][3:], split_dn[2][3:], split_dn[3][4:]
}

func GetRNsFromDN(dn string) (string, string, string) {
	split_dn := strings.Split(dn, "/")
	return split_dn[1][3:], split_dn[2][3:], split_dn[3][4:]
}
