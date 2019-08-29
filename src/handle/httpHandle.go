package yaohaoNoticeHandle

import (
	"net/http"
	"strconv"
	"strings"
	yaohaoNoticeData "xcxYaohaoNoticeServer/src/data"
	yaohaoNoticeDef "xcxYaohaoNoticeServer/src/define"

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
	sglog.Debug("type=%s,keys:%s", reqType, keys)
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
		errcode, randomNum := RequireConfirmFromClient(title, token, cardType, code, phone, leftTime)
		if errcode != yaohaoNoticeDef.YAOHAO_NOTICE_OK {
			w.Write([]byte(getErrorCodeStr(errcode)))
			sglog.Debug("require sms error,final code=%d", errcode)
			return
		} else {
			str := "{\"errcode\":0," + "\"randomCode\":\"" + randomNum + "\"}"
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
