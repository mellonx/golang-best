package services

import (
	"context"
	"errors"

	"github.com/mellonx/golang-best/internal/models"
	"github.com/mellonx/golang-best/internal/repositories"
	"github.com/mellonx/golang-best/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	Create(ctx context.Context, req *CreateUserRequest) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	List(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
	Update(ctx context.Context, id uint, req *UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id uint) error
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	IsActive *bool  `json:"is_active"`
}

// userService 用户服务实现
type userService struct {
	repo repositories.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, req *CreateUserRequest) (*models.User, error) {
	// 检查用户名是否存在
	if _, err := s.repo.GetByUsername(ctx, req.Username); err == nil {
		logger.WarnCtx(ctx, "Username already exists", "username", req.Username)
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否存在
	if _, err := s.repo.GetByEmail(ctx, req.Email); err == nil {
		logger.WarnCtx(ctx, "Email already exists", "email", req.Email)
		return nil, errors.New("email already exists")
	}

	// 加密密码
	logger.DebugCtx(ctx, "Hashing password for user", "username", req.Username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to hash password", "error", err.Error())
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		IsActive: true,
	}

	logger.DebugCtx(ctx, "Creating user in database", "username", req.Username)
	if err := s.repo.Create(ctx, user); err != nil {
		logger.ErrorCtx(ctx, "Failed to create user in database", "error", err.Error())
		return nil, err
	}

	return user, nil
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	logger.DebugCtx(ctx, "Fetching user from database", "user_id", id)
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch user", "user_id", id, "error", err.Error())
		return nil, err
	}
	return user, nil
}

// List 获取用户列表
func (s *userService) List(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	offset := (page - 1) * pageSize
	logger.DebugCtx(ctx, "Fetching user list from database", "page", page, "page_size", pageSize, "offset", offset)

	users, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch user list", "error", err.Error())
		return nil, 0, err
	}

	// 注意：这里应该获取总数，简化示例直接返回0
	return users, int64(len(users)), nil
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, id uint, req *UpdateUserRequest) (*models.User, error) {
	logger.DebugCtx(ctx, "Fetching user for update", "user_id", id)
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch user for update", "user_id", id, "error", err.Error())
		return nil, err
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	logger.DebugCtx(ctx, "Updating user in database", "user_id", id)
	if err := s.repo.Update(ctx, user); err != nil {
		logger.ErrorCtx(ctx, "Failed to update user", "user_id", id, "error", err.Error())
		return nil, err
	}

	return user, nil
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint) error {
	logger.DebugCtx(ctx, "Deleting user from database", "user_id", id)
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "Failed to delete user", "user_id", id, "error", err.Error())
		return err
	}
	return nil
}
