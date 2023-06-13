package main

import (
	"net/http"

	"github.com/anodazDev/go_dev/controller"
	"github.com/anodazDev/go_dev/middleware"
	"github.com/anodazDev/go_dev/service"
	"github.com/gin-gonic/gin"
)

func main() {
	var loginService service.LoginService = service.StaticLoginService()
	var jwtService service.JWTService = service.JWTAuthService()
	var loginController controller.LoginController = controller.LoginHandler(loginService, jwtService)
	// var profileService service.ProfileService

	server := gin.Default()

	server.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	v1 := server.Group("/v1", middleware.AuthorizeJWT())
	// v1.Use(middleware.AuthorizeJWT())
	{
		v1.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})
	}

	profile := server.Group("/", middleware.AuthorizeJWT())
	{
		email, username, tel := service.ProfileUser()
		profile.GET("/profile", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"Email": email, "username": username, "tel": tel})
		})
	}

	port := "8080"
	server.Run(":" + port)

}
