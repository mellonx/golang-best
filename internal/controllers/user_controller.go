package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mellonx/golang-best/internal/services"
	"github.com/mellonx/golang-best/pkg/logger"
	"github.com/mellonx/golang-best/pkg/response"
)

// UserController 用户控制器
type UserController struct {
	service services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// List 获取用户列表
func (ctrl *UserController) List(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	logger.InfoCtx(ctx, "Getting user list", "page", page, "page_size", pageSize)

	users, total, err := ctrl.service.List(ctx, page, pageSize)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get user list", "error", err.Error())
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.InfoCtx(ctx, "User list retrieved successfully", "count", len(users))
	response.Success(c, gin.H{
		"data":       users,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// GetByID 根据ID获取用户
func (ctrl *UserController) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID", "error", err.Error())
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	logger.InfoCtx(ctx, "Getting user by ID", "user_id", id)

	user, err := ctrl.service.GetByID(ctx, uint(id))
	if err != nil {
		logger.ErrorCtx(ctx, "User not found", "user_id", id, "error", err.Error())
		response.Error(c, http.StatusNotFound, "user not found")
		return
	}

	logger.InfoCtx(ctx, "User retrieved successfully", "user_id", id)
	response.Success(c, user)
}

// Create 创建用户
func (ctrl *UserController) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid request body", "error", err.Error())
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.InfoCtx(ctx, "Creating user", "username", req.Username, "email", req.Email)

	user, err := ctrl.service.Create(ctx, &req)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to create user", "username", req.Username, "error", err.Error())
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.InfoCtx(ctx, "User created successfully", "user_id", user.ID, "username", user.Username)
	response.Success(c, user)
}

// Update 更新用户
func (ctrl *UserController) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID", "error", err.Error())
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid request body", "error", err.Error())
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.InfoCtx(ctx, "Updating user", "user_id", id)

	user, err := ctrl.service.Update(ctx, uint(id), &req)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to update user", "user_id", id, "error", err.Error())
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.InfoCtx(ctx, "User updated successfully", "user_id", id)
	response.Success(c, user)
}

// Delete 删除用户
func (ctrl *UserController) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID", "error", err.Error())
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	logger.InfoCtx(ctx, "Deleting user", "user_id", id)

	if err := ctrl.service.Delete(ctx, uint(id)); err != nil {
		logger.ErrorCtx(ctx, "Failed to delete user", "user_id", id, "error", err.Error())
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.InfoCtx(ctx, "User deleted successfully", "user_id", id)
	response.Success(c, nil)
}
