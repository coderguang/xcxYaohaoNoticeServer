package yaohaoNoticeHandle

import (
	"strconv"
	yaohaoNoticeData "xcxYaohaoNoticeServer/src/data"
	yaohaoNoticeDb "xcxYaohaoNoticeServer/src/db"
	yaohaoNoticeDef "xcxYaohaoNoticeServer/src/define"
	yaohaoNoticeSms "xcxYaohaoNoticeServer/src/sms"

	"github.com/coderguang/GameEngine_go/sgstring"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sgregex"

	"github.com/coderguang/GameEngine_go/sglog"
)

func AddNoticeRequire(title string, token string, code string, phone string, lefttime int) yaohaoNoticeDef.YaoHaoNoticeError {

	if lefttime <= 0 || lefttime > 3 {
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_LEFT_TIME
	}

	if ok, errcode, _ := yaohaoNoticeData.IsDataAlreadyExist(title, token, code, phone); ok {
		sglog.Debug("data already exist,title:%s,token:%s,code:%s,phone:%s", title, token, code, phone)
		return errcode
	}

	if !sgregex.AllNum(code) {
		sglog.Error("error code,title:%s,code:%s", title, code)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CODE
	}
	if !sgregex.CNMobile(phone) {
		sglog.Error("error phone,title:%s,phone:%s", title, phone)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_PHONE
	}

	data := new(yaohaoNoticeDef.SData)
	data.Title = title
	data.Token = token
	data.Code = code
	data.Phone = phone
	data.EndDt = sgtime.New()
	data.EndDt.AddDate(0, lefttime, 0)
	data.Desc = ""
	data.RenewTimes = 0
	data.Status = yaohaoNoticeDef.YAOHAO_NOTICE_STATUS_NORMAL
	data.NoticeTimes = 0

	yaohaoNoticeData.AddNoticeData(data)

	yaohaoNoticeDb.InsertOrUpdateData(data)

	sglog.Info("add notice data ok,token=%s", data.Token)
	data.ShowMsg()

	return yaohaoNoticeDef.YAOHAO_NOTICE_OK
}

func CheckCardTypeValid(cardType int) bool {
	return 1 == cardType || 2 == cardType
}

func CheckCodeValid(title string, code string) bool {
	switch title {
	case "guangzhou":
		if !sgregex.AllNum(code) {
			return false
		}
		if len(code) != 13 {
			return false
		}
		return true
	}
	return false
}

