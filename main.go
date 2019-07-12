package main

import (
	"html/template"
	"net/http"
)

func main() {
	var tmpl = template.Must(template.ParseFiles("base.tmpl.html"))

	var rtr = NewRouter(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(404)
		tmpl.Execute(w, struct{ Content template.HTML }{template.HTML("<h1>404 - Page Not Found</h1><br>")})
	})

	rtr.Add("/t/:text", http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		var text = req.URL.Query().Get(":text")
		tmpl.Execute(wr, struct{ Content string }{text})
	}), true)
	rtr.Add("/h/:html", http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		var html = req.URL.Query().Get(":html")
		tmpl.Execute(wr, struct{ Content template.HTML }{template.HTML(html)})
	}), true)
	rtr.ServeMux.Handle("/", rtr.err)

	http.ListenAndServe(":80", rtr)
}
