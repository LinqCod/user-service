package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/linqcod/user-service/internal/domain"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type UserHandler struct {
	logger      *zap.SugaredLogger
	validate    *validator.Validate
	UserUsecase domain.UserUsecase
}

func NewUserHandler(logger *zap.SugaredLogger, validate *validator.Validate, userUsecase domain.UserUsecase) *UserHandler {
	return &UserHandler{
		logger:      logger,
		validate:    validate,
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var user domain.User

	if err := c.ShouldBind(&user); err != nil {
		h.logger.Errorf("error while binding request body: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if err := h.validate.Struct(user); err != nil {
		h.logger.Errorf("error while validating user body: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if err := h.UserUsecase.Create(c, &user); err != nil {
		h.logger.Errorf("error while creating user: %v", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.SuccessResponse{
		Message: "User created successfully",
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Errorf("error while getting user id: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if err := h.UserUsecase.Delete(c, int64(id)); err != nil {
		h.logger.Errorf("error while deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "User deleted successfully",
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Errorf("error while getting user id: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var user domain.User

	if err := c.ShouldBind(&user); err != nil {
		h.logger.Errorf("error while binding request body: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	user.Id = int64(id)

	if err := h.UserUsecase.Update(c, &user); err != nil {
		h.logger.Errorf("error while updating user: %v", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "User updated successfully",
	})
}

func (h *UserHandler) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Errorf("error while getting user id: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	user, err := h.UserUsecase.GetById(c, int64(id))
	if err != nil {
		h.logger.Errorf("error while getting user by id: %v", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessWithDataResponse{
		Message: "User got successfully",
		Data:    user,
	})
}

func (h *UserHandler) GetFilteredUsers(c *gin.Context) {
	count := c.Query("count")
	h.logger.Infof("count: %s", count)
	_, err := strconv.Atoi(count)
	if err != nil && count != "" {
		h.logger.Errorf("error while converting users count to number: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	nationality := c.Query("nationality")

	minAge := c.Query("minAge")
	_, err = strconv.Atoi(minAge)
	if err != nil && minAge != "" {
		h.logger.Errorf("error while converting users min Age to number: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	maxAge := c.Query("maxAge")
	_, err = strconv.Atoi(maxAge)
	if err != nil && maxAge != "" {
		h.logger.Errorf("error while converting users max Age to number: %v", err)
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	gender := c.Query("gender")

	users, err := h.UserUsecase.GetFilteredUsers(c, count, nationality, minAge, maxAge, gender)
	if err != nil {
		h.logger.Errorf("error while getting filtered users: %v", err)
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessWithDataResponse{
		Message: "Users got successfully",
		Data:    users,
	})
}
