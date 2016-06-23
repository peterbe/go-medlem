package main

import (
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"go.mozilla.org/mozldap"
	"log"
	"net/mail"
	"os"
	// "crypto/tls"
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
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"error": "No emails supplied. See docs",
		})
		ctx.SetStatusCode(iris.StatusBadRequest) // 400
		return
	} else {

		client, err := mozldap.NewTLSClient(
			ldapURI,
			ldapUsername,
			ldapPassword,
			"/Users/peterbe/dev/MOZILLA/MEDLEM/ldap-bind/medlem/ldapproxy-medlem.crt",
			"/Users/peterbe/dev/MOZILLA/MEDLEM/ldap-bind/medlem/ldapproxy-medlem.key",
			// "/Users/peterbe/dev/MOZILLA/MEDLEM/ldap-bind/medlem/ldapproxy-medlem.csr",
			"",
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		results := make(map[string]bool)
		for _, email := range emails {
			results[email] = false
		}

		mailFilter := ""
		for _, email := range emails {
			// XXX need escaping on the email. See how python-ldap does it in ldap.filter.filter_format
			mailFilter += fmt.Sprintf(
				// "(&(|(mail=%s)(emailAlias=%s))(objectClass=mozComPerson))",
				"(&(mail=%s)(objectClass=mozComPerson))",
				email, email,
			)
		}
		mailFilter = fmt.Sprintf(
			"(|%s)", mailFilter,
		)
		log.Println("mailFilter:", mailFilter)
		entries, searchErr := client.Search(
			"",
			mailFilter,
			// nil,  // use this to see all available/possible columns
			// []string{"mail", "employeeType", "emailAlias"},
			[]string{"mail"},
		)
		if searchErr != nil {
			log.Fatal(searchErr)
		}
		// log.Println(entries)
		for _, entry := range entries {
			for _, attr := range entry.Attributes {
				log.Println(attr.Name, ":", attr.Values)
				if attr.Name == "mail" {
					results[attr.Values[0]] = true
				}
			}
		}
		ctx.JSON(iris.StatusOK, results)
	}
}

var (
	ldapURI      string
	ldapUsername string
	ldapPassword string
)

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

	ldapURI = os.Getenv("LDAP_URI")
	ldapUsername = os.Getenv("LDAP_USERNAME")
	ldapPassword = os.Getenv("LDAP_PASSWORD")
	if ldapURI == "" {
		log.Fatal("$LDAP_URI must be set")
	}
	if ldapUsername == "" {
		log.Fatal("$LDAP_USERNAME must be set")
	}
	if ldapPassword == "" {
		log.Fatal("$LDAP_PASSWORD must be set")
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
