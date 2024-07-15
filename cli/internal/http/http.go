package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var url string = "http://localhost:3000/ping"


func Makecall(){
	resp, err := http.Get(url)
	if err != nil {
	log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
 	//Convert the body to type string
	sb := string(body)
	fmt.Println(sb)
}