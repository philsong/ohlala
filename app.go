package main

import (
	"github.com/QLeelulu/goku"
	"github.com/philsong/ohlala/golink"
	_ "github.com/philsong/ohlala/golink/controllers"       // notice this!! import controllers
	_ "github.com/philsong/ohlala/golink/controllers/admin" // notice this!! import controllers
	"github.com/philsong/ohlala/golink/middlewares"
	"log"
)

func main() {

	rt := &goku.RouteTable{Routes: golink.Routes}
	middlewares := []goku.Middlewarer{
		new(middlewares.UtilMiddleware),
		new(middlewares.ConfessMiddleware),
	}
	s := goku.CreateServer(rt, middlewares, golink.Config)
	goku.Logger().Logln("Server start on", s.Addr)

	log.Fatal(s.ListenAndServe())
}
