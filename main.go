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
		accessToken, refresh_token, err := loginController.Login(ctx)

		if accessToken != "" && err == nil {
			ctx.SetCookie("jwt_token", accessToken, 3600, "/", "localhost", false, true)
			ctx.JSON(http.StatusOK, gin.H{
				"accessToken":  accessToken,
				"refreshToken": refresh_token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	user := server.Group("/user", middleware.AuthorizeJWT())
	{
		email, username, tel := service.ProfileUser()
		user.GET("/profile", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"Email": email, "username": username, "tel": tel})
		})
		user.GET("/logout", func(ctx *gin.Context) {
			ctx.SetCookie("jwt_token", "", -1, "/", "", false, true)

			ctx.JSON(http.StatusOK, gin.H{
				"message": "Logged out successfully",
			})
		})
		user.POST("/refresh", func(ctx *gin.Context) {

			refreshToken := ctx.GetHeader("Refresh-Token")

			newAccessToken, newRefreshToken, err := jwtService.RefreshToken(refreshToken)
			ctx.SetCookie("jwt_token", newAccessToken, 3600, "/", "localhost", false, true)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid refresh token",
				})
				return
			}

			// Return the new access token
			ctx.JSON(http.StatusOK, gin.H{
				"accessToken":  newAccessToken,
				"refreshToken": newRefreshToken,
			})
		})
	}

	port := "8080"
	server.Run(":" + port)

}
