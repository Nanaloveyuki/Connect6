package http

import (
	"errors"
	stdhttp "net/http"

	domaingithub "connect6/backend/internal/domain/github"
	"connect6/backend/internal/dto"
	"connect6/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetByUsername(c *gin.Context) {
	user, err := h.userService.GetUser(c.Request.Context(), c.Param("username"))
	if err != nil {
		statusCode := stdhttp.StatusInternalServerError
		if errors.Is(err, domaingithub.ErrUserNotFound) {
			statusCode = stdhttp.StatusNotFound
		}

		c.JSON(statusCode, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(stdhttp.StatusOK, dto.UserResponse{
		Login:      user.Login,
		Name:       user.Name,
		AvatarURL:  user.AvatarURL,
		ProfileURL: user.ProfileURL,
		Bio:        user.Bio,
	})
}

