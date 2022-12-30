package core

import "runtime/debug"

func PrintPanicStack() {
	if r := recover(); r != nil {
		Log.Errorf("%v: %s", r, debug.Stack())
	}
}
