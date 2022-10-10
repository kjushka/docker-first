package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var counter uint64 = 0

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/")
		w.Write([]byte(strconv.FormatUint(counter, 10)))
	})

	mux.HandleFunc("/stat", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/stat")
		w.Write([]byte(strconv.FormatUint(counter, 10)))
		counter++
	})

	name, _ := os.LookupEnv("NAME")
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/about")
		t, err := template.New("about").Parse(`
			<h3>Hello, {{ .Name }}!</h3>
			<b>Hostname:</b> {{ .Hostname}}<br/>
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err = t.Execute(w, struct {
			Name, Hostname string
		}{
			Name:     name,
			Hostname: r.Host,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Println("start server")
	panic(http.ListenAndServe(":80", mux))
}
