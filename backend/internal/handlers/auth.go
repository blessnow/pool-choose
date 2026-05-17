package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yuchi/cycle-stock/internal/models"
	"github.com/yuchi/cycle-stock/internal/repository"
)

var Password = "dayuchi"

// Login 登录处理
func Login(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "请求格式错误"})
		return
	}

	if req.Password != Password {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "密码错误"})
		return
	}

	// 创建会话
	token := uuid.New().String()
	session := &models.Session{
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := repository.CreateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "创建会话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "token": token})
}

// AuthMiddleware 认证中间件
func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "未登录"})
		c.Abort()
		return
	}

	session, err := repository.GetSession(token)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "会话已过期"})
		c.Abort()
		return
	}

	c.Set("session", session)
	c.Next()
}

// Logout 登出处理
func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		repository.DeleteSession(token)
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}