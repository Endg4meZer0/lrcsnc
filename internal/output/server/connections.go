package server

import (
	"net"
	"sync"
)

type Connections struct {
	M             sync.Mutex
	Map           map[uint]net.Conn
	FirstFreeSlot uint
}

func (cs *Connections) FindFirstAvailableSlot() {
	cs.M.Lock()
	defer cs.M.Unlock()

	var i uint = 0
	for i = range uint(len(cs.Map)) {
		if _, ok := cs.Map[i]; !ok {
			cs.FirstFreeSlot = i
			return
		}
	}
	cs.FirstFreeSlot = i + 1
}
