package safe

import (
	"log"
	"runtime/debug"
)

func Go(goroutine func()) {
	GoWithRecover(goroutine, func(err interface{}) {
		log.Printf("Error in Go routine: %s\nStack: %s\n", err, debug.Stack())
	})
}

func GoWithRecover(goroutine func(), customRecover func(err interface{})) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				customRecover(err)
			}
		}()
		goroutine()
	}()
}
