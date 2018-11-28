package errors

import (
	"runtime"

	"github.com/matchvs/gameServer-go/src/log"

	"github.com/davecgh/go-spew/spew"
)

// 产生panic时的调用栈打印
func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		log.LogE("recover:%v", x)
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			log.LogE("frame %v:[func:%v,file:%v,line:%v]", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}

		for k := range extras {
			log.LogE("EXRAS#%v DATA:%v", k, spew.Sdump(extras[k]))
		}
	}
}
