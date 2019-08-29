package yaohaoNoticeConfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	yaohaoNoticeDef "xcxYaohaoNoticeServer/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

func ReadConfig(configfile string) *yaohaoNoticeDef.Config {
	config, err := ioutil.ReadFile(configfile)
	if err != nil {
		sglog.Fatal("read config error")
		os.Exit(1)
	}
	t := new(yaohaoNoticeDef.Config)
	p := &t
	err = json.Unmarshal([]byte(config), p)
	if err != nil {
		sglog.Fatal("parse config error")
		os.Exit(1)
	}
	return t
}
