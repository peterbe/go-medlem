package main

import (
	"github.com/kataras/iris"
	"log"
	"os"
)

func index(ctx *iris.Context) {
	// ctx.Render("index.html", struct { Name string }{ Name: "iris" })
	// maybe ctx.Render("index.html", nil)
	var context struct{}
	ctx.Render("index.html", context)
}

func helloworld(ctx *iris.Context) {
	ctx.Write("Hi %s\n", "world")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	debug := os.Getenv("DEBUG") == "true"
	if debug {
		log.Println("Running in debug mode")
		iris.Config.Render.Template.IsDevelopment = true
	}

	api := iris.New()
	api.Get("/helloworld", helloworld)
	api.Get("/", index)
	api.Listen("0.0.0.0:" + port)
}
