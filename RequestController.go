package main

import(
	"net/http"
)

var RequestController = map[string]map[string]func(http.ResponseWriter, *http.Request) {

	"GET": GetFunctions,

}

var GetFunctions = map[string]func(http.ResponseWriter, *http.Request) {
	
	
	
}