//请求验证下发
func RequireConfirmFromClient(title string, token string, cardType int, code string, phone string, lefttime int) (yaohaoNoticeDef.YaoHaoNoticeError, string) {
	randomCode := ""
	if lefttime <= 0 || lefttime > 3 {
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_LEFT_TIME, randomCode
	}

	if !CheckCardTypeValid(cardType) {
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_HTTP_REQ_CARD_TYPE, randomCode
	}

	if !CheckCodeValid(title, code) {
		sglog.Error("require error code,title:%s,code:%s", title, code)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CODE, randomCode
	}

	if !sgregex.CNMobile(phone) {
		sglog.Error("require error phone,title:%s,phone:%s", title, phone)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_PHONE, randomCode
	}

	if !yaohaoNoticeData.CanBindPhone(phone) {
		sglog.Error("require error phone,bind too many token,title:%s,phone:%s", title, phone)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_PHONE_BIND_TOO_MANY, randomCode
	}

	now := sgtime.New()

	if ok, errcode, existData := yaohaoNoticeData.IsDataAlreadyExist(title, token, code, phone); ok {
		sglog.Debug("require data already exist,title:%s,token:%s,code:%s,phone:%s", title, token, code, phone)

		if errcode == yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TITLE {
			return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TITLE, randomCode
		}

		if existData.Status == yaohaoNoticeDef.YAOHAO_NOTICE_STATUS_GM_LIMIT {
			return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_GM_LIMIT, randomCode
		}

		if existData.IsStillValid() {

			if !existData.IsDataChange(code, phone, cardType) {
				//没有信息变更
				if errcode == yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TOKEN {
					return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TOKEN_STILL_VALID, randomCode
				} else if errcode == yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CODE {
					return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CODE_STILL_VALID, randomCode
				} else if errcode == yaohaoNoticeDef.YAOHAO_NOTICE_ERR_PHONE {
					return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_PHONE_STILL_VALID, randomCode
				}
			}
		}
	}

	oldData := yaohaoNoticeData.GetRequireData(title, token)

	if oldData != nil {

		now := sgtime.New()

		if !oldData.IsDataChange(code, phone, cardType) {
			//没有信息变更
			if now.GetTotalSecond()-oldData.RequireDt.GetTotalSecond() <= int64(yaohaoNoticeDef.YAOHAO_NOTICE_REQUIRE_VALID_TIME) {
				if oldData.Status == int(yaohaoNoticeDef.YaoHaoNoticeRequireStatus_Answer_Complete) {
					return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_REQUIRE_HAD_CONFIRM, randomCode
				}
				if oldData.AnswerTimes >= yaohaoNoticeDef.YAOHAO_NOTICE_CONFIRM_TIMES {
					return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CONFIRM_MORE_TIMES, randomCode
				} else {
					return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_REQUIRE_WAIT_ANSWER, randomCode
				}
			} else if now.GetTotalSecond()-oldData.RequireDt.GetTotalSecond() <= int64(yaohaoNoticeDef.YAOHAO_NOTICE_REQUIRE_UNLOCK_TIME) {
				if oldData.RequireTimes >= yaohaoNoticeDef.YAOHAO_NOTICE_REQUIRE_MAX_TIMES {
					return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_REQUIRE_HAD_LOCK, randomCode
				}
			} else {
				oldData.RequireTimes = 0 //reset
				oldData.AnswerTimes = 0
			}
		} else {
			//有数据变更
			if now.GetTotalSecond()-oldData.RequireDt.GetTotalSecond() <= int64(yaohaoNoticeDef.YAOHAO_NOTICE_SMS_TIME_LIMIT) {
				return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_REQUIRE_HAD_LOCK, randomCode
			}
			if now.GetTotalSecond()-oldData.RequireDt.GetTotalSecond() <= int64(yaohaoNoticeDef.YAOHAO_NOTICE_REQUIRE_UNLOCK_TIME) {
				if oldData.RequireTimes >= yaohaoNoticeDef.YAOHAO_NOTICE_REQUIRE_MAX_TIMES {
					return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_REQUIRE_HAD_LOCK, randomCode
				}
			} else {
				oldData.RequireTimes = 0
				oldData.AnswerTimes = 0
			}
		}

		oldData.RequireDt = now
		oldData.Status = int(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_REQUIRE_WAIT_ANSWER)
		oldData.Token = token
		oldData.Code = code
		oldData.CardType = cardType
		oldData.Phone = phone
		oldData.LeftTime = lefttime
		oldData.RandomNum = sgstring.RandNumStringRunes(yaohaoNoticeDef.YAOHAO_NOTICE_RANDOM_NUM_LENGTH)
		oldData.RequireTimes++

		yaohaoNoticeData.AddOrUpdateRequireData(oldData)

		smsCode := yaohaoNoticeSms.SendConfirmMsg(oldData.Phone, oldData.RandomNum)

		if yaohaoNoticeDef.YAOHAO_NOTICE_OK != smsCode {
			return smsCode, randomCode
		}
	} else {

		//针对绑定后取消的限制

		if !yaohaoNoticeData.CanGetRequire(token) {
			sglog.Info("title:%s,token:%s,require too fast,limit it", title, token)
			return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_SMS_CLIENT, randomCode
		}

		newRequireData := new(yaohaoNoticeDef.SRequireData)
		newRequireData.Title = title
		newRequireData.Status = int(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_REQUIRE_WAIT_ANSWER)
		newRequireData.AnswerTimes = 0
		newRequireData.RequireDt = now
		newRequireData.Token = token
		newRequireData.CardType = cardType
		newRequireData.Code = code
		newRequireData.Phone = phone
		newRequireData.LeftTime = lefttime
		newRequireData.RandomNum = sgstring.RandNumStringRunes(yaohaoNoticeDef.YAOHAO_NOTICE_RANDOM_NUM_LENGTH)
		newRequireData.RequireTimes = 0

		randomCode = newRequireData.RandomNum
		yaohaoNoticeData.AddOrUpdateRequireData(newRequireData)

		smsCode := yaohaoNoticeSms.SendConfirmMsg(newRequireData.Phone, newRequireData.RandomNum)

		if yaohaoNoticeDef.YAOHAO_NOTICE_OK != smsCode {
			return smsCode, randomCode
		}
	}
	yaohaoNoticeData.AddRequireTimeLimits(token)
	return yaohaoNoticeDef.YAOHAO_NOTICE_OK, randomCode
}

