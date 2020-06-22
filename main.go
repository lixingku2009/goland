package main

import (
	"world/world"
	_ "net/http/pprof"
	"runtime"
	"net/http"
	"strconv"
)

const MaxWorld  = 1

func main()  {
	WorldMgr := make([]*world.World, 0)
	for i := 0; i < MaxWorld; i++ {
		w := new(world.World)
		w.Init(i)
		WorldMgr = append(WorldMgr, w)
	}

	for _, w := range WorldMgr {
		go w.WorldUpdateLoop()
	}

	startProfile(8888)

	select {
	}
}

func startProfile(port int) {
	runtime.SetBlockProfileRate(1)
	//远程获取pprof数据
	go func() {
		_ = http.ListenAndServe("localhost:" + strconv.Itoa(port), nil)
	}()
}