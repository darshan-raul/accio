package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
)


type Film struct {

	Title string
	Director string
}

func main(){

	fmt.Println("hello")
	h1 := func( w http.ResponseWriter, r *http.Request){
		tmpl := template.Must(template.ParseFiles("index.html"))
		films := map[string][]Film{
			"Films" :{
				{ Title: "hi", Director: "hi"},
				{ Title: "hi", Director: "hi"},
				{ Title: "hi", Director: "hi"},

			},
		}
		tmpl.Execute(w,films)
	}

	h2 := func( w http.ResponseWriter, r *http.Request){

		log.Print("got request")
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")

		htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s </li>", title, director )
		tmpl,_ := template.New("t").Parse(htmlStr)
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/",h1)
	http.HandleFunc("/add-film/",h2)
	log.Fatal(http.ListenAndServe(":8000",nil))

}