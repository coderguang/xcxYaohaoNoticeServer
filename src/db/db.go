package yaohaoNoticeDb

import (
	"database/sql"
	"strconv"
	yaohaoNoticeData "xcxYaohaoNoticeServer/src/data"
	yaohaoNoticeDef "xcxYaohaoNoticeServer/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgmysql"
)

var globalmysqldb *sql.DB
var globalmysqlStmt *sql.Stmt
var globalmysqlRequireStmt *sql.Stmt

func InitDbConnection(user string, pwd string, url string, port string, dbname string) {
	conn, err := sgmysql.Open(user, pwd, url, port, dbname, "utf8")
	if err != nil {
		sglog.Fatal("connection to db %s error,%e", url, err)
	}
	globalmysqldb = conn

	dataSql := "replace into " + yaohaoNoticeData.GetTableName() + "(token_id,name,title,card_type,code,phone,end_dt,tips,renew_times,status,notice_times) values(?,?,?,?,?,?,?,?,?,?,?)"

	globalmysqlStmt, err = globalmysqldb.Prepare(dataSql)
	if err != nil {
		sglog.Fatal("init replace sql error,%s", err)
	}

	requireSql := "replace into " + yaohaoNoticeData.GetRequireDataTableName() + "(token_id,title,name,require_time,final_login,share_times,tips) values(?,?,?,?,?,?,?)"

	globalmysqlRequireStmt, err = globalmysqldb.Prepare(requireSql)
	if err != nil {
		sglog.Fatal("init require replace sql error,%s", err)
	}

	sglog.Info("InitDbConnection complete")
}

func LoadNoticeDataFromDb() {
	sqlStr := "select * from " + yaohaoNoticeData.GetTableName()
	rows, rowsErr := globalmysqldb.Query(sqlStr)
	if rowsErr != nil {
		sglog.Error("init notice data error,err=%e", rowsErr)
		return
	}
	defer rows.Close()
	initSize := 0
	for rows.Next() {
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		_ = rows.Scan(scanArgs...)

		data := new(yaohaoNoticeDef.SData)

		for i, col := range values {
			if col != nil {
				field := columns[i]
				fieldValue := string(col.([]byte))
				switch field {
				case "token_id":
					data.Token = fieldValue
				case "name":
					data.Name = fieldValue
				case "title":
					data.Title = fieldValue
				case "code":
					data.Code = fieldValue
				case "phone":
					data.Phone = fieldValue
				case "end_dt":
					data.EndDt = sgtime.New()
					data.EndDt.Parse(fieldValue, sgtime.FORMAT_TIME_NORMAL)
				case "tips":
					data.Desc = fieldValue
				case "renew_times":
					data.RenewTimes, _ = strconv.Atoi(fieldValue)
				case "status":
					data.Status = fieldValue
				case "notice_times":
					data.NoticeTimes, _ = strconv.Atoi(fieldValue)
				case "card_type":
					data.CardType, _ = strconv.Atoi(fieldValue)
				}
			}

		}
		initSize++
		yaohaoNoticeData.AddNoticeData(data)

	}

	sglog.Info("load data from db complete,size=%d", initSize)
}

func InsertOrUpdateData(data *yaohaoNoticeDef.SData) {
	globalmysqlStmt.Exec(data.Token, data.Name, data.Title, data.CardType, data.Code, data.Phone, data.EndDt.NormalString(), data.Desc, data.RenewTimes, data.Status, data.NoticeTimes)
}

func LoadNoticeRequireDataFromDb() {
	sqlStr := "select * from " + yaohaoNoticeData.GetRequireDataTableName()
	rows, rowsErr := globalmysqldb.Query(sqlStr)
	if rowsErr != nil {
		sglog.Error("init require data error,err=%s", rowsErr)
		return
	}
	defer rows.Close()
	initSize := 0
	for rows.Next() {
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		_ = rows.Scan(scanArgs...)

		data := new(yaohaoNoticeDef.SNoticeRequireData)

		for i, col := range values {
			if col != nil {
				field := columns[i]
				fieldValue := string(col.([]byte))
				switch field {
				case "token_id":
					data.Openid = fieldValue
				case "name":
					data.Name = fieldValue
				case "title":
					data.Title = fieldValue
				case "require_time":
					data.RequireTimes, _ = strconv.Atoi(fieldValue)
				case "final_login":
					data.FinalLogin = sgtime.New()
					data.FinalLogin.Parse(fieldValue, sgtime.FORMAT_TIME_NORMAL)
				case "share_times":
					data.ShareTimes, _ = strconv.Atoi(fieldValue)
				case "desc":
					data.Desc = fieldValue
				}
			}

		}
		initSize++
		yaohaoNoticeData.AddOrUpdateNoticeRequireData(data)

	}

	sglog.Info("load require data from db complete,size=%d", initSize)
}

func InsertOrUpdateRequireData(data *yaohaoNoticeDef.SNoticeRequireData) {
	globalmysqlStmt.Exec(data.Openid, data.Title, data.Name, data.RequireTimes, data.FinalLogin.NormalString(), data.ShareTimes, data.Desc)
}
