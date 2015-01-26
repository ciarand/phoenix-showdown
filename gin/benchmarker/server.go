package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type People []Person
type Person struct {
	Name string
}

// templates is a simple struct we've created to hold all our different
// templates.
type templates struct {
	IndexTemplate *template.Template
}

// tpls is the package instance of templates we'll be using.
var tpls templates

// The way we're doing templates actually deserves some discussion. Go's
// standard template library (used in both text/template and html/template) is
// really good in a lot of ways, but it's also very unopinionated. So this is
// how we're going to define a view chain:
//
// layout.tmpl contains the bare HTML structure. It calls {{template "view" .}}
// index.tmpl  contains the definition of the "view" template
// bio.tmpl    contains a helper partial used by index.tmpl
//
// If we wanted to add another template, all we'd have to do is add one to the
// struct above and make sure it gets parsed in the correct order below.
func init() {
	// we parse these in init so if they're incorrect we fail immediately
	tpls = templates{
		IndexTemplate: template.Must(template.ParseFiles(
			"./templates/layout.tmpl",
			"./templates/bio.tmpl",
			"./templates/index.tmpl")),
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// create a new
	g := gin.New()

	// recover from panics, skip the logger in benchmarks
	g.Use(gin.Recovery() /*, gin.Logger()*/)

	// TODO: Unclear on Martini vs. Gin semantics here.
	//
	// serve the public dir to /public
	// g.Static("/public", "./public")

	// first route
	g.GET("/:title", func(c *gin.Context) {
		title := c.Params.ByName("title")

		members := People{
			Person{Name: "Chris McCord"},
			Person{Name: "Matt Sears"},
			Person{Name: "David Stump"},
			Person{Name: "Ricardo Thompson"},
		}

		// context for template
		context := struct {
			Title   string
			Members People
		}{
			title,
			members,
		}

		c.Writer.Header().Set("Content-Type", "text/html")
		c.Writer.WriteHeader(200)
		if err := tpls.IndexTemplate.ExecuteTemplate(c.Writer, "layout", context); err != nil {
			c.Error(err, nil)
		}
	})

	return http.ListenAndServe(":3000", g)
}
