package logx

import (
	"fmt"
	"log"
	"myapi"
	"os"
	"runtime"
	"strconv"
)

type Logx struct {
	Db myapi.DBBOX
}

//Log 指定文件夹
func (L *Logx) Log(s string) {

	if L.Db == nil {
		if s != "" {
			xp := L.logxp(2)
			file, err := os.OpenFile("defaltLog.log", os.O_APPEND|os.O_CREATE, 666)
			if err != nil {
				fmt.Println(xp+"***Log fail:", err)
			}
			defer file.Close()

			logger := log.New(file, xp, log.LstdFlags) // 日志文件格式:log包含时间及文件行数
			log.Println(s)
			logger.Println(s)
		}
		return
	}

	err := L.Db.SecurityRunSQL("insert into log values (?,default)", s)
	if err != nil {
		log.Println(s, "***write log error***", err)
	}
}
func (L *Logx) logxp(i int) string {
	_, f, l, ok := runtime.Caller(i)
	if !ok {
		f = "???"
		l = 0
	}
	out := f + ":" + strconv.Itoa(l) + "||"
	return out
}
