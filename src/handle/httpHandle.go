package yaohaoNoticeHandle

import (
	"net/http"
	"strconv"
	"strings"
	yaohaoNoticeData "xcxYaohaoNoticeServer/src/data"
	yaohaoNoticeDef "xcxYaohaoNoticeServer/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
)

type wx_xcx_yaohao_notice_handler struct{}

func getErrorCodeStr(code yaohaoNoticeDef.YaoHaoNoticeError) string {
	codeStr := strconv.Itoa(int(code))
	str := "{\"errcode\":" + codeStr + "}"
	sglog.Debug("return is %s", str)
	return str
}

func doLogic(w http.ResponseWriter, r *http.Request) {

	sglog.Info("get require from client,times=%d", yaohaoNoticeData.GetTotalRequireTimes())
	r.ParseForm()

	if len(r.Form["key"]) <= 0 {
		w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_HTTP_NO_KEY))) // not param keys
		sglog.Debug("no key in this handle")
		return
	}

	rawkeys := r.Form["key"][0]
	keys := strings.Split(rawkeys, ",")
	reqType := keys[0]
	sglog.Debug("type=%s,keys:%s,total:%s", reqType, keys, r.Form)
	if reqType == "require" {
		// ?key=require,title,token,type,code,phone,lefttime
		paramSizeMin := 7
		if len(keys) < paramSizeMin {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_HTTP_PARAM_NUM))) // not param keys
			sglog.Debug("require not enough params")
			return
		}
		leftTime, err := strconv.Atoi(keys[paramSizeMin-1])
		if err != nil {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_LEFT_TIME))) // not param keys
			sglog.Debug("require not enough params")
			return
		}
		index := 1
		title := keys[index]
		index++
		token := keys[index]
		index++
		cardType, _ := strconv.Atoi(keys[index])
		index++
		code := keys[index]
		index++
		phone := keys[index]

		flag, openid := yaohaoNoticeData.GetWxOpenid(title, token)
		if !flag {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TOKEN)))
			sglog.Error("unknow openid require ,title:%s,token:%s", title, token)
			return
		}
		token = openid

		errcode, randomNum := RequireConfirmFromClient(title, token, cardType, code, phone, leftTime)
		if errcode != yaohaoNoticeDef.YAOHAO_NOTICE_OK {
			w.Write([]byte(getErrorCodeStr(errcode)))
			sglog.Debug("require sms error,final code=%d", errcode)
			return
		} else {
			str := getErrorCodeStr(0)
			if 0 == yaohaoNoticeData.GetSendSmsFlag() {
				sglog.Info("send sms by config flag,this sms would't real send to phone %s,will return by http", phone)
				str = "{\"errcode\":0," + "\"randomCode\":\"" + randomNum + "\"}"
			}
			w.Write([]byte(str)) // not param keys
			sglog.Debug("require sms ok")
			return
		}
	} else if reqType == "confirm" {
		// ?key=confirm,title,token,randomcode
		if len(keys) < 4 {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_HTTP_PARAM_NUM))) // not param keys
			sglog.Debug("confirm not enough params")
			return
		}
		title := keys[1]
		token := keys[2]
		randomcode := keys[3]

		flag, openid := yaohaoNoticeData.GetWxOpenid(title, token)
		if !flag {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_TOKEN)))
			sglog.Error("unknow openid require ,title:%s,token:%s", title, token)
			return
		}
		token = openid

		errcode := ConfireRequireFromClient(title, token, randomcode)
		if errcode != yaohaoNoticeDef.YAOHAO_NOTICE_OK {
			w.Write([]byte(getErrorCodeStr(errcode))) // not param keys
			sglog.Debug("confirm sms error")
			return
		}
	} else if reqType == "data" {
		// data,title,time,cardType,len,detail
		//?key = data, title, time,type, totalnum, data
		if len(keys) < 5 {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_HTTP_PARAM_NUM))) // not param keys
			sglog.Debug("data not enough at least params")
			return
		}
		title := keys[1]
		time := keys[2]
		cardType, _ := strconv.Atoi(keys[3])
		totalSize := keys[4]
		datas := keys[5:len(keys)]
		RecvDataFromYaoHaoServer(title, cardType, time, totalSize, datas)
	} else if reqType == "openid" {
		//?key =openid,title,code,openid
		if len(keys) < 4 {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_OPEN_ID_PARAM_NUM))) // not param keys
			sglog.Debug("openid not enough at least params")
			return
		}

		openidData := new(yaohaoNoticeDef.SWxOpenid)
		title := keys[1]
		openidData.Code = keys[2]
		openidData.Openid = keys[3]
		openidData.Time = sgtime.New()

		sglog.Info("receive openid,title:%s,code:%s,openid:%s", title, openidData.Code, openidData.Openid)

		yaohaoNoticeData.AddWxOpenid(title, openidData)
	} else if reqType == "getData" {
		//?key =getData,title,code
		if len(keys) < 3 {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_OPEN_ID_PARAM_NUM))) // not param keys
			sglog.Debug("getData not enough at least params")
			return
		}
		title := keys[1]
		code := keys[2]

		flag, sdata := yaohaoNoticeData.GetNoticeDataByTitleAndCode(title, code)
		if flag {
			w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_NOT_BIND_DATA)))
			sglog.Debug("not bind data before,title:%s,code:%s", title, code)
			return
		} else {
			str := "{\"errcode\":0," + "\"code\":\"" + sdata.Code + ",\"phone\":\"" + sdata.Phone + "\"}"
			w.Write([]byte(str))
			sglog.Debug("find bind data,title:%s,Code:%s,phone:%s", title, sdata.Code, sdata.Phone)
			return
		}

	} else {
		w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_ERR_HTTP_REQ_TYPE))) // not param keys
		sglog.Debug("type error")
		return
	}

	w.Write([]byte(getErrorCodeStr(yaohaoNoticeDef.YAOHAO_NOTICE_OK))) // ok
	return
}

func (h *wx_xcx_yaohao_notice_handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	doLogic(w, r)
}

func HttpNoticeServer(port string) {

	http.Handle("/", &wx_xcx_yaohao_notice_handler{})
	listenport := "0.0.0.0:" + port
	sglog.Info("start notice http server,port is %s", port)
	http.ListenAndServe(listenport, nil)
}