func ConfireRequireFromClient(title string, token string, randomCode string) yaohaoNoticeDef.YaoHaoNoticeError {
	oldData := yaohaoNoticeData.GetRequireData(title, token)

	if oldData == nil {
		sglog.Debug("no require ,title:%s,token:%s,randomCode:%s", title, token, randomCode)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CONFIRM_NOT_REQUIRE
	}

	now := sgtime.New()

	if now.GetTotalSecond()-oldData.RequireDt.GetTotalSecond() <= int64(yaohaoNoticeDef.YAOHAO_NOTICE_REQUIRE_VALID_TIME) {
		if oldData.Status == int(yaohaoNoticeDef.YaoHaoNoticeRequireStatus_Answer_Complete) {
			sglog.Debug("had answer already ,title:%s,token:%s,randomCode:%s", title, token, randomCode)
			return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_REQUIRE_HAD_CONFIRM
		}
	} else {
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_HTTP_RANDOM_CODE_TIME_OUT
	}
	if oldData.AnswerTimes >= yaohaoNoticeDef.YAOHAO_NOTICE_CONFIRM_TIMES {
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CONFIRM_MORE_TIMES
	}
	oldData.AnswerTimes++
	if oldData.RandomNum != randomCode {
		sglog.Debug("error randomcode ,title:%s,token:%s,randomCode:%s", title, token, randomCode)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CONFIRM_RANDOMCODE
	}
	//验证通过
	//保证 title,token 组合 必须唯一
	if ok, errcode, existData := yaohaoNoticeData.IsDataAlreadyExist(title, oldData.Token, oldData.Code, oldData.Phone); ok {
		sglog.Debug("require data already be bind by other,title:%s,token:%s,code:%s,phone:%s", title, oldData.Token, oldData.Code, oldData.Phone)

		if existData.Status == yaohaoNoticeDef.YAOHAO_NOTICE_STATUS_GM_LIMIT {
			return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_GM_LIMIT
		}

		if existData.IsStillValid() && !existData.IsDataChange(oldData.Code, oldData.Phone, oldData.CardType) {
			if errcode == yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TOKEN {
				return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TOKEN_STILL_VALID
			} else if errcode == yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CODE {
				return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_CODE_STILL_VALID
			} else if errcode == yaohaoNoticeDef.YAOHAO_NOTICE_ERR_PHONE {
				return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_PHONE_STILL_VALID
			}
		} else {

			//旧数据已过期,更改数据

			yaohaoNoticeData.DelPhoneBind(existData.Phone)
			yaohaoNoticeData.AddPhoneBind(oldData.Phone)

			existData.Code = oldData.Code
			existData.Phone = oldData.Phone
			firstOfMonth := sgtime.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
			firstOfMonth.AddDate(0, oldData.LeftTime, 0)
			existData.EndDt = firstOfMonth
			existData.CardType = oldData.CardType
			existData.Desc = ""
			existData.RenewTimes++
			existData.Status = yaohaoNoticeDef.YAOHAO_NOTICE_STATUS_NORMAL

			yaohaoNoticeDb.InsertOrUpdateData(existData)

		}
	} else { //原本不存在数据

		oldData.Status = int(yaohaoNoticeDef.YaoHaoNoticeRequireStatus_Answer_Complete)

		noticeData := new(yaohaoNoticeDef.SData)
		noticeData.Token = oldData.Token
		noticeData.Name = ""
		noticeData.Title = oldData.Title
		noticeData.Code = oldData.Code
		noticeData.Phone = oldData.Phone
		noticeData.CardType = oldData.CardType

		firstOfMonth := sgtime.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
		firstOfMonth.AddDate(0, oldData.LeftTime, 0)
		noticeData.EndDt = firstOfMonth
		noticeData.Desc = ""
		noticeData.RenewTimes = 0
		noticeData.Status = yaohaoNoticeDef.YAOHAO_NOTICE_STATUS_NORMAL
		noticeData.NoticeTimes = 0

		yaohaoNoticeData.AddNoticeData(noticeData)
		yaohaoNoticeDb.InsertOrUpdateData(noticeData)
	}
	yaohaoNoticeData.RemoveRequireData(title, token)
	return yaohaoNoticeDef.YAOHAO_NOTICE_OK
}

