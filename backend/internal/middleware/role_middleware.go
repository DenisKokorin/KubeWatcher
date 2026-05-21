package middleware

import (
	"k8s-mon/internal/roles"
	"k8s-mon/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type RoleMiddleware struct {
	JWT jwt.TokenService
}

func NewRoleMiddleware(JWT jwt.TokenService) *RoleMiddleware {
	return &RoleMiddleware{JWT: JWT}
}

// RequireRole middleware checks if the user has a specific role
func (r *RoleMiddleware) RequireRole(requiredRoles ...roles.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// First, authenticate the token (same as Auth middleware)
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

			claims, err := r.JWT.ParseToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid token",
				})
			}

			// Check if user's role is in the list of required roles
			userRole := roles.Role(claims.Role)
			hasRequiredRole := false
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "Insufficient permissions",
				})
			}

			// Set claims in context for next handler
			c.Set("uuid", claims.UUID)
			c.Set("role", claims.Role)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

// RequirePermission middleware checks if the user has a specific permission
func (r *RoleMiddleware) RequirePermission(requiredPermissions ...roles.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// First, authenticate the token
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

			claims, err := r.JWT.ParseToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "Invalid token",
				})
			}

			// Check if user's role has any of the required permissions
			userRole := roles.Role(claims.Role)
			hasPermission := false
			for _, requiredPerm := range requiredPermissions {
				if userRole.HasPermission(requiredPerm) {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": "Insufficient permissions",
				})
			}

			// Set claims in context for next handler
			c.Set("uuid", claims.UUID)
			c.Set("role", claims.Role)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}
