package main

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gin-gonic/gin"
	vaultApi "github.com/hashicorp/vault/api"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	api "gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/api"
	hashicorps "gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/hashicorps"
)

// @title Swagger Example API
// @version 1.0.1
func main() {

	vaultConfig := vaultApi.DefaultConfig()

	hashicorpsClient, err := vaultApi.NewClient(vaultConfig)
	if err != nil {
		fmt.Printf(err.Error())
	}

	secretManager := secretsmanager.New(session.New())
	keystore := hashicorps.NewKeyStore(hashicorpsClient, secretManager)

	tokenName := os.Getenv("VAULT_TOKEN_NAME")
	err = keystore.Init(tokenName)
	if err != nil {
		fmt.Printf(err.Error())
	}

	signTxResource := api.SignTxResourceFactory(keystore)
	generateWalletResource := api.GenerateWalletResourceFactory(keystore)

	port := os.Getenv("PORT")
	app := gin.Default()

	// Attach the swagger handler
	app.GET("/apidocs/", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.GET("/generateWallet", generateWalletResource)
	app.POST("/signTx", signTxResource)
	app.Run(":" + port)

}
