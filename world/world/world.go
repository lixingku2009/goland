package world

import (
	"time"
	"go.uber.org/atomic"
	"sync"
	"world/gtime"
)

// 服务器状态
type worldStatus int32
const (
	Stopped			worldStatus = 0
	Running			worldStatus = 1
	Paused			worldStatus = 2
)

type World struct {
	id				int
	status			atomic.Int32
	loopCounter		atomic.Uint64
	players			sync.Map				// key:playerId
	playerLen		int						// 玩家人数
	chFunc			chan func()				// 处理队列
	gameTime		*gtime.GameTime			// 游戏时间
	updateTime		*gtime.WorldUpdateTime	//
}

const (
	SleepConst		= 20 * time.Millisecond		// 50帧及20ms一帧
	PausedTime		= time.Second				// 游戏暂停多久重新查询一遍
	MaxMsgLen		= 100						// 消息队列长度
	PushMsgTimeOut	= time.Second				// 推送消息的超时时间
)

func (w *World) Init(id int) {
	w.id = id
	w.chFunc = make(chan func(), MaxMsgLen)
	w.gameTime = new(gtime.GameTime)
	w.gameTime.Init()
	w.updateTime = new(gtime.WorldUpdateTime)
	w.updateTime.Init()
	w.Running()
}

// 游戏世界的维护
func (w *World) WorldUpdateLoop() {
	var (
		realCurrTime 		time.Time					// update前
		realPrevTime = 		time.Now().Add(-SleepConst)	// update后
		executionTimeDiff	time.Duration				// update花费时长
		diff				time.Duration				// world需要update的时长
	)

	for {
		if w.isStopped() {
			w.close()
			break
		}
		if w.isPaused() {
			time.Sleep(PausedTime)
			realPrevTime = time.Now().Add(-executionTimeDiff)
			continue
		}

		w.loopCounter.Inc()
		realCurrTime = time.Now()
		diff = realCurrTime.Sub(realPrevTime)
		w.update(diff)											// 更新游戏世界
		realPrevTime = realCurrTime
		executionTimeDiff = time.Now().Sub(realCurrTime)

		// 获取更新世界花了多长时间，如果更新花费的时间少于SleepConst，请等待SleepConst世界更新时间
		if executionTimeDiff < SleepConst {
			time.Sleep(SleepConst - executionTimeDiff)
		}
	}
}

func (w *World) close() {
	// todo something
}

func (w *World) Stop() {
	w.status.Swap(int32(Stopped))
}

func (w *World) Running() {
	w.status.Swap(int32(Running))
}

func (w *World) Paused() {
	w.status.Swap(int32(Paused))
}

func (w *World) isStopped() bool {
	return w.status.Load() == int32(Stopped)
}

func (w *World) isPaused() bool {
	return w.status.Load() == int32(Paused)
}

func (w *World) update(diff time.Duration) {
	w.gameTime.UpdateGameTime()
	// currentGameTime := w.gameTime.GetGameTime()
	w.updateTime.UpdateWithDiff(diff)
	w.updateTime.RecordUpdateTime(w.gameTime.GetGameTime(), diff, w.playerLen)

	// 先处理消息
	w.handFunc()

	// 再游戏自身的更新
	w.updateWorld()
}


func (w *World) updateWorld() {
	// todo do something
}

func (w *World) handFunc() {
	mLen := len(w.chFunc)
	for i := 0; i < int(mLen); i++ {
		f := <-w.chFunc
		f()
	}
}

// !如果是RPG游戏,这个f绝对要避免一切阻塞, 里面的日志要放到一个日志协程
func (w *World) PushFunc(f func()) bool {
	ti := time.NewTimer(PushMsgTimeOut)
	select {
	case w.chFunc<-f:
	case <-ti.C:
		return false
	}

	return true
}