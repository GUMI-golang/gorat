package fwrat

import "sync"

var (
	lock = new(sync.Mutex)
	islock = false
)