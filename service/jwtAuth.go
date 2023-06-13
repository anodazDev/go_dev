package service

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwt service
type JWTService interface {
	GenerateToken(email string, isUser bool) (string, string, error)
	ValidateToken(token string) (*jwt.Token, error)
	RefreshToken(refreshToken string) (string, string, error)
}
type authCustomClaims struct {
	Name string `json:"name"`
	User bool   `json:"user"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey       string
	issuer          string
	accessTokenExp  time.Duration
	refreshTokenExp time.Duration
}

// auth-jwt
func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey:       getSecretKey(),
		issuer:          "Bikash",
		accessTokenExp:  time.Hour * 1,
		refreshTokenExp: time.Hour * 24 * 7,
	}
}

func getSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}
func (service *jwtServices) GenerateToken(email string, isUser bool) (string, string, error) {

	accessClaims := &authCustomClaims{
		Name: email,
		User: isUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(service.accessTokenExp).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(service.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshClaims := &authCustomClaims{
		Name: email,
		User: isUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(service.refreshTokenExp).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(service.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return accessTokenString, refreshTokenString, nil
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
}

func (service *jwtServices) RefreshToken(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to parse refresh token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	email, ok := claims["name"].(string)
	if !ok {
		return "", "", fmt.Errorf("failed to extract email from refresh token claims")
	}

	accessToken, refreshToken, err := service.GenerateToken(email, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	return accessToken, refreshToken, nil
}
