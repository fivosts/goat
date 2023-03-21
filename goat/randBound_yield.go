package goat

import (
	_ "fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type SharedCounter struct {
	cnt int
	sync.Mutex
}

var bound SharedCounter

func randBound_yield() {
	thr, err := strconv.Atoi(os.Getenv("GOATRSBOUND"))
	if err != nil || thr < 1 {
		return
	}

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 1 {
		bound.Lock()
		if bound.cnt < thr {
			bound.cnt++
			bound.Unlock()
			runtime.Gosched()
		} else {
			bound.Unlock()
		}
	}
}
