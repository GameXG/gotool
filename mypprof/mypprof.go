package mypprof

import "time"
import (
	"sync"
	"sync/atomic"

	"os"

	"fmt"

	"encoding/csv"

	"strconv"

	log "bitbucket.org/jack/log4go"
)

type FuncTypeS struct {
	Min   int64
	Max   int64
	Conut int64
	m     sync.Mutex
	Total float64
}

var Enable = int32(0)

var FuncTimeM sync.RWMutex
var FuncTimeMap map[string]*FuncTypeS

func init() {
	FuncTimeMap = make(map[string]*FuncTypeS)
	go loop()
}

func loop() {
	tick := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-tick.C:
			if atomic.LoadInt32(&Enable) == 0 {
				continue
			}

			func() {
				f, err := os.Create(fmt.Sprintf("FuncTime-%v.log", time.Now().Format("2006-01-02 15-04-05")))
				if err != nil {
					log.Error("创建 FuncTime 文件错误，%v", err)
				}
				defer f.Close()

				FuncTimeM.Lock()
				old := FuncTimeMap
				FuncTimeMap = make(map[string]*FuncTypeS)
				FuncTimeM.Unlock()

				w := csv.NewWriter(f)
				defer w.Flush()
				w.Write([]string{"函数名", "最小值", "最大值", "平均值", "次数"})

				for name, value := range old {
					value.m.Lock()
					if err := w.Write([]string{name,
						strconv.FormatInt(value.Min, 10),
						strconv.FormatInt(value.Max, 10),
						strconv.FormatInt(int64(value.Total / float64(value.Conut)), 10),
						strconv.FormatInt(value.Conut, 10)}); err != nil {
						log.Error("写 FuncTime 错误，%v", err)
					}
					value.m.Unlock()
				}

				log.Info("写入FuncTime完成，宫写入 %v 条。", len(old))
			}()
		}
	}
}

// 软件内部的耗时
// 单位纳秒
func LogFuncTime(name string, message string, start time.Time, timeLimit int64) {
	if atomic.LoadInt32(&Enable) == 0 {
		return
	}

	end := time.Now()

	FuncTimeM.RLock()
	s, _ := FuncTimeMap[name]
	FuncTimeM.RUnlock()

	if s == nil {
		FuncTimeM.Lock()
		if s, _ = FuncTimeMap[name]; s == nil {
			s = &FuncTypeS{}
			FuncTimeMap[name] = s
		}
		FuncTimeM.Unlock()
	}

	t := end.UnixNano() - start.UnixNano()

	if t > atomic.LoadInt64(&s.Max) {
		atomic.StoreInt64(&s.Max, t)
	}

	min := atomic.LoadInt64(&s.Min)
	if t < min || min == 0 {
		atomic.StoreInt64(&s.Min, t)
	}

	s.m.Lock()
	s.Total += float64(t)
	s.m.Unlock()

	atomic.AddInt64(&s.Conut, 1)

	if t > timeLimit {
		log.Error("函数 %v 执行超时，%v > %v ，%v", name, t, timeLimit, message)
	}
}

// 需要统计到 TAP 网卡到软件的耗时
// 每秒打印一个收到的包
// 非线程安全
type PrintPack struct {
	printTapPackTime int64
	name             string
}

func NewPrintPack(name string) *PrintPack {
	return &PrintPack{
		name: name,
	}
}

func (p *PrintPack) PrintTapPack(pack []byte) {
	if atomic.LoadInt32(&Enable) == 0 {
		return
	}
	now := time.Now()
	if p.printTapPackTime != now.Unix() {
		p.printTapPackTime = now.Unix()
		log.Info("[PrintTapPack-%v] %v %v", p.name, now, pack)
	}
}

// 允许调试客户端关闭压缩

// 调试客户端发出的时间
