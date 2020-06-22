// 这个模块是用来自己维护服务器的tick
// 频繁调用time.Now()的服务可以用
// time.Now()代用100万次大约10ms,效率非常高
// 而自己维护的时间有精度问题
// 自己斟酌是否有自己维护时间的必要
// 经过测试,似乎没有必要
package stime

import (
	"time"
	"sync"
)

const (
	loopDuration	= 2 * time.Millisecond		// 时间精度,最短请设定2ms,再低时间的精度会急剧下降
	resetDuration	= time.Second				// 校准间隔
)

type serverTime struct {
	t				time.Time
	loopDuration	time.Duration
	sync.RWMutex
}

var sTime	= &serverTime{}

func Init()  {
	go loop()
}

func loop() {
	sTime.t = time.Now()
	for {
		sTime.Lock()
		if sTime.loopDuration >= resetDuration {
			sTime.t = time.Now()
			sTime.loopDuration = 0
		} else {
			sTime.t.Add(loopDuration)
			sTime.loopDuration += loopDuration
		}
		sTime.Unlock()
		time.Sleep(loopDuration)
	}
}

func Now() time.Time {
	//return time.Now()
	sTime.RLock()
	defer sTime.RUnlock()
	return sTime.t
}
