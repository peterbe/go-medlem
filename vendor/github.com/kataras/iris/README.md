[![Iris Logo](http://iris-go.com/assets/iris_full_logo_2.png)](http://iris-go.com)

[![Travis Widget]][Travis] [![Release Widget]][Release] [![Report Widget]][Report] [![License Widget]][License] [![Gitter Widget]][Gitter] [![Documentation Widget]][Documentation]

[Travis Widget]: https://img.shields.io/travis/tmrts/boilr.svg?style=flat-square
[Travis]: http://travis-ci.org/kataras/iris
[License Widget]: https://img.shields.io/badge/license-Apache%20License%202.0-E91E63.svg?style=flat-square
[License]: https://github.com/kataras/iris/blob/master/LICENSE
[Release Widget]: https://img.shields.io/badge/release-v3.0.0--rc.2-blue.svg?style=flat-square
[Release]: https://github.com/kataras/iris/releases
[Gitter Widget]: https://img.shields.io/badge/chat-on%20gitter-00BCD4.svg?style=flat-square
[Gitter]: https://gitter.im/kataras/iris
[Report Widget]: https://img.shields.io/badge/report%20card-A%2B-F44336.svg?style=flat-square
[Report]: http://goreportcard.com/report/kataras/iris
[Documentation Widget]: https://img.shields.io/badge/documentation-reference-5272B4.svg?style=flat-square
[Documentation]: https://www.gitbook.com/book/kataras/iris/details
[Language Widget]: https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square
[Language]: http://golang.org
[Platform Widget]: https://img.shields.io/badge/platform-Any--OS-gray.svg?style=flat-square

The fastest web framework for Go.

[![Benchmark Wizzard Processing Time Horizontal Graph](https://raw.githubusercontent.com/iris-contrib/website/cf71811e6acb2f9bf1e715e25660392bf090b923/assets/benchmark_horizontal_transparent.png)](https://github.com/smallnest/go-web-framework-benchmark)


```sh
$ cat test_json.go
```
```go
package main

import (
	"github.com/kataras/iris"
)

func main() {

	// render JSON
	iris.Get("/hi_json", func(c *iris.Context) {
		c.JSON(iris.StatusOK, iris.Map{
			"Name":  "Iris",
			"Born":  "13 March 2016",
			"Stars": 2440,
		})
	})
	iris.Listen(":8080")
}
```

> Learn about [configuration](https://kataras.gitbooks.io/iris/content/configuration.html) and [render](https://kataras.gitbooks.io/iris/content/render.html).


```sh
$ cat test_party.go
```
```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

func main() {
	// logger middleware
	log := logger.New(iris.Logger)

	// group routes by path prefix and middleware sharing
	group := iris.Party("/users", log)
	{
		group.Get("/", func(c *iris.Context) {
			// return all users or render a template
		})

		group.Get("/:userID", func(c *iris.Context) {
			// return a user with ID `c.Param("userID")`
		})

		group.Delete("/:userID", func(c *iris.Context) {
			//delete a user with ID `c.Param("userID")`
		})
	}

	// using static subdomains
	subdomain := iris.Party("account.", log, myAuthMiddleware).Layout("layouts/subdomain.html")
	{
		subdomain.Get("/", func(c *iris.Context) {
			// render a template with a context of {username: "myusername"}
			c.Render("account/index.html", iris.Map{ // we can also use a struct
				"username": c.Session().GetString("username"),
			})
		})

		subdomain.Post("/edit", func(c *iris.Context) {
			//...
		})
	}

	// using dynamic subdomains
	dynamicSub := iris.Party("*.")
	{
		// middleware on route, called before the final handler
		dynamicSub.Get("/", log, func(c *iris.Context) {
			c.Write("Hello from subdomain: %s", c.Subdomain())
		})
	}

	iris.Listen(":8080")
}

// using high level sessions inside a custom middleware
func myAuthMiddleware(c *iris.Context) {
	s := c.Session()

	if s.GetString("username") == "myusername" && s.GetString("passowrd") == "mypassword" {
		c.Next()
	} else {
		c.EmitError(iris.StatusUnauthorized)
	}
}

```
> Learn about [named parameters](https://kataras.gitbooks.io/iris/content/named-parameters.html), [parties](https://kataras.gitbooks.io/iris/content/party.html) and [subdomains](https://kataras.gitbooks.io/iris/content/subdomains.html).

Installation
------------
 The only requirement is Go 1.6

`$ go get -u github.com/kataras/iris/iris`

 >If you are connected to the Internet through China [click here](https://kataras.gitbooks.io/iris/content/install.html)

FAQ
------------
You can find answers by exploring [these questions](https://github.com/kataras/iris/issues?q=label%3Aquestion).


Features
------------
- Focus on high performance
- Robust routing & static, wildcard subdomains
- View system supporting [5+](https://kataras.gitbooks.io/iris/content/render_templates.html) template engines
- Highly scalable Websocket API with custom events
- Sessions support with GC, memory & redis providers
- Middlewares & Plugins were never be easier
- Full REST API
- Custom HTTP Errors
- Typescript compiler + Browser editor
- Content negotiation & streaming
- Transport Layer Security
- [Reload](https://github.com/kataras/iris/tree/master/iris#run) on source code changes
- and much more


Docs & Community
------------

<a href="https://www.gitbook.com/book/kataras/iris/details"><img align="right" width="185" src="http://iris-go.com/assets/book/cover_1.png"></a>


- Read the [book](https://www.gitbook.com/book/kataras/iris/details) or [wiki](https://github.com/kataras/iris/wiki)

- Take a look at the [examples](https://github.com/iris-contrib/examples)

- [HISTORY](https://github.com//kataras/iris/tree/master/HISTORY.md) file is your friend.


If you'd like to discuss this package, or ask questions about it, feel free to

* Post an issue or  idea [here](https://github.com/kataras/iris/issues)
* [Chat]( https://gitter.im/kataras/iris) with us

Open debates

 - [Contribute: New website and logo for Iris](https://github.com/kataras/iris/issues/153)
 - [E-book Cover - Which one you suggest?](https://github.com/kataras/iris/issues/67)

**TIP** Be sure to read the [history](HISTORY.md) for Migrating from 2.x to 3.x.

Philosophy
------------

The Iris philosophy is to provide robust tooling for HTTP, making it a great solution for single page applications, web sites, hybrids, or public HTTP APIs.

Iris does not force you to use any specific ORM or template engine. With support for the most used template engines, you can quickly craft the perfect application.

Benchmarks
------------

[This Benchmark suite](https://github.com/smallnest/go-web-framework-benchmark) aims to compare the whole HTTP request processing between Go web frameworks.

![Benchmark Wizzard Processing Time Horizontal Graph](https://raw.githubusercontent.com/iris-contrib/website/cf71811e6acb2f9bf1e715e25660392bf090b923/assets/benchmark_horizontal_transparent.png)

[Please click here to view all detailed benchmarks.](https://github.com/smallnest/go-web-framework-benchmark)

Testing
------------

Iris suggests you to use [this](https://github.com/gavv/httpexpect) new  suite to test your API.
[Httpexpect](https://github.com/gavv/httpexpect) supports fasthttp & Iris after [recommandation](https://github.com/gavv/httpexpect/issues/2). Its author is very active so I believe its a promising library. You can view examples [here](https://github.com/gavv/httpexpect/blob/master/example/iris_test.go) and [here](https://github.com/kataras/iris/blob/master/tests/router_test.go).

Versioning
------------

Current: **v3.0.0-rc.2**
>  Iris is an active project


Read more about Semantic Versioning 2.0.0

 - http://semver.org/
 - https://en.wikipedia.org/wiki/Software_versioning
 - https://wiki.debian.org/UpstreamGuide#Releases_and_Versions


Todo
------------
> for the next release 'v3'

- [x] [Dynamic/Wildcard subdomains](https://kataras.gitbooks.io/iris/content/subdomains.html).
- [x] Create server & client side (js) library for .on('event', func action(...)) / .emit('event')... (like socket.io but supports only websocket).
- [x] Find and provide support for the most stable template engine and be able to change it via the configuration, keep html/templates  support.
- [x] Extend, test and publish to the public the [Iris' cmd](https://github.com/kataras/iris/tree/master/iris).
- [x] Route naming and html url func, requested [here](https://github.com/kataras/iris/issues/165).


If you're willing to donate click [here](DONATIONS.md)

People
------------
A big thanks goes to [all people](https://github.com/kataras/iris/issues?utf8=%E2%9C%93&q=label%3A%22feature+request%22) who help building this framework with feature-requests, bug reports and more!

The author of Iris is [@kataras](https://github.com/kataras).


License
------------

This project is licensed under the Apache License 2.0.

License can be found [here](https://github.com/kataras/iris/blob/master/LICENSE).
