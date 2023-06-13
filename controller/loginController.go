package controller

import (
	"fmt"

	"github.com/anodazDev/go_dev/datatransfer"
	"github.com/anodazDev/go_dev/service"
	"github.com/gin-gonic/gin"
)

// interfact
type LoginController interface {
	Login(ctx *gin.Context) (string, string, error)
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

func (controller *loginController) Login(ctx *gin.Context) (string, string, error) {
	var credential datatransfer.LoginCredentials
	err := ctx.ShouldBind(&credential)
	if err != nil {
		return "", "", err
	}

	isUserAuthenticated := controller.loginService.LoginUser(credential.Email, credential.Password)

	if isUserAuthenticated {
		return controller.jwtAuth.GenerateToken(credential.Email, true)
	}
	return "", "", fmt.Errorf("failed to login: ")
}
