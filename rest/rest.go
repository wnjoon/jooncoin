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
	port = utils.SetPort(_port)

	handlers(router)

	utils.PrintConnectionInformation(port)
	log.Fatal(http.ListenAndServe(port, router))
}

func handlers(router *mux.Router) {
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET")
	// router.HandleFunc("/block/{hash:[a-z0-9]*}", block).Methods("GET", "POST")
	router.HandleFunc("/block/{hash:[a-z0-9]+}", block).Methods("GET")
	router.HandleFunc("/block", block).Methods("POST")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/transaction", transaction)
	router.HandleFunc("/mywallet", myWallet)

}
