package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mellonx/golang-best/internal/services"
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := ctrl.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// GetByID 根据ID获取用户
func (ctrl *UserController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := ctrl.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "user not found")
		return
	}

	response.Success(c, user)
}

// Create 创建用户
func (ctrl *UserController) Create(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := ctrl.service.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, user)
}

// Update 更新用户
func (ctrl *UserController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := ctrl.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, user)
}

// Delete 删除用户
func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
