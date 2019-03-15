package main

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	vaultApi "github.com/hashicorp/vault/api"
	//ginSwagger "github.com/swaggo/gin-swagger"
	//"github.com/swaggo/gin-swagger/swaggerFiles"
	api "gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/secretStore"
)

// @title Swagger Example API
// @version 1.0.1
func main() {

	config := secretstore.NewConfig()
	hashicorpsSS := secretStore.NewHashicorps(config)
	awsSS = secretStore.NewAWS(7)
	tokenName := os.Getenv("VAULT_TOKEN_NAME")

	err = hashicorpsSS.Init(awsSS, tokenName)
	if err != nil {
		fmt.Printf(err.Error())
	}

	keystore := key.NewBasicKeyStore(hashicorpsSS)

	signTxResource := api.SignTxResourceFactory(keystore)
	generateWalletResource := api.GenerateWalletResourceFactory(keystore)

	port := os.Getenv("PORT")
	app := gin.Default()

	// Attach the swagger handler
	//app.GET("/apidocs/any*", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/generateWallet", generateWalletResource)
	app.POST("/signTx", signTxResource)
	app.Run(":" + port)

}
