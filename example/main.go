package main

import (
	"html/template"
	"net/http"

	router "github.com/SimplyCodin/gorouter"
)

func main() {
	var tmpl = template.Must(template.ParseFiles("base.tmpl.html"))

	var rtr = router.New(func(wr http.ResponseWriter, req *http.Request) {
		wr.WriteHeader(404)
		err := tmpl.Execute(wr, struct{ Content template.HTML }{
			template.HTML(`
			<header id ="err404" center>
				<h1>404</h1><h2>Page Not Found</h2><h2><a href="/">Go Home</a></h2>
			</header>
			<style>
				html, body, h1 { margin: 0; }
				body { background: url("https://unsplash.com/photos/oFdd0VEPQX4/download") black; background-position-y: center; color: white; }
				a { color: hsla(289, 21%, 60%, 1); text-shadow: 4px 4px 4px rgba(0,0,0,0.5); text-decoration: none; }
				#err404 h1 { font-size: 150px; font-weight: 700; }
				#err404 h2, #err404 a { font-size: 45px; }
			</style>
			`),
		})
		checkErr(err)
	})

	rtr.Add("/t/:text", http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		var text = req.URL.Query().Get(":text")
		err := tmpl.Execute(wr, struct{ Content string }{
			text,
		})
		checkErr(err)
	}), true)

	rtr.Add("/h/:html", http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		var html = req.URL.Query().Get(":html")
		err := tmpl.Execute(wr, struct{ Content template.HTML }{
			template.HTML(html),
		})
		checkErr(err)
	}), true)

	rtr.Add("/", nil, true, http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		err := tmpl.Execute(wr, struct{ Content template.HTML }{
			template.HTML(`
				<h1>GoRouter Example - <a href="https://github.com/SimplyCodin/gorouter">GitHub</a></h1>
				<h2>Please don't use this on a production server</h2>
				<h3>Routes</h3>
				<ul>
					<li><code>/h/:html</code> - Returns HTML in a template</li>
					<li><code>/t/:text</code> - Returns text in a template</li>
				</ul>
			`),
		})
		checkErr(err)
	}))

	err := http.ListenAndServe(":8080", rtr)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		println(err.Error())
	}
}
