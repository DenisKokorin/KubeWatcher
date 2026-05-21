package handler

import (
	"errors"
	"k8s-mon/internal/dto"
	service "k8s-mon/internal/user_service"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const AccessTokenDuration = 60 * 24 * 60 * 60
const RefreshTokenDuration = 60 * 24 * 60 * 60

type UserHandler struct {
	srv service.UserServiceInterface
}

func NewUserHandler(srv service.UserServiceInterface) UserHandler {
	return UserHandler{srv: srv}
}

func (u *UserHandler) Register(c echo.Context) error {
	var userDTO dto.UserDTO

	err := c.Bind(&userDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}

	uid, err := u.srv.CreateUser(c.Request().Context(), userDTO.Name, userDTO.Email, userDTO.Position, userDTO.Team, userDTO.Role, userDTO.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "user already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"uid": uid.String(),
	})
}

func (u *UserHandler) CreateUser(c echo.Context) error {
	var userDTO dto.UserDTO

	err := c.Bind(&userDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}

	if userDTO.Role == "" {
		userDTO.Role = "user"
	}

	uid, err := u.srv.CreateUser(c.Request().Context(), userDTO.Name, userDTO.Email, userDTO.Position, userDTO.Team, userDTO.Role, userDTO.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "user already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"uid": uid.String(),
	})
}

func (u *UserHandler) DeleteUser(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	err = u.srv.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user deleted successfully",
	})
}

func (u *UserHandler) UploadUserDocument(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "file is required",
		})
	}

	uploadedFile, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to open uploaded file",
		})
	}
	defer uploadedFile.Close()

	// Upload file using storage service
	fileInfo, err := u.srv.UploadUserProfileDocument(c.Request().Context(), userID, uploadedFile, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to upload document",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"fileName":  fileInfo.FileName,
			"fileSize":  fileInfo.FileSize,
			"fileType":  fileInfo.FileType,
			"objectKey": fileInfo.ObjectKey,
		},
	})
}

func (u *UserHandler) GetUserDocumentURL(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	presignedURL, err := u.srv.GetUserDocumentURL(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "document not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]string{
			"url": presignedURL,
		},
	})
}

func (u *UserHandler) DeleteUserDocument(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	err = u.srv.DeleteUserDocument(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to delete document",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "document deleted successfully",
	})
}

func (u *UserHandler) Login(c echo.Context) error {
	var loginReq dto.LoginDTO

	err := c.Bind(&loginReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}
	ip := c.RealIP()
	AccessToken, RefreshToken, err := u.srv.Login(c.Request().Context(), loginReq.Email, loginReq.Password, c.Request().UserAgent(), ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	err = setCookie(c, "AccessToken", AccessToken, AccessTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	err = setCookie(c, "RefreshToken", RefreshToken, RefreshTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"AccessToken":  AccessToken,
		"RefreshToken": RefreshToken,
	})
}

func (u *UserHandler) Logout(c echo.Context) error {
	// First, get the refresh token before clearing cookies
	refresh, err := c.Cookie("RefreshToken")

	// Delete the JWT session from database/Redis if refresh token exists
	if err == nil && refresh.Value != "" {
		if err := u.srv.Logout(c.Request().Context(), refresh.Value); err != nil {
			// Log the error but don't fail logout - still clear cookies
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "failed to invalidate session",
			})
		}
	}

	// Clear the cookies after logout
	err = setCookie(c, "AccessToken", "", 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	err = setCookie(c, "RefreshToken", "", 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logged out successfully",
	})
}

func setCookie(c echo.Context, key, value string, duration int) error {
	cookie := http.Cookie{
		Name:     key,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   duration,
		Value:    value,
		Path:     "/",
	}
	c.SetCookie(&cookie)
	return nil
}

func (u *UserHandler) Refresh(c echo.Context) error {
	oldRefresh, err := c.Cookie("RefreshToken")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	refreshUUID, err := uuid.Parse(oldRefresh.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	expiredAccesToken, err := c.Cookie("AccessToken")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	newRefresh, newAccess, err := u.srv.RefreshToken(c.Request().Context(), refreshUUID, c.Request().UserAgent(), c.RealIP(), expiredAccesToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshSession) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "invalid refresh session",
			})
		}
		if errors.Is(err, service.ErrTokenExpired) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "token expired",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	err = setCookie(c, "RefreshToken", newRefresh.String(), RefreshTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	err = setCookie(c, "AccessToken", newAccess, AccessTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"AccessToken":  newAccess,
		"RefreshToken": newRefresh,
	})
}

func (u *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get("uuid")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "unauthorized",
		})
	}

	userUUID, ok := userID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}

	uid, err := uuid.Parse(userUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}

	user, err := u.srv.GetUser(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"uuid":     user.UUID,
		"name":     user.Name,
		"email":    user.Email,
		"position": user.Position,
		"team":     user.Team,
		"role":     user.Role,
	})
}

func (u *UserHandler) GetAllUsers(c echo.Context) error {
	users, err := u.srv.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
	})
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func (u *UserHandler) UpdateUserRole(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid user id",
		})
	}

	var req UpdateRoleRequest
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}

	err = u.srv.UpdateUserRole(c.Request().Context(), userID, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "role updated successfully",
	})
}
