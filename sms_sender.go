package luosimao

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SMSSender struct {
	sendUrl      string
	sendBatchUrl string
	statusUrl    string
	auth         Authorization
	Proxy        string

	clientPool sync.Pool
}

func NewSMSSender(auth Authorization, proto ProtocalType, timeout time.Duration) *SMSSender {
	sender := new(SMSSender)
	switch proto {
	case JSON:
		{
			sender.sendUrl = SMSServerURL + "send.json"
			sender.sendBatchUrl = SMSServerURL + "send_batch.json"
			sender.statusUrl = SMSServerURL + "status.json"
		}
	case XML:
		{
			sender.sendUrl = SMSServerURL + "send.xml"
			sender.sendBatchUrl = SMSServerURL + "send_batch.xml"
			sender.statusUrl = SMSServerURL + "status.xml"
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

func (p *SMSSender) Send(req SMSRequest) (resp Response, err error) {
	params := url.Values{}
	params.Add("mobile", req.Mobile)
	params.Add("message", req.Message)

	request, err := http.NewRequest("POST", p.sendUrl, bytes.NewBuffer([]byte(params.Encode())))
	if err != nil {
		log.Fatal(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", p.auth.BasicAuthorization())

	client := p.clientPool.Get().(*http.Client)
	defer p.clientPool.Put(client)

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return resp, errors.New(fmt.Sprintf("invalid http status code error: %s", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		log.Fatal(err)
		return
	}

	return
}

func (p *SMSSender) BatchSend(req BatchSMSRequest) (resp Response, err error) {
	params := url.Values{}
	params.Add("mobile_list", req.Mobiles)
	params.Add("message", req.Message)
	params.Add("time", req.Time)

	request, err := http.NewRequest("POST", p.sendUrl, bytes.NewBuffer([]byte(params.Encode())))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
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

func (p *SMSSender) Status() (status Status, err error) {
	params := url.Values{}
	request, err := http.NewRequest("POST", p.sendUrl, bytes.NewBuffer([]byte(params.Encode())))
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
