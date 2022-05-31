package rest

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wnjoon/jooncoin/utils"
)

// Create new type for url to adjust MarshalText
var port string

// Description for REST API
// To Use in documentation

func Start(_port int) {
	router := mux.NewRouter()
	port = setPort(_port)

	handlers(router)

	utils.PrintConnectionInformation(port)
	log.Fatal(http.ListenAndServe(port, router))
}
