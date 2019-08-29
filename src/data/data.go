package yaohaoNoticeData

import (
	yaohaoNoticeConfig "xcxYaohaoNoticeServer/src/config"
	yaohaoNoticeDef "xcxYaohaoNoticeServer/src/define"

	"github.com/coderguang/GameEngine_go/sgstring"

	"github.com/coderguang/GameEngine_go/sglog"
)

var globalCfg *yaohaoNoticeDef.Config
var globalTokenMap *yaohaoNoticeDef.SecureSData
var globalRequireMap *yaohaoNoticeDef.SecureSRequireData

var globalRequireTimes int

func InitConfig(configfile string) {
	globalCfg = yaohaoNoticeConfig.ReadConfig(configfile)

	globalTokenMap = new(yaohaoNoticeDef.SecureSData)
	globalTokenMap.Data = make(map[string](map[string]*yaohaoNoticeDef.SData))

	globalRequireMap = new(yaohaoNoticeDef.SecureSRequireData)
	globalRequireMap.Data = make(map[string](map[string]*yaohaoNoticeDef.SRequireData))

	globalRequireTimes = 0
}

func GetTotalRequireTimes() int {
	return globalRequireTimes
}

func AddTotalRequireTimes() {
	globalRequireTimes++
}

func AddNoticeData(data *yaohaoNoticeDef.SData) {
	addNoticeDataByToken(data)
}

// func AddOrUpdateNoticeData(data *yaohaoNoticeDef.SData) *yaohaoNoticeDef.SData {

// 	if v, ok := globalTokenMap[data.Title]; ok {
// 		if data, okex := v[data.Token]; okex {
// 			//以token为主key

// 		}
// 	}

// }

func addNoticeDataByToken(data *yaohaoNoticeDef.SData) bool {
	keyData := data.Token

	globalTokenMap.Lock.Lock()
	defer globalTokenMap.Lock.Unlock()

	if v, ok := globalTokenMap.Data[data.Title]; ok {
		if detail, okex := v[keyData]; okex {
			sglog.Error("token:%s already regist", keyData)
			detail.ShowMsg()
			return false
		}
		v[keyData] = data
	} else {
		tmp := make(map[string]*yaohaoNoticeDef.SData)
		tmp[keyData] = data
		globalTokenMap.Data[data.Title] = tmp
	}
	return true
}

func GetTableName() string {
	return globalCfg.DbTable
}

func GetListenPort() string {
	return globalCfg.ListenPort
}

func GetSmsInfo() (string, string) {
	return globalCfg.SmsKey, globalCfg.SmsSecret
}

func GetDbConnectionData() (string, string, string, string, string) {
	return globalCfg.DbUser, globalCfg.DbPwd, globalCfg.DbUrl, globalCfg.DbPort, globalCfg.DbName
}

func IsDataAlreadyExist(title string, token string, code string, phone string) (bool, yaohaoNoticeDef.YaoHaoNoticeError, *yaohaoNoticeDef.SData) {
	if !IsValidTitle(title) {
		sglog.Debug("tile:%s not in config", title)
		return true, yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TITLE, nil
	}

	globalTokenMap.Lock.Lock()
	defer globalTokenMap.Lock.Unlock()

	if v, ok := globalTokenMap.Data[title]; ok {
		if data, okex := v[token]; okex {
			return true, yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TOKEN, data
		}
	}

	return false, yaohaoNoticeDef.YAOHAO_NOTICE_OK, nil
}

func IsValidTitle(title string) bool {
	return sgstring.EqualWithOr(title, globalCfg.Title)
}

func AddOrUpdateRequireData(data *yaohaoNoticeDef.SRequireData) {

	globalRequireMap.Lock.Lock()
	defer globalRequireMap.Lock.Unlock()

	if v, ok := globalRequireMap.Data[data.Title]; ok {
		v[data.Token] = data
	} else {
		tmp := make(map[string]*yaohaoNoticeDef.SRequireData)
		tmp[data.Token] = data
		globalRequireMap.Data[data.Title] = tmp
	}
	sglog.Info("add or update require data complete")
	data.ShowMsg()
}

func GetRequireData(title string, token string) *yaohaoNoticeDef.SRequireData {

	globalRequireMap.Lock.Lock()
	defer globalRequireMap.Lock.Unlock()

	if v, ok := globalRequireMap.Data[title]; ok {
		if vv, okex := v[token]; okex {
			return vv
		}
	}
	return nil
}

func GetSendSmsFlag() int {
	return globalCfg.SendSms
}

func ChangeSendSmsFlag(newFlag int) {
	globalCfg.SendSms = newFlag
}

func GetNeedNoticeData(title string) map[string]*yaohaoNoticeDef.SData {

	globalTokenMap.Lock.Lock()
	defer globalTokenMap.Lock.Unlock()

	if v, ok := globalTokenMap.Data[title]; ok {
		return v
	}
	return nil
}
