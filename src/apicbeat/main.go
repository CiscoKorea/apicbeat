package main

import (
	"os"
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"apicbeat/beater"
	"apicbeat/apic"
)

var Name = "apicbeat"

func main() {
	beatwork()
//	unitwork()
}

func beatwork() {
	if err := beat.Run(Name, "", beater.New()); err != nil {
		os.Exit(1)
	}
}

func unitwork() {
	s := apic.NewSession("10.72.86.21", "admin", "1234Qwer")
	fmt.Println(*s, "\n")
	
	{
		event, _ := apic.GetFvTenantEndPoints(s)
		for i, value := range event {
			fmt.Println(i, value)
		}
	}
	
	fmt.Println("")
	
	{
		event, _ := apic.GetFvTenantEndPointDNs(s)
		for i, value := range event {
			fmt.Println(i, value)
		}
	}
	
	fmt.Println("")
	
	{
		event, _ := apic.GetFvTenantHealths(s)
		for i, value := range event {
			fmt.Println(i, value)
		}
	}
	
	fmt.Println("")
	
	{
		event, _ := apic.GetFvTenantHealthCurs(s)
		for i, value := range event {
			fmt.Println(i, value)
		}
	}
	
	fmt.Println("")
	
	{
		event, _ := apic.GetFaultInfos(s)
		for i, value := range event {
			fmt.Println(i, value)
		}
	}
	
	
//	{ url := "class/fvAEPg.json?rsp-subtree-include=health"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
	
//	{ url := "class/fvAEPg.json?query-target=subtree&target-subtree-class=fvCEp&query-target-filter=wcard(fvCEp.dn, \"INNOTEK\")"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
	
//	{ url := "class/fvAEPg.json?rsp-subtree-include=faults"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "class/fvAEPg.json?query-target-filter=wcard(fvAEPg.dn, \"tn-01.SHINHAN.*\")&rsp-subtree-include=faults"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
	
//	{ url := "class/fvTenant.json"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	
//	{ url := "class/fvTenant.json?query-target-filter=eq(fvTenant.name, \"01.SHINHAN\")"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	
//	{ url := "mo/uni/tn-01.SHINHAN.json"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "mo/uni/tn-01.SHINHAN.json?query-target=subtree&target-subtree-class=fvAEPg&rsp-subtree-include=health"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "mo/uni/tn-01.SHINHAN.json?rsp-subtree-include=health,required"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "mo/uni/tn-01.SHINHAN.json?query-target=subtree&target-subtree-class=healthInst"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	
//	{ url := "class/fvTenant.json?query-target-filter=wcard(fvTenant.name, \"SHINHAN\")"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "class/fvTenant.json?query-target-filter=wcard(fvTenant.name, \"SHINHAN\")&rsp-subtree-include=health"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//
//	{ url := "class/fvAEPg.json?query-target-filter=wcard(fvAEPg.dn, \"SHINHAN\")"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "class/fvAEPg.json?query-target-filter=and(wcard(fvAEPg.dn, \"SHINHAN\"),wcard(fvAEPg.dn, \"WEB\"))"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "class/fvAEPg.json?query-target-filter=wcard(fvAEPg.dn, \"tn-.*01.SHINHAN.*\")"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	
//	{ url := "class/fvAEPg.json?query-target-filter=wcard(fvAEPg.dn, \"tn-01.SHINHAN.*\")&rsp-subtree-include=health"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	
//	{ url := "class/fvAEPg.json?query-target-filter=wcard(fvAEPg.dn, \"tn-01.SHINHAN.*\")&rsp-subtree-include=health&rsp-subtree-filter=gt(healthInst.cur, \"90\")"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "class/fvAEPg.json?query-target-filter=wcard(fvAEPg.dn, \"tn-01.SHINHAN.*\")&rsp-subtree-include=health&rsp-subtree-filter=bw(healthInst.cur, \"80\", \"90\")"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	{ url := "class/fvAEPg.json?rsp-subtree-include=health"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
//	
//	{ url := "class/fvTenant.json?query-target-filter=wcard(fvAEPg.dn, \"tn-01.SHINHAN.*\")&rsp-subtree-include=health"; data, _ := s.RawGet(url); fmt.Println(url, "\n", data, "\n") }
	
}
