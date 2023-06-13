package controller

import (
	"github.com/anodazDev/go_dev/datatransfer"
	"github.com/anodazDev/go_dev/service"
	"github.com/gin-gonic/gin"
)

// interfact
type LoginController interface {
	Login(ctx *gin.Context) string
}

type loginController struct {
	loginService service.LoginService
	jwtAuth      service.JWTService
}

func LoginHandler(loginService service.LoginService,
	jwtAuth service.JWTService) LoginController {
	return &loginController{
		loginService: loginService,
		jwtAuth:      jwtAuth,
	}
}

func (controller *loginController) Login(ctx *gin.Context) string {
	var credential datatransfer.LoginCredentials
	err := ctx.ShouldBind(&credential)
	if err != nil {
		return "no data found"
	}
	isUserAuthenticated := controller.loginService.LoginUser(credential.Email, credential.Password)
	if isUserAuthenticated {
		return controller.jwtAuth.GenerateToken(credential.Email, true)

	}
	return ""
}
