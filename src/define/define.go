package yaohaoNoticeDef

import (
	"sync"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgtime"
)

const YAOHAO_NOTICE_REQUIRE_VALID_TIME int = 300
const YAOHAO_NOTICE_REQUIRE_UNLOCK_TIME int = 1800
const YAOHAO_NOTICE_CONFIRM_TIMES int = 3
const YAOHAO_NOTICE_RANDOM_NUM_LENGTH int = 4
const YAOHAO_NOTICE_REQUIRE_MAX_TIMES int = 3

const (
	YAOHAO_NOTICE_STATUS_NORMAL                     = "0"
	YAOHAO_NOTICE_STATUS_CANCEL                     = "1"
	YAOHAO_NOTICE_STATUS_TIME_OUT                   = "2"
	YAOHAO_NOTICE_STATUS_GM_LIMIT                   = "3"
	YAOHAO_NOTICE_STATUS_CANCEL_BY_GM_BECASURE_LUCK = "4"
)

type YaoHaoNoticeError int

const (
	YAOHAO_NOTICE_OK                                    YaoHaoNoticeError = iota //0
	YAOHAO_NOTICE_ERR_DATA_EXISTS                                                //1
	YAOHAO_NOTICE_ERR_TITLE                                                      //2
	YAOHAO_NOTICE_ERR_TOKEN                                                      //3
	YAOHAO_NOTICE_ERR_PHONE                                                      //4
	YAOHAO_NOTICE_ERR_CODE                                                       //5
	YAOHAO_NOTICE_ERR_LEFT_TIME                                                  //6
	YAOHAO_NOTICE_ERR_GM_LIMIT                                                   //7
	YAOHAO_NOTICE_ERR_CONFIRM_MORE_TIMES                                         //8
	YAOHAO_NOTICE_ERR_TOKEN_STILL_VALID                                          //9
	YAOHAO_NOTICE_ERR_CODE_STILL_VALID                                           //10
	YAOHAO_NOTICE_ERR_PHONE_STILL_VALID                                          //11
	YAOHAO_NOTICE_ERR_REQUIRE_HAD_CONFIRM                                        //12 已应答
	YAOHAO_NOTICE_ERR_REQUIRE_WAIT_ANSWER                                        //13 等待应答
	YAOHAO_NOTICE_ERR_REQUIRE_HAD_LOCK                                           //14 锁定
	YAOHAO_NOTICE_ERR_CONFIRM_NOT_REQUIRE                                        //15 未请求
	YAOHAO_NOTICE_ERR_CONFIRM_RANDOMCODE                                         //16 错误的验证码
	YAOHAO_NOTICE_ERR_HTTP_NO_KEY                                                //17
	YAOHAO_NOTICE_ERR_HTTP_PARAM_NUM                                             //18
	YAOHAO_NOTICE_ERR_HTTP_REQ_TYPE                                              //19
	YAOHAO_NOTICE_ERR_HTTP_REQ_MAX_TIMES                                         //20
	YAOHAO_NOTICE_ERR_HTTP_RANDOM_CODE_TIME_OUT                                  //21
	YAOHAO_NOTICE_ERR_SMS_CLIENT                                                 //22
	YAOHAO_NOTICE_ERR_SMS_PROCESS                                                //23
	YAOHAO_NOTICE_ERR_SMS_OTHER                                                  //24
	YAOHAO_NOTICE_ERR_SMS_RESULT_PARSE_ERROR                                     //25
	YAOHAO_NOTICE_ERR_YAOHAO_SERVER_INT_FORMAT                                   //26
	YAOHAO_NOTICE_ERR_YAOHAO_SERVER_DATA_SIZE_NOT_MATCH                          //27
	YAOHAO_NOTICE_ERR_YAOHAO_SERVER_TIME_NOT_MATCH                               //28
	YAOHAO_NOTICE_ERR_HTTP_REQ_CARD_TYPE                                         //29
	YAOHAO_NOTICE_ERR_OPEN_ID_PARAM_NUM                                          //30
)

