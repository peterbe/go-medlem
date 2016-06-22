package main

import (
	"errors"
	"github.com/kataras/iris"
	"log"
	"net/mail"
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

// private function
func getEmails(ctx *iris.Context) ([]string, error) {
	// var emails []string
	// var error nil
	jsonUsers := UsersJSON{}
	var emails []string
	jsonErr := ctx.ReadJSON(&jsonUsers)

	log.Println(ctx.RequestHeader("Content-Type"))
	if jsonErr != nil {
		// XXX return a BadRequest if the attempt really was JSON
		// log.Println("Error when reading JSON body: " + jsonErr.Error())

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

	var parsedEmails []string
	for _, email := range emails {
		e, err := mail.ParseAddress(email)
		if err != nil {
			return nil, errors.New(
				"Not a valid email address: " + email + " (" + err.Error() + ")",
			)
		}
		// log.Println(e.Name, e.Address)
		// parsedEmails.append(e.Address)
		parsedEmails = append(parsedEmails, e.Address)
	}

	return parsedEmails, nil

}

func IsStaff(ctx *iris.Context) {
	emails, err := getEmails(ctx)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"error": "No emails supplied. See docs",
		})
		ctx.SetStatusCode(iris.StatusBadRequest) // 400
	} else if len(emails) == 0 {
		// ctx.Write("No emails supplied. See docs\n")
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"error": "No emails supplied. See docs",
		})
		ctx.SetStatusCode(iris.StatusBadRequest) // 400
		return
	} else {
		results := make(map[string]bool)
		for _, email := range emails {
			// log.Println("EMAIL:", email)
			results[email] = false
		}
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
		log.Println("Running in debug mode")
		// iris.Logger.Infof("Running in debug mode")
		iris.Config.Render.Template.IsDevelopment = true
	}

	api := iris.New()
	api.Get("/helloworld", helloworld)
	api.Get("/staff", IsStaff)
	api.Post("/staff", IsStaff)
	api.Get("/contribute.json", func(ctx *iris.Context) {
		ctx.ServeFile("./contribute.json", false)
	})
	api.Static("/node_modules", "./node_modules/", 1)
	api.Static("/static", "./static/", 1)
	api.Get("/", index)
	api.Listen("0.0.0.0:" + port)
}
