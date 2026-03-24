package main

import (
	"context"
	"testing"

	"github.com/mellonx/golang-best/internal/models"
	"github.com/mellonx/golang-best/internal/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewUserRepository(db)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		IsActive: true,
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_ = repo.Create(context.Background(), user)

	// 测试获取
	found, err := repo.GetByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, found.Username)
	assert.Equal(t, user.Email, found.Email)
}