type YaoHaoNoticeRequireStatus int

const (
	YaoHaoNoticeRequireStatus_Wait_Answer YaoHaoNoticeRequireStatus = iota
	YaoHaoNoticeRequireStatus_Answer_Complete
	YaoHaoNoticeRequireStatus_Wait_ReAnswer //应答错误再次等待
)

type Config struct {
	Title      []string `json:"title"`
	DbUrl      string   `json:"dbUrl"`
	DbPort     string   `json:"dbPort"`
	DbUser     string   `json:"dbUser"`
	DbPwd      string   `json:"dbPwd"`
	DbName     string   `json:"dbName"`
	DbTable    string   `json:"dbTable"`
	ListenPort string   `json:"listenPort"`
	SmsKey     string   `json:"smsKey"`
	SmsSecret  string   `json:"smsSecret"`
	SendSms    int      `json:"sendSms"`
}

type SData struct {
	Token       string
	Name        string
	Title       string
	CardType    int
	Code        string
	Phone       string
	EndDt       *sgtime.DateTime
	Desc        string
	RenewTimes  int
	Status      string
	NoticeTimes int
}

type SecureSData struct {
	Data map[string](map[string]*SData)
	Lock sync.RWMutex
}

func (data *SData) ShowMsg() {
	sglog.Debug("token:%s", data.Token)
	sglog.Debug("Name:%s", data.Name)
	sglog.Debug("Title:%s", data.Title)
	sglog.Debug("CardType:%d", data.CardType)
	sglog.Debug("Code:%s", data.Code)
	sglog.Debug("Phone:%s", data.Phone)
	sglog.Debug("EndDt:%s", data.EndDt.NormalString())
	sglog.Debug("Desc:%s", data.Desc)
	sglog.Debug("RenewTimes:%s", data.RenewTimes)
}

func (data *SData) IsDataChange(code string, phone string, cardType int) bool {
	if data.Code == code && data.Phone == phone && data.CardType == cardType {
		return false
	}
	return true
}

func (data *SData) IsStillValid() bool {
	now := sgtime.New()
	if data.Status == YAOHAO_NOTICE_STATUS_NORMAL && now.Before(data.EndDt) {
		return true
	}
	return false
}

type SRequireData struct {
	Token        string
	Title        string
	CardType     int
	Code         string
	Phone        string
	RequireDt    *sgtime.DateTime //请求时间
	RandomNum    string
	AnswerTimes  int //回应次数
	Status       int
	LeftTime     int
	RequireTimes int
}

type SecureSRequireData struct {
	Data map[string](map[string]*SRequireData)
	Lock sync.RWMutex
}

func (data *SRequireData) IsDataChange(code string, phone string, cardType int) bool {
	if data.Code == code && data.Phone == phone && data.CardType == cardType {
		return false
	}
	return true
}

func (data *SRequireData) ShowMsg() {
	sglog.Debug("=======start===========")
	sglog.Debug("token:%s", data.Token)
	sglog.Debug("Title:%s", data.Title)
	sglog.Debug("Code:%s", data.Code)
	sglog.Debug("CardType:%d", data.CardType)
	sglog.Debug("Phone:%s", data.Phone)
	sglog.Debug("RequireDt:%s", data.RequireDt.NormalString())
	sglog.Debug("RandomNum:%s", data.RandomNum)
	sglog.Debug("AnswerTimes:%d", data.AnswerTimes)
	sglog.Debug("status:%d", data.Status)
	sglog.Debug("requireTimes:%d", data.RequireTimes)
	sglog.Debug("=======end===========")
}

type SMSData struct {
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
}

type SWxOpenid struct {
	Code   string
	Openid string
	Time   *sgtime.DateTime
}

type SecureWxOpenid struct {
	Data map[string](map[string]*SWxOpenid)
	Lock sync.RWMutex
}
