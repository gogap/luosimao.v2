package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xigang/luosimao.v2"
)

var voiceAuth luosimao.Authorization
var smsAuth luosimao.Authorization

func main() {

	voiceAuth = luosimao.Authorization{UserName: "api", Password: ""}
	smsAuth = luosimao.Authorization{UserName: "api", Password: ""}

	send_voice("13400000000", 123456)
	send_sms("13400000000", "您的验证码是:1234【日日进】")
	send_sms_batch("13400000000,13400000000", "您的验证码是:12345 【公司签名】")
	voice_status()
	sms_status()
}

func send_voice(mobile string, code int32) {
	sender := luosimao.NewVoiceSender(voiceAuth, luosimao.JSON, time.Second*5)
	resp, err := sender.Send(luosimao.VoiceRequest{Mobile: mobile, Code: code})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		b, _ := json.Marshal(resp)
		fmt.Println(string(b))
		fmt.Println(resp.ErrorDescription())
	}
}

func voice_status() {
	sender := luosimao.NewVoiceSender(voiceAuth, luosimao.JSON, time.Second*5)
	resp, err := sender.Status()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		b, _ := json.Marshal(resp)
		fmt.Println(string(b))
		fmt.Println(resp.ErrorDescription())
	}
}

func send_sms(mobile, message string) {
	sender := luosimao.NewSMSSender(smsAuth, luosimao.JSON, time.Second*5)
	fmt.Println(">>>", smsAuth)
	resp, err := sender.Send(luosimao.SMSRequest{Mobile: mobile, Message: message})
	fmt.Println("cellphone and message:", mobile, message)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		b, _ := json.Marshal(resp)
		fmt.Println(string(b))
		fmt.Println(resp.ErrorDescription())
	}
}

func sms_status() {
	sender := luosimao.NewSMSSender(smsAuth, luosimao.JSON, time.Second*5)
	resp, err := sender.Status()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		b, _ := json.Marshal(resp)
		fmt.Println(string(b))
		fmt.Println(resp.ErrorDescription())
	}
}

func send_sms_batch(mobiles string, message string) {
	sender := luosimao.NewSMSSender(smsAuth, luosimao.JSON, time.Second*5)
	resp, err := sender.BatchSend(luosimao.BatchSMSRequest{Mobiles: mobiles, Message: message, Time: ""})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		b, _ := json.Marshal(resp)
		fmt.Println(string(b))
		fmt.Println(resp.ErrorDescription())
	}
}
