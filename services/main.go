package main

import (
	"fmt"
	"runtime"
	"time"

	"./goo"
)

func main() {
	g, _ := goo.New()
	go func() {
		for {
			time.Sleep(time.Second * 10)
			fmt.Println(runtime.NumGoroutine())
		}
	}()
	g.Start(":5354")

}
