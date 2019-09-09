package main

import (
	"log"
	"os"
	"strconv"
	yaohaoNoticeData "xcxYaohaoNoticeServer/src/data"
	yaohaoNoticeDb "xcxYaohaoNoticeServer/src/db"
	yaohaoNoticeHandle "xcxYaohaoNoticeServer/src/handle"

	"github.com/coderguang/GameEngine_go/sgcmd"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgserver"
)

func ChangeSmsFlag(cmdstr []string) {
	flag, _ := strconv.Atoi(cmdstr[1])
	yaohaoNoticeData.ChangeSendSmsFlag(flag)
}

func ShowSmsFlag(cmdstr []string) {
	sglog.Info("smsFlag:%d", yaohaoNoticeData.GetSendSmsFlag())
}

func RegistNoticeData(cmdstr []string) {
	leftTime, _ := strconv.Atoi(cmdstr[5])
	yaohaoNoticeHandle.AddNoticeRequire(cmdstr[1], cmdstr[2], cmdstr[3], cmdstr[4], leftTime)
}

func RequireConfirmSms(cmdstr []string) {
	leftTime, _ := strconv.Atoi(cmdstr[6])
	cardType, _ := strconv.Atoi(cmdstr[3])
	yaohaoNoticeHandle.RequireConfirmFromClient(cmdstr[1], cmdstr[2], cardType, cmdstr[4], cmdstr[5], leftTime)
}

func ConfirmRequire(cmdstr []string) {
	yaohaoNoticeHandle.ConfireRequireFromClient(cmdstr[1], cmdstr[2], cmdstr[3])
}

func RegistCmd() {
	// ["RegistNoticeData","guangzhou","sgewgtew","4654694","18826409048","2"]
	sgcmd.RegistCmd("RegistNoticeData", "[\"RegistNoticeData\",\"title\",\"token\",\"code\",\"phone\",\"2\"] add notice data", RegistNoticeData)
	// ["RequireConfirmSms","guangzhou","sgewgtew","4654694","18826409048","2"]
	sgcmd.RegistCmd("RequireConfirmSms", "[\"RequireConfirmSms\",\"title\",\"token\",\"code\",\"phone\",\"2\"] add require data", RequireConfirmSms)
	// ["ConfirmRequire","guangzhou","sgewgtew","4654694"]
	sgcmd.RegistCmd("ConfirmRequire", "[\"ConfirmRequire\",\"title\",\"token\",,\"cardType\",\"randomcode\"] confirm require", ConfirmRequire)
	//["ChangeSmsFlag","1"]
	sgcmd.RegistCmd("ChangeSmsFlag", "[\"ChangeSmsFlag\",\"1\"] change sms real send flag,1 to send,0 to shutdown sms", ChangeSmsFlag)
	//["ShowSmsFlag"]
	sgcmd.RegistCmd("ShowSmsFlag", "[\"ShowSmsFlag\"] show current sms flag", ShowSmsFlag)
}

//==============main=============
func main() {

	sgserver.StartLogServer("debug", "./log/", log.LstdFlags, true)

	arg_num := len(os.Args) - 1
	if arg_num < 1 {
		sglog.Fatal("please input config file")
		return
	}
	configfile := os.Args[1]
	sglog.Info("read global config from ", configfile)
	yaohaoNoticeData.InitConfig(configfile)

	sglog.Info("start run yaohao notice program")
	sglog.Info("start connect to db")

	yaohaoNoticeDb.InitDbConnection(yaohaoNoticeData.GetDbConnectionData())

	yaohaoNoticeDb.LoadNoticeDataFromDb()

	go yaohaoNoticeHandle.HttpNoticeServer(yaohaoNoticeData.GetListenPort())

	go yaohaoNoticeData.ClearOpenidByTimer()

	RegistCmd()

	sgcmd.StartCmdWaitInputLoop()

}
