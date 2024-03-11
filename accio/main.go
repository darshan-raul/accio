package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	http.HandleFunc("/",h1)
	log.Fatal(http.ListenAndServe(":8000",nil))
}