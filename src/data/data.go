package yaohaoNoticeData

import (
	"time"
	yaohaoNoticeConfig "xcxYaohaoNoticeServer/src/config"
	yaohaoNoticeDef "xcxYaohaoNoticeServer/src/define"

	"github.com/coderguang/GameEngine_go/sgthread"

	"github.com/coderguang/GameEngine_go/sgstring"
	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
)

var globalCfg *yaohaoNoticeDef.Config
var globalTokenMap *yaohaoNoticeDef.SecureSData
var globalRequireMap *yaohaoNoticeDef.SecureSRequireData
var globalOpenidMap *yaohaoNoticeDef.SecureWxOpenid
var globalNoticeRequireMap *yaohaoNoticeDef.SecureNoticeRequire

var globalRequireTimes int

func InitConfig(configfile string) {
	globalCfg = yaohaoNoticeConfig.ReadConfig(configfile)

	globalTokenMap = new(yaohaoNoticeDef.SecureSData)
	globalTokenMap.Data = make(map[string](map[string]*yaohaoNoticeDef.SData))

	globalRequireMap = new(yaohaoNoticeDef.SecureSRequireData)
	globalRequireMap.Data = make(map[string](map[string]*yaohaoNoticeDef.SRequireData))

	globalOpenidMap = new(yaohaoNoticeDef.SecureWxOpenid)
	globalOpenidMap.Data = make(map[string](map[string]*yaohaoNoticeDef.SWxOpenid))

	for _, v := range globalCfg.Title {
		globalOpenidMap.Data[v] = make(map[string]*yaohaoNoticeDef.SWxOpenid)
	}

	globalRequireTimes = 0

	globalNoticeRequireMap = new(yaohaoNoticeDef.SecureNoticeRequire)
	globalNoticeRequireMap.MapData = make(map[string](map[string]*yaohaoNoticeDef.SNoticeRequireData))

	for _, v := range globalCfg.Title {
		globalNoticeRequireMap.MapData[v] = make(map[string]*yaohaoNoticeDef.SNoticeRequireData)
	}

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

func GetRequireDataTableName() string {
	return globalCfg.DbTable + "_require_data"
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

func GetWxOpenid(title string, code string) (bool, string) {
	globalOpenidMap.Lock.RLock()
	defer globalOpenidMap.Lock.RUnlock()

	if v, ok := globalOpenidMap.Data[title]; ok {
		if vv, ok := v[code]; ok {
			now := sgtime.New()
			if now.GetTotalSecond()-vv.Time.GetTotalSecond() > 3600 {
				delete(v, code)
				return false, ""
			}
			return true, vv.Openid
		}
	}
	return false, ""
}

func AddWxOpenid(title string, data *yaohaoNoticeDef.SWxOpenid) {
	globalOpenidMap.Lock.RLock()
	defer globalOpenidMap.Lock.RUnlock()

	if v, ok := globalOpenidMap.Data[title]; ok {
		if vv, ok := v[data.Code]; ok {
			now := sgtime.New()
			if now.GetTotalSecond()-vv.Time.GetTotalSecond() > 3600 {
				delete(v, data.Code)
				v[data.Code] = data
			} else {
				sglog.Error("duplicate ,title:%s,code is %s,old openid:%s,new openid:%s", title, data.Code, vv.Openid, data.Openid)
			}
		} else {
			v[data.Code] = data
		}
	}
}

func ClearOpenidByTimer() {
	for {
		{
			sglog.Info("start to run clear openid data")
			globalOpenidMap.Lock.Lock()
			now := sgtime.New()
			for k, v := range globalOpenidMap.Data {
				sglog.Debug("delete openid data by timer,title:%s ,size:%d", k, len(v))
				for kk, vv := range v {
					if now.GetTotalSecond()-vv.Time.GetTotalSecond() > 3600 {
						sglog.Debug("delete openid data,title:%s ,code:%s,openid:%s", k, vv.Code, vv.Openid)
						delete(v, kk)
					}
				}
			}
			globalOpenidMap.Lock.Unlock()
			sglog.Info("clear openid data complete")
		}
		nowTime := time.Now()
		normalTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 59, 0, nowTime.Location())
		timeInt := normalTime.Sub(nowTime)
		sleepTime := int(timeInt/time.Second) + 1 // +1 for avoid loop run in that second time
		sglog.Info("next clear timer will run after %d seconds in %s", sleepTime, normalTime.String())
		sgthread.SleepBySecond(sleepTime)
	}
}

func GetNoticeDataByTitleAndCode(title string, code string) (bool, *yaohaoNoticeDef.SData) {
	flag, openid := GetWxOpenid(title, code)
	if !flag {
		sglog.Error("can't find openid,title:%s,code:%s", title, code)
		return false, nil
	}
	globalTokenMap.Lock.Lock()
	defer globalTokenMap.Lock.Unlock()

	if v, ok := globalTokenMap.Data[title]; ok {
		if vv, ook := v[openid]; ook {
			return true, vv
		}
	}
	return false, nil
}

func AddOrUpdateNoticeRequireData(data *yaohaoNoticeDef.SNoticeRequireData) {
	globalNoticeRequireMap.Lock.Lock()
	defer globalNoticeRequireMap.Lock.Unlock()

	if v, ok := globalNoticeRequireMap.MapData[data.Title]; ok {
		v[data.Openid] = data
	} else {
		tmp := make(map[string]*yaohaoNoticeDef.SNoticeRequireData)
		tmp[data.Openid] = data
		globalNoticeRequireMap.MapData[data.Title] = tmp
	}
}

func GetNoticeRequireData(title string, openid string) *yaohaoNoticeDef.SNoticeRequireData {
	globalNoticeRequireMap.Lock.Lock()
	defer globalNoticeRequireMap.Lock.Unlock()

	if v, ok := globalNoticeRequireMap.MapData[title]; ok {
		if vv, ook := v[openid]; ook {
			return vv
		}
	}
	return nil
}
