package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wnjoon/jooncoin/blockchain"
	"github.com/wnjoon/jooncoin/utils"
	"github.com/wnjoon/jooncoin/wallet"
)

type url string
type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

// Struct for Transaction payload
type txPayload struct {
	To     string
	Amount int
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

// Print documentation how to use api
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
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the Blockchain",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get balance of address by TxOut",
		},
	}
	utils.HandleError(json.NewEncoder(rw).Encode(data))
}

// 1. Get Block of hash
// 2. Create block of data
func block(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Get block from its hash
		vars := mux.Vars(r)
		encoder := json.NewEncoder(rw)

		hash := vars["hash"]
		block, err := blockchain.GetBlock(hash)
		if err == utils.ErrBlockNotFound {
			encoder.Encode(utils.ErrorResponse{ErrorMessage: fmt.Sprint(err)})
		} else {
			encoder.Encode(block)
		}
	case "POST":
		blockchain.AddBlock(blockchain.Blockchain())
		rw.WriteHeader(http.StatusCreated)
	}
}

// Get blockchain
func blocks(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain()))
}

// Get status of blockchain
func status(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Blockchain()))
}

// Get balance of address by TxOut
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoder := json.NewEncoder(rw)

	address := vars["address"]
	total := r.URL.Query().Get("total") // If URL has a suffix named "total" and its value => switch to get a total balance or seperated balances

	switch total {
	case "true":
		amount := blockchain.BalanceOfAddressByTxOut(address, blockchain.Blockchain())
		encoder.Encode(balanceResponse{address, amount})
	default:
		err := encoder.Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain()))
		utils.HandleError(err)
	}
}

func transaction(rw http.ResponseWriter, r *http.Request) {
	var payload txPayload
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))

	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(utils.ErrorResponse{ErrorMessage: err.Error()})
		return
	}
	rw.WriteHeader(http.StatusCreated)

}

// Get mempool
func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

/*
 *
 * Minor priority functions
 */
// Marshal text of url
func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

// Create middleware for add http header of json
// Called adapter using http.HandleFunc
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

// Find wallet
func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
}
