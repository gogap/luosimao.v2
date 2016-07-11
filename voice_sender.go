package luosimao

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type VoiceSender struct {
	sendUrl   string
	statusUrl string
	auth      Authorization
	Proxy     string

	clientPool sync.Pool
}

func NewVoiceSender(auth Authorization, proto ProtocalType, timeout time.Duration) *VoiceSender {
	sender := new(VoiceSender)
	switch proto {
	case JSON:
		{
			sender.sendUrl = VoiceServerURL + "verify.json"
			sender.statusUrl = VoiceServerURL + "status.json"
		}
	case XML:
		{
			sender.sendUrl = VoiceServerURL + "verify.xml"
			sender.statusUrl = VoiceServerURL + "status.xml"
		}
	}
	sender.auth = auth
	sender.clientPool = sync.Pool{
		New: func() interface{} {
			return &http.Client{
				Timeout: timeout,
			}
		},
	}

	return sender
}

func (p *VoiceSender) Send(req VoiceRequest) (resp Response, err error) {
	strCode := fmt.Sprintf("%04d", req.Code)

	params := url.Values{}
	params.Add("mobile", req.Mobile)
	params.Add("code", strCode)

	request, err := http.NewRequest("POST", p.sendUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", p.auth.BasicAuthorization())

	client := p.clientPool.Get().(*http.Client)
	defer p.clientPool.Put(client)

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return resp, errors.New(fmt.Sprintf("invalid http status code error: %s", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return
	}

	return
}

func (p *VoiceSender) Status() (status Status, err error) {
	params := url.Values{}
	request, err := http.NewRequest("POST", p.sendUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", p.auth.BasicAuthorization())

	client := p.clientPool.Get().(*http.Client)
	defer p.clientPool.Put(client)

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return status, errors.New(fmt.Sprintf("invalid http status code error: %s", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &status); err != nil {
		return
	}
	return
}
