package main

import (
	"SycretTest/Http"
	"SycretTest/Services"
	_ "SycretTest/docs"
	"SycretTest/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a SycretTest Test.

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	viper.SetConfigFile("Config.json")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err.Error())
	}
	port := viper.GetString("api.port")
	serviceUrl := viper.GetString("service.url")

	r := gin.Default()
	middL := middleware.InitMiddleware()
	r.Use(middL.Logger())
	docService := Services.NewDocumentService(serviceUrl)

	Http.NewDocumentController(r, docService)

	r.Static("/temp/documents/response", "./temp/documents/response")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(port)

}
