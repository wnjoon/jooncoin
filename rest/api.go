package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wnjoon/jooncoin/blockchain"
	"github.com/wnjoon/jooncoin/utils"
)

type url string

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type blockBody struct {
	// Message includes data
	Message string
}

func handlers(router *mux.Router) {
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET")
	// router.HandleFunc("/block/{hash:[a-z0-9]*}", block).Methods("GET", "POST")
	router.HandleFunc("/block/{hash:[a-z0-9]+}", block).Methods("GET")
	router.HandleFunc("/block", block).Methods("POST")

}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "Get all blocks",
		},
		{
			URL:         url("/block/{hash}"),
			Method:      "GET",
			Description: "Get a block from its hash",
		},
		{
			URL:         url("/block"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
	}
	utils.HandleError(json.NewEncoder(rw).Encode(data))
}

func block(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// Get block from its hash
		vars := mux.Vars(r)
		encoder := json.NewEncoder(rw)

		hash := vars["hash"]
		block, err := blockchain.GetBlock(hash)
		if err == utils.ErrBlockNotFound {
			encoder.Encode(utils.ErrorResponse{fmt.Sprint(err)})
		} else {
			encoder.Encode(block)
		}
	case "POST":
		var body blockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&body))
		blockchain.Blockchain().AddBlock(body.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
}
