// 这个文件是记录游戏更新耗时的模块
// 看似不需要用到
package gtime

import (
	"time"
	"math"
	"log"
)

const (
	avgDiffCount = 500
)

type UpdateTime struct {
	updateTimeDataTable				[avgDiffCount]time.Duration		// 最近avgDiffCount帧的更新时长
	averageUpdateTime				time.Duration					// 最后一组(avgDiffCount次)帧平均时长
	totalUpdateTime					time.Duration					// 最后一组(avgDiffCount次)帧总时长
	updateTimeTableIndex			int								// 记录当前用到了第几帧时长
	maxUpdateTime					time.Duration					// 最大更新时长
	maxUpdateTimeOfLastTable		time.Duration					// 最后一组(avgDiffCount次)帧的最大时长
	maxUpdateTimeOfCurrentTable		time.Duration					// 当前一组(avgDiffCount次)帧的最大时长
	recordedTime					time.Time
}

func (u *UpdateTime) InitTime() {
}

func (u *UpdateTime) GetAverageUpdateTime() time.Duration {
	return u.averageUpdateTime
}

// 加权平均数
func (u *UpdateTime) GetTimeWeightedAverageUpdateTime() time.Duration {
	var sum time.Duration
	var weightSum time.Duration
	for _, diff := range u.updateTimeDataTable {
		sum += diff * diff
		weightSum += diff
	}

	if weightSum == 0 {
		return weightSum
	}

	return sum / weightSum
}

func (u *UpdateTime) GetMaxUpdateTime() time.Duration {
	return u.maxUpdateTime
}

// 这一组的最大时间
func (u *UpdateTime) GetMaxUpdateTimeOfCurrentTable() time.Duration {
	return time.Duration(math.Max(float64(u.maxUpdateTimeOfCurrentTable), float64(u. maxUpdateTimeOfLastTable)))
}

// 上次耗时
func (u *UpdateTime) GetLastUpdateTime() time.Duration {
	if u.updateTimeDataTable[u.updateTimeTableIndex] != 0 {
		return u.updateTimeDataTable[u.updateTimeTableIndex - 1]
	}

	return u.updateTimeDataTable[avgDiffCount - 1]
}

func (u *UpdateTime) UpdateWithDiff(diff time.Duration)  {
	u.totalUpdateTime = u.totalUpdateTime - u.updateTimeDataTable[u.updateTimeTableIndex] + diff
	u.updateTimeDataTable[u.updateTimeTableIndex] = diff

	if diff > u.maxUpdateTime {
		u.maxUpdateTime = diff
	}

	if diff > u.maxUpdateTimeOfCurrentTable {
		u.maxUpdateTimeOfCurrentTable = diff
	}

	u.updateTimeTableIndex++

	if u.updateTimeTableIndex >= avgDiffCount {
		u.updateTimeTableIndex = 0
		u.maxUpdateTimeOfLastTable = u.maxUpdateTimeOfCurrentTable
		u.maxUpdateTimeOfCurrentTable = 0
	}

	if u.updateTimeDataTable[avgDiffCount - 1] != 0 {
		u.averageUpdateTime = u.totalUpdateTime / avgDiffCount
		//log.Println(u.averageUpdateTime, u.totalUpdateTime, avgDiffCount)
	} else if u.updateTimeTableIndex - 1 != 0 {
		u.averageUpdateTime = time.Duration(int(u.totalUpdateTime) / u.updateTimeTableIndex)
		//log.Println(u.averageUpdateTime, u.totalUpdateTime, u.updateTimeTableIndex)
	}
}

func (u *UpdateTime) RecordUpdateTimeReset() {
	u.recordedTime = time.Now()
}

func (u *UpdateTime) recordUpdateTimeDuration(text string, minUpdateTime time.Duration) {
	thisTime := time.Now()
	diff := thisTime.Sub(u.recordedTime)

	if diff > minUpdateTime {
		log.Printf("Recover Update Time of %v: %v.", text, diff)
	}

	u.recordedTime = thisTime
}

type WorldUpdateTime struct {
	UpdateTime
	recordUpdateTimeInterval		time.Duration
	recordUpdateTimeMin				time.Duration
	lastRecordTime					time.Time
}

func (w *WorldUpdateTime) Init() {
	w.InitTime()
}

// todo 读配置
func (w *WorldUpdateTime) LoadFromConfig() {
	//w.recordUpdateTimeInterval = sConfigMgr->GetIntDefault("RecordUpdateTimeDiffInterval", 60000)
	//w.recordUpdateTimeMin = sConfigMgr->GetIntDefault("MinRecordUpdateTimeDiff", 100)
	w.recordUpdateTimeInterval = 60000 * time.Millisecond
	w.recordUpdateTimeMin = 100 * time.Millisecond
}

func (w *WorldUpdateTime) SetRecordUpdateTimeInterval(t time.Duration) {
	w.recordUpdateTimeInterval = t
}

func (w *WorldUpdateTime) RecordUpdateTime(gameTimeMs time.Time, diff time.Duration, sessionCount int) {
	if w.recordUpdateTimeInterval > 0 && diff > w.recordUpdateTimeMin {
		if gameTimeMs.Sub(w.lastRecordTime) > w.recordUpdateTimeInterval {
			log.Printf("Update time diff: %v. Players online: %v.", w.GetAverageUpdateTime(), sessionCount)
			w.lastRecordTime = gameTimeMs
		}
	}
}

func (w *WorldUpdateTime) RecordUpdateTimeDuration(text string){
	w.recordUpdateTimeDuration(text, w.recordUpdateTimeMin)
}
