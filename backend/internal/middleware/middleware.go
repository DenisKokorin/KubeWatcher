package middleware

import (
	"k8s-mon/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	JWT jwt.TokenService
}

func NewJWTMiddleware(JWT jwt.TokenService) *JWTMiddleware {
	return &JWTMiddleware{JWT: JWT}
}

func (j *JWTMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token string
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		if token == "" {
			cookie, err := c.Cookie("AccessToken")
			if err == nil && cookie.Value != "" {
				token = cookie.Value
			}
		}

		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "No auth token",
			})
		}

		claims, err := j.JWT.ParseToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid token",
			})
		}
		c.Set("uuid", claims.UUID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		return next(c)
	}
}
