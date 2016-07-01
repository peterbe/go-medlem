package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/kataras/iris"
	"go.mozilla.org/mozldap"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"sort"
)

func index(ctx *iris.Context) {
	// ctx.Render("index.html", struct { Name string }{ Name: "iris" })
	// maybe ctx.Render("index.html", nil)
	var context struct{}
	ctx.Render("index.html", context)
}

func helloworld(ctx *iris.Context) {
	log.Println(sampleMultiline)
	filename := repackageAsFilepath(sampleMultiline)
	defer os.Remove(filename) // clean up
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	ctx.Write("Hi %s\n", content)
}

type UsersForm struct {
	Emails []string `form:"email"`
}

type UsersJSON struct {
	Emails []string `json:"email"`
}

// private function
func getEmails(ctx *iris.Context) ([]string, error) {
	// var emails []string
	// var error nil
	jsonUsers := UsersJSON{}
	var emails []string
	jsonErr := ctx.ReadJSON(&jsonUsers)

	// log.Println(ctx.RequestHeader("Content-Type"))
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

func repackageAsFilepath(thing string) string {
	/* This function will always return a valid file path.
	If the supplied parameter "thing" is not already a file, but a string,
	we'll take its content and put it into a temp file. */
	_, err := os.Stat(thing)
	if err == nil {
		// it was already a valid file
		return thing
	}
	tmpfile, err := ioutil.TempFile("", "thing")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(thing)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile.Name()
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

		// perhaps we have it cached
		sort.Strings(emails)
		emailHash := md5.New()
		for _, email := range emails {
			io.WriteString(emailHash, email)
		}
		cacheKey := "staff:" + hex.EncodeToString(emailHash.Sum(nil))
		if jsonString, err := cache.Get([]byte(cacheKey)); err == nil {
			results := make(map[string]bool)
			if err := json.Unmarshal(jsonString, &results); err == nil {
				ctx.JSON(iris.StatusOK, results)
				return
			}
		}

		// ldapCertFile and ldapKeyFile needs to be actual file paths
		// to real files on the filesystem going into mozldap.NewTLSClient
		// but they might have arrived to us in the form of multiline
		// strings. E.g. as environment variables set in Heroku.
		ldapCertFilePath := repackageAsFilepath(ldapCertFile)
		if ldapCertFilePath != ldapCertFile {
			defer os.Remove(ldapCertFilePath)
		}
		ldapKeyFilePath := repackageAsFilepath(ldapKeyFile)
		if ldapKeyFilePath != ldapKeyFile {
			defer os.Remove(ldapKeyFilePath)
		}
		client, err := mozldap.NewTLSClient(
			ldapURI,
			ldapUsername,
			ldapPassword,
			ldapCertFilePath,
			ldapKeyFilePath,
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
				"(mail=%s)",
				// "(&(mail=%s)(!(employeeType=DISABLED)))",
				// email, email
				email,
			)
		}
		mailFilter = fmt.Sprintf(
			"(&(|%s)(objectClass=mozComPerson))", mailFilter,
		)
		// log.Println("mailFilter:", mailFilter)
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
		if jsonString, err := json.Marshal(results); err == nil {
			cache.Set([]byte(cacheKey), []byte(jsonString), 60*5) // 5 min caching
		}
		ctx.JSON(iris.StatusOK, results)
	}
}

var (
	ldapURI         string
	ldapUsername    string
	ldapPassword    string
	ldapCertFile    string
	ldapKeyFile     string
	sampleMultiline string
	cacheSize       = 10 * 1024 * 1014
	cache           *freecache.Cache
)

func main() {
	sampleMultiline = os.Getenv("SAMPLE_MULTILINE")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	cache = freecache.NewCache(cacheSize)

	ldapURI = os.Getenv("LDAP_URI")
	ldapUsername = os.Getenv("LDAP_USERNAME")
	ldapPassword = os.Getenv("LDAP_PASSWORD")
	ldapCertFile = os.Getenv("LDAP_CERT_FILE")
	ldapKeyFile = os.Getenv("LDAP_KEY_FILE")
	requiredKeys := []string{
		"LDAP_URI",
		"LDAP_USERNAME",
		"LDAP_PASSWORD",
		"LDAP_CERT_FILE",
		"LDAP_KEY_FILE",
	}
	for _, key := range requiredKeys {
		if os.Getenv(key) == "" {
			log.Fatal(fmt.Sprintf("$%v must be set", key))
		}
	}

	api := iris.New()

	debug := os.Getenv("DEBUG")
	log.Println(debug)
	if debug == "true" {
		log.Println("Running in debug mode")
		// iris.Logger.Infof("Running in debug mode")
		api.Config.Render.Template.IsDevelopment = true
		api.Config.Render.Rest.IndentJSON = true
	}

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
