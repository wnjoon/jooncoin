package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	router.HandleFunc("/blocks", getAllBlocks).Methods("GET")
	router.HandleFunc("/block/{height:[0-9]+}", getBlockByHeight).Methods("GET")
	router.HandleFunc("/block", createBlock).Methods("POST")

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
			URL:         url("/block/{height}"),
			Method:      "GET",
			Description: "Get a block has id",
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

func getBlockByHeight(rw http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	blockHeight, err := strconv.Atoi((mux.Vars(r))["height"])
	utils.HandleError(err)

	block, err := blockchain.GetBlockchain().GetBlockByHeight(blockHeight)

	encoder := json.NewEncoder(rw)
	if err == utils.ErrBlockNotFound {
		encoder.Encode(utils.ErrorResponse{ErrorMessage: fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func getAllBlocks(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
}

func createBlock(rw http.ResponseWriter, r *http.Request) {
	var body blockBody
	utils.HandleError(json.NewDecoder(r.Body).Decode(&body))
	blockchain.GetBlockchain().AddBlock(body.Message)
	rw.WriteHeader(http.StatusCreated)
}
