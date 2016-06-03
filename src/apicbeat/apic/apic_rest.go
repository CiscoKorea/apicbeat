package apic

import (
	"os"
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"encoding/json"
)

type Dict map[string]interface {}
type List []interface {}

type Session struct {
	Addr string `json:"address"`
	User string `json:"user"`
	Pwd string `json:"pwd"`
	Client *http.Client `json:"client"`
	Cookies []*http.Cookie `json:"cookies"`
	Auth struct {
		TotalCount string `json:"totalCount"`
		Imdata []struct {
			AAALogin struct {
				Attributes struct {
					Token string `json:"token"`
					RefreshTimeoutSeconds string `json:"refreshTimeoutSeconds"`
					MaximumLifetimeSeconds string `json:"maximumLifetimeSeconds"`
					GuiIdleTimeoutSeconds string `json:"guiIdleTimeoutSeconds"`
					RestTimeoutSeconds string `json:"restTimeoutSeconds"`
				} `json:"attributes"`
			} `json:"aaaLogin"`
		} `json:"imdata"`
	} `json:"auth"`
}

func (s* Session) Login() {
	login_user := "{ \"aaaUser\" : { \"attributes\" : { \"name\" : \""+ s.User + "\", \"pwd\" : \"" + s.Pwd + "\"}}}"
	
	resp, err := s.Client.Post(s.Addr + "aaaLogin.json", "application/json", bytes.NewBufferString(login_user))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		fmt.Println("Login Failed")
		os.Exit(1)
	}

	body_raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	err = json.Unmarshal(body_raw, &s.Auth)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	s.Cookies = resp.Cookies()
}

func NewSession (addr string, user string, passwd string) *Session {
	s := &Session{
		Addr: "https://" + addr + "/api/",
		User: user,
		Pwd: passwd,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
	
	return s
}

func (s* Session) rest (method string, url string, data *bytes.Buffer) (List, error) {
	if data == nil { data = bytes.NewBufferString("") }
	req, err := http.NewRequest(method, s.Addr + url, data)
	for i := 0; i < len(s.Cookies); i++ { req.AddCookie(s.Cookies[i]) }
	
	resp, err := s.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 403 {
		s.Login()
		return s.rest(method, url, data)
	}
	if resp.StatusCode != 200 {
		fmt.Println("GET", url, "StatusCode :", resp.StatusCode)
		return nil, err
	}
	
	body_raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	
	body := Dict{}
	err = json.Unmarshal(body_raw, &body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	
	imdata := body["imdata"].([]interface{})
	
	return imdata, nil
}

func (s* Session) RawGet (url string) (string, error) {
	req, err := http.NewRequest("GET", s.Addr + url, bytes.NewBufferString(""))
	for i := 0; i < len(s.Cookies); i++ {
		req.AddCookie(s.Cookies[i])
	}
	
	resp, err := s.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 403 {
		s.Login()
		return s.RawGet(url)
	}
	if resp.StatusCode != 200 {
		fmt.Println("GET", url, "StatusCode :", resp.StatusCode)
		return "", err
	}
	
	body_raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	
	return string(body_raw), nil
}

func (s* Session) Get (url string) (List, error) {
	return s.rest("GET", url, bytes.NewBufferString(""))
}

func (s* Session) Post (url string, data string) (List, error) {
	return s.rest("POST", url, bytes.NewBufferString(data))
}

func (s* Session) Put (url string, data string) (List, error) {
	return s.rest("PUT", url, bytes.NewBufferString(data))
}

func (s* Session) Del (url string, data string) (List, error) {
	return s.rest("DELETE", url, bytes.NewBufferString(data))
}

	
