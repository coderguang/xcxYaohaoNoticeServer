package yaohaoNoticeSms

import (
	"encoding/json"
	yaohaoNoticeData "xcxYaohaoNoticeServer/src/data"
	yaohaoNoticeDef "xcxYaohaoNoticeServer/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func SendCommonSms(phone string, signName string, templateCode string, templateParam string) yaohaoNoticeDef.YaoHaoNoticeError {
	key, secret := yaohaoNoticeData.GetSmsInfo()
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", key, secret)
	if err != nil {
		sglog.Error("get sms client error,error:%s", err)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_SMS_CLIENT
	}

	sglog.Info("now try to send sms to phone:%s,singName:%s,templateCode:%s,templateParam:%s", phone, signName, templateCode, templateParam)

	if 0 == yaohaoNoticeData.GetSendSmsFlag() {
		sglog.Info("by config flag,this sms would't real send to phone %s", phone)
		return yaohaoNoticeDef.YAOHAO_NOTICE_OK
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = signName
	request.QueryParams["TemplateCode"] = templateCode
	request.QueryParams["TemplateParam"] = templateParam

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		sglog.Error("send sms to %s,signal:%s,error:%s", phone, signName, err)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_SMS_PROCESS
	}
	result := response.GetHttpContentString()
	sglog.Info(result)

	smsData := new(yaohaoNoticeDef.SMSData)
	p := &smsData

	if err := json.Unmarshal([]byte(result), p); err != nil {
		sglog.Error("send sms get result not a json,str=%s", result)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_SMS_RESULT_PARSE_ERROR
	}

	if smsData.Code == "OK" {
		return yaohaoNoticeDef.YAOHAO_NOTICE_OK
	}
	return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_SMS_OTHER
}

func SendConfirmMsg(phone string, randomCode string) yaohaoNoticeDef.YaoHaoNoticeError {

	signalName := "汽车摇号中签查询"
	templateCode := "SMS_171751048"
	templateParams := "{\"code\":\"" + randomCode + "\"}"

	return SendCommonSms(phone, signalName, templateCode, templateParams)
}

func SendLuckMsg(phone string, code string, timestr string) yaohaoNoticeDef.YaoHaoNoticeError {

	signalName := "汽车摇号中签查询"
	templateCode := "SMS_171751053"
	templateParams := "{\"name\":\"\",\"code\":\"" + code + "\",\"time\":\"" + timestr + "\"}"

	return SendCommonSms(phone, signalName, templateCode, templateParams)
}

func SendUnLuckMsg(phone string, code string, timestr string, endDt *sgtime.DateTime) yaohaoNoticeDef.YaoHaoNoticeError {

	signalName := "汽车摇号中签查询"
	templateCode := "SMS_171751055"
	templateParams := "{\"code\":\"" + code + "\"}"

	return SendCommonSms(phone, signalName, templateCode, templateParams)
}
