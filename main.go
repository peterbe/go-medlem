package main

import (
    "log"
    "fmt"
    "os"
    "github.com/kataras/iris"
    )

func index(ctx *iris.Context){
   // ctx.Render("index.html", struct { Name string }{ Name: "iris" })
   // maybe ctx.Render("index.html", nil)
   var context struct{}
   ctx.Render("index.html", context)
}

func helloworld(ctx *iris.Context){
   ctx.Write("Hi %s\n", "world")
}


func main() {
    port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	} else {
        fmt.Sprintf("PORT=%s\n", port)
    }

    // iris.Config.Render.Template.IsDevelopment = true

    api := iris.New()
    api.Get("/helloworld", helloworld)
    api.Get("/", index)
    api.Listen("0.0.0.0:" + port)
}