func RecvDataFromYaoHaoServer(title string, cardType int, timestr string, totalSize string, datas []string) yaohaoNoticeDef.YaoHaoNoticeError {
	if !yaohaoNoticeData.IsValidTitle(title) {
		sglog.Error("RecvDataFromYaoHaoServer unvalid title:%s", title)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TITLE
	}
	if !CheckCardTypeValid(cardType) {
		sglog.Error("RecvDataFromYaoHaoServer unvalid cardType:%d", cardType)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_HTTP_REQ_CARD_TYPE
	}

	dataSize, err := strconv.Atoi(totalSize)
	if err != nil {
		sglog.Error("RecvDataFromYaoHaoServer size fromat error:%s", totalSize)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_YAOHAO_SERVER_INT_FORMAT
	}
	if dataSize != len(datas) {
		sglog.Error("RecvDataFromYaoHaoServer data size not match: size:%d,real size:%d", dataSize, len(datas))
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_YAOHAO_SERVER_DATA_SIZE_NOT_MATCH
	}

	now := sgtime.New()
	currentYearMonth := now.YearString() + now.MonthString()
	if currentYearMonth != timestr {
		sglog.Error("RecvDataFromYaoHaoServer time not match data time:%s,current time:%s", timestr, currentYearMonth)
		return yaohaoNoticeDef.YAOHAO_NOTICE_ERR_YAOHAO_SERVER_TIME_NOT_MATCH
	}

	luckBoys := make(map[string]string)

	for _, v := range datas {
		luckBoys[v] = ""
	}

	nextMonth := sgtime.New()
	nextMonth.AddDate(0, 1, 0)
	luckBoyNum := 0
	unLuckBoyNum := 0
	needNoticeData := yaohaoNoticeData.GetNeedNoticeData(title)
	if needNoticeData != nil {
		for _, v := range needNoticeData {
			if v.Status != yaohaoNoticeDef.YAOHAO_NOTICE_STATUS_NORMAL {
				sglog.Debug("title:%s,phone:%s,token:%s status not normal", title, v.Phone, v.Token)
				continue
			}
			if v.EndDt.Before(now) {
				sglog.Debug("title:%s,phone:%s,token:%s status out of time,endDt is %s,now is %s", title, v.Phone, v.Token, v.EndDt.NormalString(), now.NormalString())
				continue
			}
			if v.CardType != cardType {
				continue
			}

			v.NoticeTimes++
			if v.EndDt.Before(nextMonth) {
				v.Status = yaohaoNoticeDef.YAOHAO_NOTICE_STATUS_TIME_OUT
				yaohaoNoticeData.DelPhoneBind(v.Phone)
			}

			if _, luckFlag := luckBoys[v.Code]; luckFlag {
				errcode := yaohaoNoticeSms.SendLuckMsg(v.Phone, v.Code, timestr)
				if yaohaoNoticeDef.YAOHAO_NOTICE_OK == errcode {
					sglog.Debug("title:%s,phone:%s,token:%s send luck msg ok", title, v.Phone, v.Token)
				} else {
					sglog.Debug("title:%s,phone:%s,token:%s send luck msg error,code=%d", title, v.Phone, errcode)
				}
				v.Status = yaohaoNoticeDef.YAOHAO_NOTICE_STATUS_CANCEL_BY_GM_BECASURE_LUCK
				luckBoyNum++
			} else {
				errcode := yaohaoNoticeSms.SendUnLuckMsg(v.Phone, v.Code, timestr, v.EndDt)
				if yaohaoNoticeDef.YAOHAO_NOTICE_OK == errcode {
					sglog.Debug("title:%s,phone:%s,token:%s send unluck msg ok", title, v.Phone, v.Token)
				} else {
					sglog.Debug("title:%s,phone:%s,token:%s send unluck msg error,code=%d", title, v.Phone, v.Token, errcode)
				}
				unLuckBoyNum++
			}

			yaohaoNoticeDb.InsertOrUpdateData(v)
		}
		sglog.Info("notice tile %s totalsize:%d,luck:%d,unlock:%d", title, len(needNoticeData), luckBoyNum, unLuckBoyNum)
	} else {
		sglog.Debug("no any data need notice ,title:%s", title)
	}
	return yaohaoNoticeDef.YAOHAO_NOTICE_OK
}
