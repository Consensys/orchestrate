package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
)

// SignTxRequestBody maps fields of a signTx request
type SignTxRequestBody struct {
	Tx SerializedTx `json:"tx" binding:"required"`
	From SerializedAddress `json:"from" binding:"required"`
	Chain JsonifiedChain `json:"chain" binding:"required"`
}

// SignedTxResponse embeddes the signed tx
type SignedTxResponse struct {
	SignedTx []byte `json:"result" binding:"required"`
	Error string `json:"error" binding:"required"`
}

// GenerateWalletResponse embeddes the generated wallet address
type GenerateWalletResponse struct {
	Result string `json:"result" binding:"required"`
	Error string `json:"error" binding:"required"`
}

// SignTxResourceFactory returns the handler to sign a transaction
// @description Returns a signed transaction
// @id signTx
// @tags signature, aws, ethereum, transactions
// @summary Returns a signed transaction
// @accept application/json
// @produce json
// @param tx {object} SerializedTx 1 "Tx Object"
// @param chain {object} JsonifiedChain 1 "Chain object"
// @param from {string} SerializedAddress 1 "Address hex encoded"
// @success 200 {object}
func SignTxResourceFactory(s services.KeyStore) gin.HandlerFunc {

	return func(c *gin.Context) {
	
		var requestBody SignTxRequestBody
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(400, SignedTxResponse{
				SignedTx: []byte{},
				Error: err.Error(), 
			})
			return
		}

		raw, _, err := s.SignTx(
			requestBody.Chain.ToCoreStack(),
			requestBody.From.ToGeth(),
			requestBody.Tx.ToGeth(),
		)

		if err != nil {
			c.JSON(500, SignedTxResponse{
				SignedTx: []byte{},
				Error: err.Error(), 
			})
			return
		}

		response := SignedTxResponse{
			SignedTx: raw,
			Error: "",
		}

		c.JSON(200, response)

	}

}

// GenerateWalletResourceFactory returs a handler to generate a wallet
func GenerateWalletResourceFactory(s services.KeyStore) gin.HandlerFunc {

	return func(c *gin.Context) {

		res, err := s.GenerateWallet()
		if err != nil {
			c.JSON(500, GenerateWalletResponse{
				Result: "",
				Error: err.Error(), 
			})
			return
		}

		c.JSON(200, GenerateWalletResponse{
			Result: res.Hex(),
			Error: "",
		})
	}

}

