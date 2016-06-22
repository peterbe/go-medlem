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

type UsersForm struct {
	Emails []string `form:"email"`
}

type UsersJSON struct {
	Emails []string `json:"email"`
}

// func logger(format string, a ...interface{}) {
//     iris.Logger.Infof(format, a)
// }

func isStaff(ctx *iris.Context) {
	jsonUsers := UsersJSON{}
	var emails []string
	jsonErr := ctx.ReadJSON(&jsonUsers)

	if jsonErr != nil {
		// XXX return a BadRequest if the attempt really was JSON
		log.Println("Error when reading JSON body: " + jsonErr.Error())
	} else {
		// log.Println("Users", jsonUsers)
		emails = jsonUsers.Emails
	}
	if jsonErr != nil || len(emails) == 0 {
		// Read form form data or query string
		formUsers := UsersForm{}
		formErr := ctx.ReadForm(&formUsers)
		if formErr != nil {
			log.Println("Error when reading form: " + formErr.Error())
		} else {
			emails = formUsers.Emails
		}
	}

	if len(emails) == 0 {
		// ctx.Write("No emails supplied. See docs\n")
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"error": "No emails supplied. See docs",
		})
		ctx.SetStatusCode(iris.StatusBadRequest) // 400
	} else {
		// log.Println("Users", emails)
		results := make(map[string]bool)
		for _, email := range emails {
			// log.Println("EMAIL:", email)
			results[email] = false
		}

		// ctx.JSON(iris.StatusOK, iris.Map{
		// 	"emails":  emails,
		// })
		ctx.JSON(iris.StatusOK, results)
	}

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	debug := os.Getenv("DEBUG") == "true"
	if debug {
		// log.Println("Running in debug mode")
		iris.Logger.Infof("Running in debug mode")
		iris.Config.Render.Template.IsDevelopment = true
	}

	api := iris.New()
	api.Get("/helloworld", helloworld)
	api.Get("/staff", isStaff)
	api.Post("/staff", isStaff)
	api.Get("/", index)
	api.Listen("0.0.0.0:" + port)
}
