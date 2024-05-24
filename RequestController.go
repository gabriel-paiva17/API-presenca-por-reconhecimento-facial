package main

import(
	"net/http"
)

var RequestController = map[string]map[string]func(http.ResponseWriter, *http.Request) {

	"GET": GetFunctions,
	"POST": PostFunctions,

}

var GetFunctions = map[string]func(http.ResponseWriter, *http.Request) {
	
	
	
}


var PostFunctions = map[string]func(http.ResponseWriter, *http.Request) {
	
	
	
}