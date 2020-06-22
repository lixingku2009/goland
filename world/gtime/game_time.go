package gtime

import (
	"time"
)

type GameTime struct {
	startTime		time.Time
	gameTime		time.Time
}

func (g *GameTime) Init() {
	g.startTime = time.Now()
	g.gameTime = g.startTime
}

func (g *GameTime) UpdateGameTime() {
	g.gameTime = time.Now()
}

func (g *GameTime) GetGameTime() time.Time {
	return g.gameTime
}

func (g *GameTime) GetStartTime() time.Time {
	return g.startTime
}

func (g *GameTime) GetUptime() time.Duration {
	return g.gameTime.Sub(g.startTime)
}
