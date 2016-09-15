package main

import (
	"github.com/gamexg/gotool/ozzolog"
	"github.com/gamexg/ozzo-log"
)

var llog = ozzolog.GetLogger("main")

func main() {
	lConf := ozzolog.Config{
		LogConsoleLevel: int(log.LevelDebug),
	}

	ozzolog.Open(&lConf)
	defer ozzolog.Close()

	llog.Error("err")

}
