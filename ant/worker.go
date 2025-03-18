package ant

import "time"

type worker interface {
	run()
	finish()
	lastUsedTime() time.Time
	setLastUsedTime(t time.Time)
	inputFunc(func())
	inputArg(any)
}
