package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/api"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Swagger Example API
// @version 1.0.1
func main() {

	client := secretsmanager.New(session.New())
	keystore := services.NewAWSKeyStore(client)

	signTxResource := api.SignTxResourceFactory(keystore)
	generateWalletResource := api.GenerateWalletResourceFactory(keystore)

	port := os.Getenv("PORT")
	app := gin.Default() // create gin app
	
	// Attach the swagger handler
	app.GET("/apidocs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.GET("/generateWallet", generateWalletResource)
	app.POST("/signTx", signTxResource)
	app.Run(":" + port)

}
