package main

import (
	"github.com/freekobie/kora/docs"
	"github.com/freekobie/kora/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/api/v1"

	open := router.Group("/api/v1")
	open.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "200",
			"message": "online",
		})
	})

	// users
	open.POST("/auth/register", app.handler.CreateUser)
	open.POST("/auth/login", app.handler.LoginUser)
	open.POST("/auth/access", app.handler.GetUserAccessToken)
	open.POST("/auth/verify", app.handler.VerifyUser)
	open.POST("/auth/verify/request", app.handler.RequestVerification)

	protected := open.Group("/")
	protected.Use(middlewares.Authentication())
	{
		//users
		protected.GET("/users/:id", app.handler.GetUser)
		protected.PATCH("/users/profile", app.handler.UpdateUserData)
		protected.DELETE("/users/:id", app.handler.DeleteUser)

		// folders
		protected.POST("/folders", app.handler.CreateFolder)

		// files
		protected.POST("/files/upload", app.handler.FileUpload)
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
