package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Server struct{
	Name string
	Type string
	Region string
}

func main(){

	fmt.Println("hello world")

	h1 := func (w http.ResponseWriter , r *http.Request){
		tmpl := template.Must(template.ParseFiles("index.html"))

		servers := map[string][]Server{
			"Servers": {
				{Name: "server1",Type: "t3.micro", Region: "Mumbai"},
				{Name: "server1",Type: "t3.micro", Region: "Pune"},
				{Name: "server1",Type: "t3.medium", Region: "Hyderabad"},
				{Name: "server1",Type: "t4g.micro", Region: "Mumbai"},
			},

		}

		tmpl.Execute(w,servers)
	}
	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		name := r.PostFormValue("name")
		servertype := r.PostFormValue("type")
		fmt.Println(name)
		fmt.Println(servertype)
		//htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", name, type)
		//tmpl, _ := template.New("t").Parse(htmlStr)
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "server-list-element", Server{Name: name, Type: servertype})
	}

	h3 := func(w http.ResponseWriter, r *http.Request){
		fmt.Println("creating server")
		fmt.Print(r.Header)
	}
	http.HandleFunc("/",h1)
	http.HandleFunc("/add-server/", h2)
	http.HandleFunc("/create-server/", h3)
	log.Fatal(http.ListenAndServe(":8000",nil))
}