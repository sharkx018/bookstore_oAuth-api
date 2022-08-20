package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sharkx018/bookstore_oauth-api/src/http"
	"github.com/sharkx018/bookstore_oauth-api/src/repository/db"
	"github.com/sharkx018/bookstore_oauth-api/src/repository/rest"
	"github.com/sharkx018/bookstore_oauth-api/src/service/access_token"
)

var (
	router = gin.Default()
)

func StartApplication() {

	dbRepository := db.NewRepository()
	restRepository := rest.NewRepository()

	service := access_token.NewService(dbRepository, restRepository)
	accessTokenHandler := http.NewHandler(service)

	router.GET("/oauth/access_token/:access_token_id", accessTokenHandler.GetById)
	router.POST("/oauth/access_token/create", accessTokenHandler.Create)

	router.Run(":8080")

}
