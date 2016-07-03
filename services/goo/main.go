package main

import (
	"os"

	"./goo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5354"
	}
	port = ":" + port

	g, _ := goo.New()
	g.Run(port)
}
