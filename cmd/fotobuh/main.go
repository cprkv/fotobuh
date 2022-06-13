package main

import (
	"fotobuh/lib"
	"fotobuh/lib/db"
	"fotobuh/lib/server"
	"log"
)

func unwrap(err error) {
	if err != nil {
		log.Panicf("[unwrap error] %v", err)
	}
}

func main() {
	unwrap(lib.Config.Init())
	unwrap(db.Context.Init())
	server.StartApp()
}
