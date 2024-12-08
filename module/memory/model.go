package memory

import (
	"golang.org/x/sys/windows"
	"sync"
)

type locksMap struct {
	locks map[uintptr]*bool
	lock  *sync.RWMutex
}

type Ctx struct {
	handle  windows.Handle
	uintPtr uintptr
	maps    locksMap
}
